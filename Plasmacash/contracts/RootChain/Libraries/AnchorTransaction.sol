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
 * @title Deep Blockchains - Anchor Transaction Library (with Metamask support)
 * @author Michael Chung (michael@wolk.com)
 * @dev This library handles AnchorTx and Extra struct RLP encode/decode on chain
 */

pragma solidity ^0.4.25;

import './RLP.sol';
import './RLPEncode.sol';
import './ECRecovery.sol';
import './SafeMath.sol';

library AnchorTransaction {

    using SafeMath for uint64;
    using RLP for bytes;
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLPEncode for bytes[];
    using RLPEncode for bytes;

    struct AnchorTx {
        uint64  BlockChainID;
        uint64  BlockNumber;
        bytes32 BlockHash;
        Ownership Extra;
        bytes Sig;
    }

    struct Ownership {
        address[]  AddedOwners;
        address[]  RemovedOwners;
    }

    function getSigner(bytes memory txBytes) internal pure returns (address) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        bytes[] memory unsignedTx = new bytes[](5);
        bytes memory sig;

        for(uint i=0; i<rlpTx.length; i++) {
            if (i==4){
                sig = rlpTx[i].toData();
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        bytes memory rlpUnsignedTx =  unsignedTx.encodeList();
        return ECRecovery.recover(keccak256(rlpUnsignedTx), sig);
    }

    function parseOwnership(bytes memory extraByte) internal pure returns (Ownership memory) {
        RLP.RLPItem[] memory rlpExtra = extraByte.toRLPItem().toList(2);
        Ownership memory extra;
        extra.AddedOwners = parseAddrs(rlpExtra[0].toBytes());
        extra.RemovedOwners = parseAddrs(rlpExtra[1].toBytes());
        return extra;
    }

    function parseAddrs(bytes memory addrsBytes) internal pure returns (address[] memory addrList) {
        RLP.RLPItem memory rlpAddrList = addrsBytes.toRLPItem();
        require(rlpAddrList.isList());
        uint length = rlpAddrList.items();
        addrList = new address[](length);

        RLP.RLPItem[] memory rlpAddrs = addrsBytes.toRLPItem().toList(length);
        for(uint i=0; i<length; i++) {
            addrList[i] = rlpAddrs[i].toAddress();
        }

        return addrList;
    }

    function parseTx(bytes memory txBytes) internal pure returns (AnchorTx memory) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        AnchorTx memory txn;
        txn.BlockChainID =  uint64(rlpTx[0].toUint());
        txn.BlockNumber =  uint64(rlpTx[1].toUint());
        txn.BlockHash =  rlpTx[2].toBytes32();
        txn.Extra =  parseOwnership(rlpTx[3].toBytes());
        txn.Sig =  rlpTx[4].toData();
        return txn;
    }
}
