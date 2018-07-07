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
	"bytes"
	"fmt"
)

type Proof struct {
	key       []byte
	proof     [][]byte
	proofBits uint64
}

func (self *Proof) Check(v []byte, root []byte, defaultHashes [TreeDepth][]byte, verbose bool) bool {
	// the leaf value to start off hashing!  The value is hash(RLPEncode([]))
	debug := false
	cur := v
	p := 0

	for i := uint64(0); i < 64; i++ {

		if (uint64(1<<i) & self.proofBits) > 0 {
			if byte(0x01<<(i%8))&byte(self.key[(TreeDepth-1-i)/8]) > 0 { // i-th bit is "1", so hash with H([]) on the left
				if debug {
					fmt.Printf("C%v | [P,*] bit%v=1 | H(P[%d]:%x, C[%d]:%x) => ", i+1, i, p, self.proof[p], i, cur)
				}
				cur = Keccak256(self.proof[p], cur)
			} else { // i-th bit is "0", so hash with H([]) on the right
				if debug {
					fmt.Printf("C%v | [*,P] bit%v=0 | H(C[%d]:%x, P[%d]:%x) => ", i+1, i, i, cur, p, self.proof[p])
				}
				cur = Keccak256(cur, self.proof[p])
			}
			p++
		} else {
			if byte(0x01<<(i%8))&byte(self.key[(TreeDepth-1-i)/8]) > 0 { // i-th bit is "1", so hash with H([]) on the left
				if debug {
					fmt.Printf("C%v | [D,*] bit%v=1 | H(D[%d]:%x, C[%d]:%x) => ", i+1, i, i, defaultHashes[i], i, cur)
				}
				cur = Keccak256(defaultHashes[i], cur)
			} else {
				if debug {
					fmt.Printf("C%v | [*,D] bit%v=0 | H(C[%d]:%x, D[%d]:%x) => ", i+1, i, i, cur, i, defaultHashes[i])
				}
				cur = Keccak256(cur, defaultHashes[i])
			}
		}
		if debug {
			fmt.Printf(" %x\n", cur)
		}
	}
	res := bytes.Compare(cur, root) == 0
	if verbose {
		if res {
			fmt.Printf(" CheckProof success (proof matched root: %x)\n", root)
		} else {
			fmt.Printf(" CheckProof FAILURE (proof does NOT match root: %x)\n", root)
		}
	}
	return res
}

func (self *Proof) String() string {
	out := fmt.Sprintf("{\"token\":\"%x\",\"proofBits\":\"%x\",\"proof\":[", self.key, self.proofBits)
	for i, p := range self.proof {
		if i > 0 {
			out = out + ","
		}
		out = out + fmt.Sprintf("\"0x%x\"", p)
	}
	out = out + "]}"
	return out
}

func (p *Proof) Bytes() (out []byte) {
	out = append(out, UInt64ToByte(p.proofBits)...)
	for _, h := range p.proof {
		out = append(out, h...)
	}
	return out
}
