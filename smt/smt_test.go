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
package smt_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wolkdb/deepblockchains/smt"
)

func TestSMT(t *testing.T) {
	cs, err := smt.NewCloudstore(smt.DefaultChunkstorePath)
	if err != nil {
		t.Fatalf("[smt_test:NewCloudstore]%v", err)
	}
	defer cs.Close()

	s := smt.NewSparseMerkleTree(cs)
	nkeys := uint64(500)
	nversions := uint64(3)
	chunkHash := make(map[uint64]common.Hash)
	merkleRoot := make(map[uint64]common.Hash)
	kv := make(map[uint64]map[uint64]common.Hash)
	for ver := uint64(0); ver < nversions; ver++ {
		kv[ver] = make(map[uint64]common.Hash)
		for i := uint64(0); i < nkeys; i++ {
			storageBytesNew := uint64(3)
			k := smt.Bytes32ToUint64(smt.Computehash(smt.Uint64ToBytes32(i % 10000)))
			v := smt.Computehash([]byte(fmt.Sprintf("%d%d", i, ver)))
			kv[ver][k] = common.BytesToHash(v)
			prevBlock := ver
			err = s.Insert(smt.UIntToByte(k), v, storageBytesNew, prevBlock)
			if err != nil {
				t.Fatalf("SetKey: %v\n", err)
			}
		}
		s.Flush()
		chunkHash[ver] = s.ChunkHash()
		merkleRoot[ver] = s.MerkleRoot()
		fmt.Printf("Generated: Version %d Hash: %x Merkle Root: %x\n", ver, chunkHash[ver], merkleRoot[ver])
	}

	for ver := uint64(0); ver < nversions; ver++ {
		s = smt.NewSparseMerkleTree(cs)
		s.InitWithRoot(merkleRoot[ver])
		passes := 0
		for i := uint64(0); i < nkeys; i++ {
			k := smt.Bytes32ToUint64(smt.Computehash(smt.Uint64ToBytes32(i % 10000)))
			v1, found, proof, storageBytes, prevBlock, err := s.Get(smt.UIntToByte(k))
			if err != nil {
				fmt.Printf("err not found %x %v \n", k, err)
			} else if found {
				if bytes.Compare(kv[ver][k].Bytes(), v1) == 0 {
					checkproof := proof.Verify(v1, merkleRoot[ver].Bytes(), false)
					if checkproof {
						passes++
					} else {
						fmt.Printf("k:%x v:%x storageBytes:%d prevBlock: %d ver %d -- ", k, v1, storageBytes, prevBlock, ver)
						t.Fatalf("CHECK PROOF ==> FAILURE\n")
					}

				} else {
					t.Fatalf("k:%x v:%x sb:%d kv[k]:%x INCORRECT\n", k, v1, storageBytes, kv[k])
				}
			} else {
				fmt.Printf("k:%x not found \n", k)
			}
		}
		fmt.Printf("Version %d  -- %d/%d keys PASSED\n", ver, passes, nkeys)
	}
}
