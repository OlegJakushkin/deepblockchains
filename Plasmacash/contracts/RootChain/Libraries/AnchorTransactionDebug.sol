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
 * @title Deep Blockchains - Anchor Transaction Debug
 * @author Michael Chung (michael@wolk.com)
 * @dev Anchor Transaction Debug Tool. Helper functions are not used in production contract
 */

pragma solidity ^0.4.25;
//pragma experimental ABIEncoderV2;

import './AnchorTransaction.sol';

contract AnchorTransactionDebug {

    using AnchorTransaction for bytes;
    using RLP for bytes;
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLPEncode for bytes[];
    using RLPEncode for bytes;

    function verifyTx(bytes txBytes, address signer) public pure returns (bool isValidTx) {
        return txBytes.getSigner() == signer;
    }

    function getSigner(bytes txBytes) public pure returns (address signer) {
        return txBytes.getSigner();
    }

    function getUnsignedtxBytes(bytes memory txBytes) public pure returns (bytes unsignedTxByte) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        bytes[] memory unsignedTx = new bytes[](5);
        for(uint i=0; i<rlpTx.length; i++) {
            if (i==4){
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        return unsignedTx.encodeList();
    }

    function getShortHash(bytes memory txBytes) public pure returns (bytes32 shortHash) {
        bytes memory rlpUnsignedTx = getUnsignedtxBytes(txBytes);
        return keccak256(rlpUnsignedTx);
    }

    function getSignHash(bytes memory txBytes) public pure returns (bytes32 signedHash) {
        bytes32 shortHash = getShortHash(txBytes);
        bytes memory prefix = "\x19Ethereum Signed Message:\n32";
        return keccak256(abi.encodePacked(prefix, shortHash));
    }

    function getTxHash(bytes memory txBytes) public pure returns (bytes32 txHash) {
        return keccak256(txBytes);
    }

    function getBlockchainID(bytes memory txBytes) public pure returns (uint64 chaindID) {
        return txBytes.parseTx().BlockChainID;
    }

    function getBlockNumber(bytes memory txBytes) public pure returns (uint64 blockNumber) {
        return txBytes.parseTx().BlockNumber;
    }

    function getBlockHash(bytes memory txBytes) public pure returns (bytes32 blockHash) {
        return txBytes.parseTx().BlockHash;
    }

    function getAddedOwners(bytes memory txBytes) public pure returns (address[] addedOwners) {
        return txBytes.parseTx().Extra.AddedOwners;
    }

    function getRemovedOwners(bytes memory txBytes) public pure returns (address[] removedOwners) {
        return txBytes.parseTx().Extra.RemovedOwners;
    }

    function getSig(bytes memory txBytes) public pure returns (bytes signiture) {
        return txBytes.parseTx().Sig;
    }

    // Experimental Debug Curretnly Not Available
    // function parseAnchorTx(bytes memory txBytes) public pure returns (AnchorTransaction.AnchorTx memory anchorTx) {
    //     return txBytes.parseTx();
    // }
}
