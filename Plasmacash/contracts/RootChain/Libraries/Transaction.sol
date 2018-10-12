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
 * @title Deep Blockchains - Plasma Transaction Library (with Metamask support)
 * @author Michael Chung (michael@wolk.com)
 * @dev This library handles PlasmaTx RLP encode/decode and verifies signature on chain
 */

pragma solidity ^0.4.25;

import './RLP.sol';
import './RLPEncode.sol';
import './ECRecovery.sol';
import './SafeMath.sol';

library Transaction {

    using SafeMath for uint64;
    using RLP for bytes;
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLPEncode for bytes[];
    using RLPEncode for bytes;

    struct PlasmaTx {
        uint64  TokenID;
        uint64  Denomination;
        uint64  DepositIndex;
        uint64  PrevBlock;
        address PrevOwner;
        address Recipient;
        uint64  Allowance;
        uint64  Spent;
        //bytes Sig;
        uint64  Balance;
    }

    function getSigner(bytes memory txBytes) internal pure returns (address) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        bytes[] memory unsignedTx = new bytes[](9);
        bytes memory sig;

        for(uint i=0; i<rlpTx.length; i++) {
            if (i==4){
                unsignedTx[i] = rlpTx[i].toBytes();
            }else if (i==8){
                sig = rlpTx[i].toData();
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        bytes memory rlpUnsignedTx =  unsignedTx.encodeList();
        address signer = ECRecovery.recover(keccak256(rlpUnsignedTx), sig);
        return signer;
     }

    function verifyTx(bytes memory txBytes) internal pure returns (bool) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        if ((rlpTx[3].toUint() == 0)) {
            //If prevBlock = 0, check whether last 8 bytes of keccak256(recipient, depositIndex, denomination) == tokenID
            return uint256(keccak256(abi.encodePacked(rlpTx[5].toAddress(),uint64(rlpTx[2].toUint()),uint64(rlpTx[1].toUint())))) % (2**64) == (rlpTx[0].toUint());
        }
        if (uint64(rlpTx[6].toUint()).uint64Add(uint64(rlpTx[7].toUint())) > uint64(rlpTx[1].toUint())) {
            return false;
        }

        bytes[] memory unsignedTx = new bytes[](9);
        bytes memory sig;
        address prevOwner;
        for(uint i=0; i<9; i++) {
            if (i==4){
                prevOwner = rlpTx[i].toAddress();
                unsignedTx[i] = rlpTx[i].toBytes();
            }else if (i==8){
                sig = rlpTx[i].toData();
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        bytes memory rlpUnsignedTx = unsignedTx.encodeList();
        return (ECRecovery.recover(keccak256(rlpUnsignedTx), sig) == prevOwner);
     }

    function parseTx(bytes memory txBytes) internal pure returns (PlasmaTx memory) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        PlasmaTx memory txn;
        txn.TokenID =  uint64(rlpTx[0].toUint());
        txn.Denomination =  uint64(rlpTx[1].toUint());
        txn.DepositIndex =  uint64(rlpTx[2].toUint());
        txn.PrevBlock =  uint64(rlpTx[3].toUint());
        txn.PrevOwner =  rlpTx[4].toAddress();
        txn.Recipient =  rlpTx[5].toAddress();
        txn.Allowance =  uint64(rlpTx[6].toUint());
        txn.Spent =  uint64(rlpTx[7].toUint());
        txn.Balance = uint64(rlpTx[1].toUint()).uint64Sub(uint64(rlpTx[6].toUint())).uint64Sub(uint64(rlpTx[7].toUint()));
        return txn;
    }
}
