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
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestCheckProof(t *testing.T) {
	var p1, p2 Proof

	p1Leaf := common.Hex2Bytes("4fdcebb3247a9a715e416e68439e563e8faf57c804642441b93724d5b4fe0878")
	p1Root := common.Hex2Bytes("f361d4563fec05b5262f16c96aa062924256f61bd7482213ae23bf8bb2ad2e69")
	p1, _ = ToProof(uint64(0x9af84bc1208918b), common.Hex2Bytes("e000000000000000fb0d81010243cb5171ab9e619ca2a996a2f6eb2505a80b2e0252349c0f8d09105a581781cccb429b0de65eb3866f36039f615688ac29da96d9e12da69edabf97d77bd62537b7ed25202ba195360d56e3c8021109df6646f0f1cab6a6e130801a"))
	r1 := p1.Verify(p1Leaf, p1Root, false)
	fmt.Printf("result p1: %v\n", r1)
	fmt.Printf("p1: %+v leaf: %x root: %x\n", p1, p1Leaf, p1Root)

	p2Leaf := common.Hex2Bytes("b2c958036da87b5e289cc04cd8fb40341f78f43485445687d74be69395d11dc7")
	p2Root := common.Hex2Bytes("a35a6b6586809dfc915a23b3b289c67b05f23a82fee5a5f080c56ddb9103ef75")
	p2, _ = ToProof(uint64(0x69eb463bc4f6b2df), common.Hex2Bytes("0000000000000000"))
	r2 := p2.Verify(p2Leaf, p2Root, false)

	fmt.Printf("p2: %+v leaf: %x root: %x\n", p2, p2Leaf, p2Root)
	fmt.Printf("result p2: %v\n", r2)
}
