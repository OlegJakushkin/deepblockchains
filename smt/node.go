// Copyright 2018 Wolk Inc.
// This file is part of the SMT library.
//
// The SMT library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The SMT library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
package smt

import (
	"bytes"
	"fmt"
)

const (
	TreeDepth     = 64
	debug         = false
	nchildren     = 16
	bytesPerChild = 88
	chunkSize     = 4096
	bitsPerPiece  = 4
)

var GlobalDefaultHashes [TreeDepth][]byte

func init() {
	GlobalDefaultHashes = ComputeDefaultHashes()
}

func keypiece(k []byte, i int) uint8 {
	if i%2 == 0 { // top 4 bits if i is even
		return (k[i/2] >> 4) & 0x0F
	} else { // bottom 4 bits if i is odd
		return (k[i/2]) & 0x0F
	}
}

func nkeypieces(k []byte) int {
	return len(k) * 2
}

type Node struct {
	level        int    // position of node for information/logging only
	chunkHash    []byte // for cloud storage
	merkleRoot   []byte // for proofs
	unloaded     bool
	terminal     bool
	key          []byte
	dirty        bool
	storageBytes uint64
	blockNum     uint64
	mrcache      [bitsPerPiece + 1][nchildren][]byte
	children     []*Node
}

func NewNode(level int, parent *Node) *Node {
	return &Node{
		children:     make([]*Node, nchildren),
		terminal:     false,
		dirty:        false,
		unloaded:     false,
		level:        level,
		storageBytes: 0,
	}
}

func (n *Node) SetHash(hash []byte) {
	n.unloaded = true
	n.chunkHash = hash
}

func (n *Node) generateProof(cs Cloudstore, k []byte, v []byte, i int, p *Proof) (ok bool) {
	if n.unloaded {
		n.load(cs)
	}
	idx := keypiece(k, i)
	if n.children[idx] == nil {
		return false // not a member of the tree!
	}
	if n.children[idx].unloaded {
		n.children[idx].load(cs)
	}
	if n.children[idx].terminal {
		p.proofBits = 0
	} else {
		ok = n.children[idx].generateProof(cs, k, v, i+1, p)
		if !ok {
			return false
		}
	}

	n.computeMerkleRootCache()

	for level := bitsPerPiece; level > 0; level-- {
		// fmt.Printf("idx %08b @ level %d ==> PROOF OUTPUT %x\n", idx, level, n.mrcache[level][idx])
		sister_index := idx
		if idx&1 > 0 {
			sister_index = idx - 1
		} else {
			sister_index = idx + 1
		}
		p0 := n.mrcache[level][sister_index]
		if bytes.Compare(p0, GlobalDefaultHashes[TreeDepth-level]) == 0 {
			// fmt.Printf(" ---- %x GlobalDefaultHashes[%d]: %x\n", p0, TreeDepth-level, GlobalDefaultHashes[TreeDepth-level])
		} else {
			p.proof = append(p.proof, n.mrcache[level][sister_index])
			p.proofBits |= (uint64(1) << uint64(n.level-level+1))
			// fmt.Printf(" ^^^^ proof[%d]: %x @ level %d n.level=%d [%d] -- proofBits: (%x)\n", len(p.proof)-1, p.proof[len(p.proof)-1], level, n.level, n.level-level+1, p.proofBits)
		}
		idx = idx >> 1
	}
	// fmt.Printf("generateProof %x k[%d] = %x ---> len(proof) = %d\n", k, i, keypiece(k, i), len(p.proof))
	return true
}

func (n *Node) computeMerkleRootCache() {
	// now for each of 8...0 levels, hash the level of "leaves" into  n.mrcache
	newleaves_cnt := nchildren / 2
	startlevel := bitsPerPiece - 1
	for level := startlevel; level >= 0; level-- {
		dh := GlobalDefaultHashes[n.level-startlevel]
		for i := 0; i < newleaves_cnt; i++ {
			if level == startlevel {
				// now, using the children's merkle roots (or, dh if child is nil),
				l := n.children[i*2]
				r := n.children[i*2+1]
				if l != nil && r != nil {
					n.mrcache[bitsPerPiece][i*2] = l.merkleRoot[:]
					n.mrcache[bitsPerPiece][i*2+1] = r.merkleRoot[:]
				} else if l != nil && r == nil {
					n.mrcache[bitsPerPiece][i*2] = l.merkleRoot[:]
					n.mrcache[bitsPerPiece][i*2+1] = dh
				} else if l == nil && r != nil {
					n.mrcache[bitsPerPiece][i*2] = dh
					n.mrcache[bitsPerPiece][i*2+1] = r.merkleRoot[:]
				} else {
					n.mrcache[bitsPerPiece][i*2] = dh
					n.mrcache[bitsPerPiece][i*2+1] = dh
				}
			}
			n.mrcache[level][i] = Computehash(n.mrcache[level+1][i*2], n.mrcache[level+1][i*2+1])
			if bytes.Compare(n.mrcache[level+1][i*2], n.mrcache[level+1][i*2+1]) != 0 {
				if debug {
					fmt.Printf(" G mrcache[level %d][i %d]: %x", level, i, n.mrcache[level][i])
					fmt.Printf(" = Keccak256( %02x: %x, %02x: %x )\n", i*2, n.mrcache[level+1][i*2], i*2+1, n.mrcache[level+1][i*2+1])
				}
			}
		}
		newleaves_cnt = newleaves_cnt / 2
	}

	n.merkleRoot = make([]byte, 32)
	copy(n.merkleRoot[:], n.mrcache[0][0][:])
}

