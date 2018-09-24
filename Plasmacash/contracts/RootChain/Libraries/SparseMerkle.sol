/** Copyright 2018 Wolk Inc.
* This file is part of the Plasmacash library.
*
* The plasmacash library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The Plasmacash library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
*
* You should have received a copy of the GNU Lesser General Public License
* along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
*/


/**
 * @title  sparseMerkle
 * @author Michael Chung (michael@wolk.com)
 * @dev sparse merkle tree implementation.
 */

pragma solidity ^0.4.25;

 contract SparseMerkle {

     bytes32[65] defaultHashes;
     constructor() public {
         defaultHashes[0] = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;
         setDefaultHashes(0,64);
     }

     function setDefaultHashes(uint8 startIndex, uint8 endIndex) internal {
         require(defaultHashes[startIndex] != 0 && defaultHashes[endIndex] == 0);
         for (uint8 i = startIndex; i < endIndex ; i += 1) {
             defaultHashes[i+1] =  keccak256(abi.encodePacked(defaultHashes[i], defaultHashes[i]));
         }
     }

     function checkMembership(bytes32 leaf, bytes32 root, uint64 tokenID, bytes proof) public view returns (bool) {
         bytes32 computedHash = getRoot(leaf, tokenID, proof);
         return (computedHash == root);
     }

     function getRoot(bytes32 leaf, uint64 index, bytes proof) internal view returns (bytes32) {
         require((proof.length - 8) % 32 == 0 && proof.length <= 2056);
         bytes32 proofElement;
         bytes32 computedHash = leaf;
         uint16 p = 8;
         uint64 proofBits;
         assembly { proofBits := div(mload(add(proof, 32)), exp(256, 24))}

         for (uint8 d = 0; d <= 63; d++ ) {
             if (proofBits % 2 == 0) {
                 proofElement = defaultHashes[d];
             } else {
                 p +=32;
                 require(proof.length >= p);
                 assembly { proofElement := mload(add(proof, p)) }
             }
             if (index % 2 == 0) {
                 computedHash = keccak256(abi.encodePacked(computedHash, proofElement));
             }else{
                 computedHash = keccak256(abi.encodePacked(proofElement, computedHash));
             }
             proofBits = proofBits / 2;
             index = index / 2;
         }
         return computedHash;
     }
 }
