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
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func ComputeDefaultHashes() (defaultHashes [TreeDepth][]byte) {
	empty := make([]byte, 0)
	defaultHashes[0] = Keccak256(empty)
	for level := 1; level < TreeDepth; level++ {
		defaultHashes[level] = Keccak256(defaultHashes[level-1], defaultHashes[level-1])
	}
	return defaultHashes
}

func Keccak256(data ...[]byte) []byte {
	hasher := sha3.NewKeccak256()
	for _, b := range data {
		hasher.Write(b)
	}
	return hasher.Sum(nil)
}

// helper stuff here for a while
func IntToByte(i int64) (k []byte) {
	k = make([]byte, 8)
	binary.BigEndian.PutUint64(k, uint64(i))
	return k
}

func UIntToByte(i uint64) (k []byte) {
	k = make([]byte, 8)
	binary.BigEndian.PutUint64(k, uint64(i))
	return k
}

func UInt64ToByte(i uint64) (k []byte) {
	k = make([]byte, 8)
	binary.BigEndian.PutUint64(k, uint64(i))
	return k
}

func Uint64ToBytes32(i uint64) (k []byte) {
	k = make([]byte, 32)
	binary.BigEndian.PutUint64(k[24:32], uint64(i))
	return k
}

func BytesToUint64(inp []byte) uint64 {
	return binary.BigEndian.Uint64(inp)
}

func Bytes32ToUint64(k []byte) (out uint64) {
	h := k[0:8]
	return BytesToUint64(h)
}