func (n *Node) computeMerkleRoot(cs Cloudstore) []byte {
	n.load(cs)

	if n.terminal {
		// build merkle root hash of a terminal going from level 0 to the terminal Level
		n.merkleRoot = make([]byte, 32)
		copy(n.merkleRoot[:], n.chunkHash[:])
		// the leaf value to start off hashing!  The value is hash(RLPEncode([]))
		cur := n.chunkHash
		// for each bit 0 up to the terminal level, use the bit # to hash cur on the left (0) or right (1)
		for i := uint(0); i <= uint(n.level); i++ {
			if byte(0x01<<(i%8))&byte(n.key[(TreeDepth-1-i)/8]) > 0 { // i-th bit is "1", so hash with H([]) on the left
				if debug && (i < 1 || i > 54) {
					fmt.Printf(" mr %x bit %3d (%08b)=1 hash(GlobalDefaultHashes[%d]:%x, cur:%x) => ", n.key, i, i, i, GlobalDefaultHashes[i], cur)
				}
				cur = Computehash(GlobalDefaultHashes[i], cur)
			} else { // i-th bit is "0", so hash with H([]) on the right
				if debug && (i < 1 || i > 54) {
					fmt.Printf(" mr %x bit %3d (%08b)=0 hash(cur:%x, GlobalDefaultHashes[%d]:%x) => ", n.key, i, i, cur, i, GlobalDefaultHashes[i])
				}
				cur = Computehash(cur, GlobalDefaultHashes[i])
			}
			if debug && (i < 1 || i > 54) {
				fmt.Printf(" %x\n", cur)
			}
		}
		// copy the answer and return!
		copy(n.merkleRoot, cur)
		return n.merkleRoot
	} else {
		// ok, we are not the terminal, so for each child, compute THEIR merkle root at the child Level
		for i := 0; i < nchildren; i++ {
			if n.children[i] != nil {
				n.children[i].computeMerkleRoot(cs)
			}
		}
	}

	// now for each of 8...0 levels, hash the level of "leaves" into  n.mrcache
	newleaves_cnt := nchildren / 2
	startlevel := bitsPerPiece - 1
	for level := startlevel; level >= 0; level-- {
		dh := GlobalDefaultHashes[n.level-startlevel]
		for i := 0; i < newleaves_cnt; i++ {
			if level == startlevel {
				// now, using the children's merkle roots (or, dh if child is nil),
				l := n.children[i*2]
				r := n.children[i*2+1]
				if l != nil && r != nil {
					n.mrcache[bitsPerPiece][i*2] = l.merkleRoot[:]
					n.mrcache[bitsPerPiece][i*2+1] = r.merkleRoot[:]
				} else if l != nil && r == nil {
					n.mrcache[bitsPerPiece][i*2] = l.merkleRoot[:]
					n.mrcache[bitsPerPiece][i*2+1] = dh
				} else if l == nil && r != nil {
					n.mrcache[bitsPerPiece][i*2] = dh
					n.mrcache[bitsPerPiece][i*2+1] = r.merkleRoot[:]
				} else {
					n.mrcache[bitsPerPiece][i*2] = dh
					n.mrcache[bitsPerPiece][i*2+1] = dh
				}
			}
			n.mrcache[level][i] = Computehash(n.mrcache[level+1][i*2], n.mrcache[level+1][i*2+1])
			if bytes.Compare(n.mrcache[level+1][i*2], n.mrcache[level+1][i*2+1]) != 0 {
				if debug {
					fmt.Printf(" G mrcache[level %d][i %d]: %x", level, i, n.mrcache[level][i])
					fmt.Printf(" = Keccak256( %02x: %x, %02x: %x )\n", i*2, n.mrcache[level+1][i*2], i*2+1, n.mrcache[level+1][i*2+1])
				}
			}
		}
		newleaves_cnt = newleaves_cnt / 2
	}
	// finally, store the answer
	n.computeMerkleRootCache()
	return n.merkleRoot
}

