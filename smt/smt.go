/** This file is part of the SMT library.
*
* The SMT library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The SMT library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
*
* You should have received a copy of the GNU Lesser General Public License
* along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
 */
package smt

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type SparseMerkleTree struct {
	Cloudstore    Cloudstore
	root          *Node
	DefaultHashes [TreeDepth][]byte
}

func NewSparseMerkleTree(cs Cloudstore) *SparseMerkleTree {
	var self SparseMerkleTree
	self.root = NewNode(TreeDepth-1, nil)
	self.Cloudstore = cs
	self.DefaultHashes = ComputeDefaultHashes()
	return &self
}

func (self *SparseMerkleTree) Init(hash common.Hash) {
	self.root.chunkHash = hash.Bytes()
	self.root.unloaded = true
}

func (self *SparseMerkleTree) Flush() {
	self.root.computeMerkleRoot(self.Cloudstore, self.DefaultHashes)
	self.root.flush(self.Cloudstore)
}

func (self *SparseMerkleTree) Delete(k []byte) error {
	self.root.delete(k, 0)
	return nil
}

func (self *SparseMerkleTree) GenerateProof(k []byte, v []byte) (p *Proof) {
	var pr Proof
	pr.key = k
	self.root.generateProof(self.Cloudstore, k, v, 0, self.DefaultHashes, &pr)
	return &pr
}

func (self *SparseMerkleTree) Get(k []byte) (v0 []byte, found bool, p *Proof, storageBytes uint64, prevBlock uint64, err error) {
	v0, found, storageBytes, prevBlock, err = self.root.get(self.Cloudstore, k, 0)
	if found {
		var pr Proof
		pr.key = k
		ok := self.root.generateProof(self.Cloudstore, k, v0, 0, self.DefaultHashes, &pr)
		if !ok {
			return v0, found, &pr, storageBytes, prevBlock, fmt.Errorf("NO proof")
		}
		return v0, found, &pr, storageBytes, prevBlock, nil
	}
	return v0, found, nil, storageBytes, prevBlock, nil
}

func (self *SparseMerkleTree) ChunkHash() common.Hash {
	return common.BytesToHash(self.root.chunkHash)
}

func (self *SparseMerkleTree) MerkleRoot() common.Hash {
	return common.BytesToHash(self.root.merkleRoot)
}

func (self *SparseMerkleTree) Insert(k []byte, v []byte, storageBytesNew uint64, blockNum uint64) error {
	// fmt.Printf(" SMT insert(k:%x, v:%x, bn: %d)\n", k, v, blockNum)
	self.root.insert(self.Cloudstore, k, v, 0, storageBytesNew, blockNum)
	return nil
}

func (self *SparseMerkleTree) Dump() {
	self.root.dump(nil)
}
