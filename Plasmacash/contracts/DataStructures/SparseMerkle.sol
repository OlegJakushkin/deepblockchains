pragma solidity ^0.4.20;

/**
 * @title  sparseMerkle
 * @author Michael Chung (michael@wolk.com)
 * @dev sparse merkle tree implementation.
 */

 contract SparseMerkle {

     bytes32[65] defaultHashes;
     constructor() public {
         defaultHashes[0] = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;
         setDefaultHashes(0,64);
     }

     function setDefaultHashes(uint8 startIndex, uint8 endIndex) internal {
         require(defaultHashes[startIndex] != 0 && defaultHashes[endIndex] == 0);
         for (uint8 i = startIndex; i < endIndex ; i += 1) {
             defaultHashes[i+1] =  keccak256(defaultHashes[i], defaultHashes[i]);
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
                 computedHash = keccak256(computedHash, proofElement);
             }else{
                 computedHash = keccak256(proofElement, computedHash);
             }
             proofBits = proofBits / 2;
             index = index / 2;
         }
         return computedHash;
     }
 }