func (n *Node) delete(k []byte, i int) (ok bool, err error) {
	if i >= nkeypieces(k) {
		return false, fmt.Errorf(" we hit the bottom WOW!")
	} else {
		idx := keypiece(k, i)
		if n.children[idx] == nil {
			return false, nil
		} else {
			if n.children[idx].terminal {
				if bytes.Compare(k, n.children[idx].key) == 0 {
					n.children[idx] = nil
					n.dirty = true
					// TODO: if this blanks everything out, do something!
				} else {
					return false, nil
				}
			} else {
				ok, err = n.children[idx].delete(k, i+1)
				if ok {
					n.dirty = true
				}
				return ok, err
			}
		}
	}
	return false, nil
}

func (n *Node) insert(cs Cloudstore, k []byte, v []byte, i int, storageBytesNew uint64, blockNum uint64) error {
	if n.unloaded {
		n.load(cs)
	}

	if i >= nkeypieces(k) {
		// we hit the bottom WOW!
	} else {
		n.dirty = true
		idx := keypiece(k, i)
		if n.children[idx] == nil {
			//fmt.Printf(" -- new child %d %x\n", i, keypiece(k,i))
			n.children[idx] = NewNode(n.level-bitsPerPiece, n)
			n.children[idx].terminal = true
			n.children[idx].key = k
			n.children[idx].chunkHash = v
			n.children[idx].dirty = true
			n.children[idx].storageBytes = storageBytesNew
			n.children[idx].blockNum = blockNum
		} else {
			if n.children[idx].terminal {
				//fmt.Printf(" -- SPLIT child @ level %d : %02x (%x)\n", n.children[idx].level, keypiece(k,i), k)
				if bytes.Compare(n.children[idx].key, k) == 0 {
					if bytes.Compare(n.children[idx].chunkHash, v) == 0 {
						// nothing changed...
					} else {
						n.children[idx].dirty = true
						n.children[idx].chunkHash = v
						n.children[idx].blockNum = blockNum
						n.children[idx].storageBytes += storageBytesNew
					}
				} else {
					// two keys!
					n.children[idx].dirty = true
					n.children[idx].terminal = false
					n.children[idx].insert(cs, n.children[idx].key, n.children[idx].chunkHash, i+1, n.children[idx].storageBytes, n.children[idx].blockNum)
					n.children[idx].insert(cs, k, v, i+1, storageBytesNew, blockNum)
				}
			} else {
				n.children[idx].insert(cs, k, v, i+1, storageBytesNew, blockNum)
			}
		}
		tot := uint64(0)
		for i := 0; i < nchildren; i++ {
			if n.children[i] != nil {
				tot += n.children[i].storageBytes
			}
		}
		n.storageBytes = tot
	}
	return nil
}

// load from SWARM using self.chunkHash; node chunks are saved in "flush" operations
// node chunks are nchildren rows, each of 40 bytes: 8 byte keys and 32 byte hashes.  If the 32 byte hash is 0, then there is no child.
func (self *Node) load(cs Cloudstore) bool {
	if !self.unloaded {
		return false
	}
	chunk, ok, err := cs.GetChunk(self.chunkHash)
	//fmt.Printf("\nThe chunk retrieved using hash (%x): %+x (%+v) and ERR: %+v OK: %+v", self.chunkHash, chunk, chunk, err, ok)
	if err != nil {
		return false
	} else if !ok {
		return true
	} else {
		blank_key := make([]byte, 8)
		blank_hash := make([]byte, 32)
		for j := 0; j < nchildren; j++ {
			//TODO: check that these chunk indices are valid before calling
			p_key := chunk[j*bytesPerChild+0 : j*bytesPerChild+8]
			p_hash := chunk[j*bytesPerChild+8 : j*bytesPerChild+40]
			if bytes.Compare(p_hash, blank_hash) != 0 {
				self.children[j] = NewNode(self.level-bitsPerPiece, self)
				if bytes.Compare(p_key, blank_key) == 0 { // its NOT a terminal
					self.children[j].terminal = false
					self.children[j].unloaded = true
				} else { // it IS a terminal
					self.children[j].unloaded = false
					self.children[j].terminal = true
					self.children[j].key = p_key // 8 bytes
				}
				self.children[j].chunkHash = p_hash
				self.children[j].storageBytes = Bytes32ToUint64(chunk[j*bytesPerChild+40 : j*bytesPerChild+48])
				self.children[j].merkleRoot = chunk[j*bytesPerChild+48 : j*bytesPerChild+80]
				self.children[j].blockNum = Bytes32ToUint64(chunk[j*bytesPerChild+80 : j*bytesPerChild+88])
			}
		}
	}
	self.unloaded = false
	return true
}

func (self *Node) get(cs Cloudstore, k []byte, i int) (v []byte, ok bool, storageBytes uint64, blockNum uint64, err error) {
	self.load(cs)
	//TODO: load siblings along with desired chunk / child
	idx := keypiece(k, i)
	if self.children[idx] != nil {
		if self.children[idx].terminal {
			if bytes.Equal(k, self.children[idx].key) {
				return self.children[idx].chunkHash, (bytes.Compare(self.children[idx].key, k) == 0), self.children[idx].storageBytes, self.children[idx].blockNum, nil
			}
		} else {
			v, ok, storageBytes, blockNum, err = self.children[idx].get(cs, k, i+1)
			if err != nil {
				return v, ok, storageBytes, blockNum, err
			} else {
				return v, ok, storageBytes, blockNum, err
			}
		}
	}
	return v, false, 0, 0, nil
}

func (self *Node) flush(cs Cloudstore) (err error) {
	if self.dirty {
		// compute hash!
		chunk := make([]byte, chunkSize)
		for i := 0; i < nchildren; i++ {
			if self.children[i] != nil {
				if self.children[i].terminal {
					copy(chunk[i*bytesPerChild:i*bytesPerChild+8], self.children[i].key)
					copy(chunk[i*bytesPerChild+8:i*bytesPerChild+40], self.children[i].chunkHash)
					copy(chunk[i*bytesPerChild+40:i*bytesPerChild+48], UIntToByte(self.children[i].storageBytes))
					copy(chunk[i*bytesPerChild+48:i*bytesPerChild+80], self.children[i].merkleRoot)
					copy(chunk[i*bytesPerChild+80:i*bytesPerChild+88], UIntToByte(self.children[i].blockNum))
				} else {
					// recursive call
					err = self.children[i].flush(cs)
					if err != nil {
						return err
					} else {
						// the top level hash of the child has been computed, so write it into this chunk
						copy(chunk[i*bytesPerChild+8:i*bytesPerChild+40], self.children[i].chunkHash)
						copy(chunk[i*bytesPerChild+40:i*bytesPerChild+48], UIntToByte(self.children[i].storageBytes))
						copy(chunk[i*bytesPerChild+48:i*bytesPerChild+80], self.children[i].merkleRoot)
						copy(chunk[i*bytesPerChild+80:i*bytesPerChild+88], UIntToByte(self.children[i].blockNum))
					}
				}
			} else {
				blank := make([]byte, bytesPerChild)
				copy(chunk[i*bytesPerChild:(i+1)*bytesPerChild], blank)
			}
		}
		// store newly developed chunk to Cloudstore
		chunkID := Computehash(chunk)
		//err := cs.StoreChunk(chunkID, chunk)
		go func() {
			err := cs.SetChunk(chunkID, chunk)
			if err != nil {
				//return err
			}
		}()

		self.chunkHash = chunkID
		self.dirty = false
	}
	return nil
}

func (self *Node) flushRoot(cs Cloudstore) (err error) {
	//create (merkleroot, chunkHash) mapping
	err = cs.SetChunk(self.merkleRoot, self.chunkHash)
	if err != nil {
		return err
	}
	return nil
}

func (self *Node) dump(prefix []byte) error {
	for i := 0; i < (TreeDepth-1-self.level)/bitsPerPiece; i++ {
		fmt.Printf("  ")
	}
	fmt.Printf("[")
	for j := 0; j < len(prefix); j++ {
		pr := fmt.Sprintf("%x", prefix[j])
		if len(pr) == 2 {
			pr = pr[1:2]
		}
		fmt.Printf("%s", pr)
	}
	fmt.Printf("] ")

	fmt.Printf("Level %d Hash %x StorageBytes: %d  MerkleRoot: %x", self.level, self.chunkHash, self.storageBytes, self.merkleRoot)
	if self.unloaded {
		fmt.Printf(" (UNLOADED)")
	}
	if self.terminal {
		fmt.Printf(" ** TERM KEY: %x VAL: %x StorageBytes: %d blockNum: %d MerkleRoot: %x\n", self.key, self.chunkHash, self.storageBytes, self.blockNum, self.merkleRoot)
	} else {
		fmt.Printf("\n")
	}
	for i := 0; i < nchildren; i++ {
		if self.children[i] != nil {
			out := make([]byte, len(prefix)+1)
			copy(out, prefix)
			out[len(prefix)] = byte(i)
			self.children[i].dump(out)
		}
	}
	return nil
}
