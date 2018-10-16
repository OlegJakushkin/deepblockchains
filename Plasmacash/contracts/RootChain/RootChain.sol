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
 * @title  PlasmaCash-MVP by Wolk
 * @author Michael Chung (michael@wolk.com)
 */


pragma solidity ^0.4.25;

import "./Libraries/Transaction.sol";
import "./Libraries/PriorityQueue.sol";
import "./Libraries/SparseMerkle.sol";
import "./Libraries/SafeMath.sol";

contract RootChain {

    event Deposit(address _depositor, uint64 indexed _depositIndex, uint64 _denomination,  uint64 indexed _tokenID);
    event StartExit(address _exiter, uint64 indexed _depositIndex, uint64 _denomination, uint64 indexed _tokenID, uint256 indexed _timestamp);
    event PublishedBlock(bytes32 _rootHash, uint64 indexed _blknum, uint64 indexed _currentDepositIndex);
    event FinalizedExit(address _exiter, uint64 indexed _depositIndex, uint64 _denomination,  uint64 indexed _tokenID, uint256 indexed _timestamp);
    event Challenge(address _challenger, uint64 indexed _tokenID, uint256 indexed _timestamp);

    /* test-only Events */
    event ExitStarted(uint256 indexed _priority);
    event ExitCompleted(uint256 indexed _priority);
    event CurrtentExit(uint64 depID, uint64 tokenID, uint256 exitableTS);
    event ExitTime(uint256 exitableTS, uint256 cuurrentTS);


    using Transaction for bytes;
    using SafeMath for uint64;
    address public authority;
    uint64 public currentDepositIndex;
    uint64 public currentBlkNum;
    mapping(uint64 => bytes32) public childChain;
    mapping(uint64 => uint64) public depositIndex;
    mapping(uint64 => uint64) public depositBalance;
    mapping(uint64 => Exit) public exits;
    PriorityQueue exitsQueue;
    SparseMerkle  smt;

    struct Exit {
        uint64  prevBlk;
        uint64  exitBlk;
        address exitor;
        uint64 balance;
        uint exitableTS;
        uint bond;
    }

    modifier isAuthority() {
        require(msg.sender == authority);
        _;
    }

   constructor() public {
        authority = msg.sender;
        currentDepositIndex = 0;
        currentBlkNum = 0;
        exitsQueue = new PriorityQueue();
        smt = new SparseMerkle();
    }

    // @dev Allows Plasma chain operator to submit block root
    // @param blkRoot The root of a child chain block
    // @param blknum The child chain block number
    function submitBlock(bytes32 _blkRoot, uint64 _blknum) public isAuthority {
        /* test-only: mutable blockroot */
        childChain[_blknum] = _blkRoot;
        currentBlkNum = _blknum;
        //require(currentBlkNum + 1 == _blknum);
        //currentBlkNum += 1;
        emit PublishedBlock(_blkRoot, _blknum, currentDepositIndex);
    }

    // @dev Allows anyone to deposit eth into the Plasma chain, Reject tokendeposit for now.
    function deposit() public payable {
        require (msg.value < (2**64 - 2) ); //18.446744073709551615Eth
        uint64 depositAmount = uint64(msg.value % (2 ** 64)) ;
        uint64 tokenID = uint64(uint256(keccak256(abi.encodePacked(msg.sender, currentDepositIndex, depositAmount))) % (2 ** 64));
        require (depositAmount > 0 && depositBalance[tokenID] == 0);
        depositBalance[tokenID] = depositAmount;
        depositIndex[currentDepositIndex] = tokenID;
        emit Deposit(msg.sender, currentDepositIndex, depositAmount, tokenID);
        currentDepositIndex += 1;
    }

    // @dev Allows original owner to submit withdraw request [bond will be added in futre]
    function depositExit(uint64 _depositIndex) public {
        require(exits[tokenID].exitableTS == 0);
        uint64 tokenID = depositIndex[_depositIndex];
        uint64 denomination = depositBalance[tokenID];
        require(uint64(uint256(keccak256(abi.encodePacked(msg.sender, _depositIndex, denomination))) % (2 ** 64)) == tokenID);
        require(denomination > 0);
        Exit memory etx = Exit(0, 1, msg.sender, denomination, block.timestamp + 60, msg.value); //depositExit validation uniquely checked by exit.prevBlk
        addExitToQueue(tokenID, _depositIndex, etx);
        emit StartExit(msg.sender, _depositIndex, denomination, tokenID, block.timestamp);
    }

    // @ dev Takes in the transaction transfering ownership to the current owner and the proofs necesary to prove there inclusion
    function startExit(bytes prevTxBytes, bytes prevProof, uint64 prevBlk, bytes txBytes, bytes proof, uint64 blk) public {
        Transaction.PlasmaTx memory exitTx = txBytes.parseTx();
        uint64 tokenID = exitTx.TokenID;
        require(tokenID > 0 && exits[tokenID].exitableTS == 0);
        Transaction.PlasmaTx memory prevTx = prevTxBytes.parseTx();
        require(exitTx.Recipient == msg.sender, "unauth exit");
        require(exitTx.PrevOwner == prevTx.Recipient, "invalid signer");
        require(exitTx.TokenID == prevTx.TokenID, "tokenID mismatch");

        require(txBytes.verifyTx(), "exitTx sig failure");
        require(prevTxBytes.verifyTx(), "prevTx sig failure");
        require(prevBlk < blk, "Invlid Exit order");
        require(prevBlk == exitTx.PrevBlock, "potential challengeBetween");

        //checkMembership(leaf, root, tokenID, proof);
        require(smt.checkMembership(keccak256(prevTxBytes),childChain[prevBlk], tokenID, prevProof), "prevTx non member");
        require(smt.checkMembership(keccak256(txBytes),childChain[blk], tokenID, proof), "exitTx non member");

        //StartExit. bond(currently not required)
        Exit memory etx = Exit(prevBlk, blk, msg.sender, exitTx.Balance, block.timestamp + 60, msg.value);
        addExitToQueue(tokenID, exitTx.DepositIndex, etx);
        emit StartExit(msg.sender, exitTx.DepositIndex, exitTx.Denomination, tokenID, block.timestamp);

    }


    // @ dev Submit proof that the exiting uid has been spent on the child chain either between the prevTx and tx or after the exit has been triggered. proving that the owner of the exiting uid is illegitimate
    function challenge(bytes txBytes, bytes proof, uint64 blk) public {

        Transaction.PlasmaTx memory challengeTx = txBytes.parseTx();
        uint64 tokenID = challengeTx.TokenID;
        Exit memory etx = exits[tokenID];
        require(tokenID > 0 && etx.exitableTS > 0, "exit not exist!");

        if (blk > etx.exitBlk) {
          if (etx.prevBlk == 0) {
            require(challengeTx.PrevBlock != 0, "Deposit cannot be challenged by itself");
            require(challengeTx.PrevOwner == etx.exitor, "Invalid challengeAfter cond");
          }else{
            require(challengeTx.PrevOwner == etx.exitor && challengeTx.PrevBlock == etx.exitBlk, "Invalid challengeAfter cond");  // outer range double spend
          }

        }else if (etx.prevBlk < blk && blk < etx.exitBlk){
          require(challengeTx.Recipient != etx.exitor, "Invlid challengeBetween cond"); // any legitimate challengeBetween will invalidate transaction after blk
        }

        require(txBytes.verifyTx(), "sig failure");
        require(smt.checkMembership(keccak256(txBytes),childChain[blk], tokenID, proof), "non member");

        //valid challenge, etx removed from queue
        delete exits[tokenID];
        emit Challenge(msg.sender, tokenID, block.timestamp);
        //TODO: claim bond attached to exit
    }

    // @ dev challenge a pending exit by proving that a faulty transaction has been incorrectly included in the chain. Therefore any subsequent transaction of this token are considered invalid
    function challengeBefore(bytes txBytes, bytes proof, uint64 blk, bytes faultyTxBytes, bytes faultyProof, uint64 faultyBlk) public {

        Transaction.PlasmaTx memory lastValidTx = txBytes.parseTx();
        uint64 tokenID = lastValidTx.TokenID;
        Exit memory etx = exits[tokenID];
        require(tokenID > 0 && etx.exitableTS > 0, "exitTX not exist!");
        Transaction.PlasmaTx memory faultyTx = faultyTxBytes.parseTx();
        require(lastValidTx.TokenID == faultyTx.TokenID, "tokenID mismatch");
        require(lastValidTx.Recipient != faultyTx.PrevOwner, "Invlid challengeBefore cond");
        require(txBytes.verifyTx(), "lastValidTx sig failure");
        require(faultyTxBytes.verifyTx(), "faultyTx sig failure");
        require(blk == faultyTx.PrevBlock, "jump sequence");
        require(blk < faultyBlk && faultyBlk <= etx.prevBlk, "Invlid transition sequence");

        //checkMembership(leaf, root, tokenID, proof);
        require(smt.checkMembership(keccak256(txBytes),childChain[blk],lastValidTx.TokenID, proof), "lastTx non member");
        require(smt.checkMembership(keccak256(faultyTxBytes),childChain[faultyBlk],faultyTx.TokenID, faultyProof), "faultyTx non member");

        //valid challenge, etx removed from queue
        delete exits[tokenID];
        emit Challenge(msg.sender, lastValidTx.TokenID, block.timestamp);
        //TODO: initate new exit afterchallengeBefore
    }


    function finalizeExits() public {
        uint64 depID;
        uint64 tokenID;
        uint256 exitableTS;
        uint256 currenTS = block.timestamp;
        uint8 counter = 0;
        (depID, tokenID, exitableTS) = getNextExit();
        Exit memory currentExit;

        emit CurrtentExit(depID, tokenID, exitableTS);
        emit ExitTime(exitableTS, currenTS);


        while (exitableTS < currenTS && counter <= 20) {
            currentExit = exits[tokenID];

            //Guard against Invalid Exit cancelled by challenge;
            if (currentExit.exitor != 0x0) {
                uint64 denomination = depositBalance[tokenID];
                uint64 exitorBalance = currentExit.balance;
                uint64 operatorBalance = denomination.uint64Sub(exitorBalance);

                delete depositBalance[tokenID];
                delete depositIndex[depID];
                //TODO: refund bond attached to valid exit
                currentExit.exitor.transfer(exitorBalance);
                authority.transfer(operatorBalance);
                emit FinalizedExit(currentExit.exitor, depID, exitorBalance, tokenID, currenTS);
            }

            emit ExitCompleted(exitsQueue.getMin());
            exitsQueue.delMin();
            delete exits[tokenID];
            counter = counter+1;

            if (exitsQueue.currentSize() > 0) {
                (depID, tokenID, exitableTS) = getNextExit();
                emit CurrtentExit(depID, tokenID, exitableTS);
                emit ExitTime(exitableTS, currenTS);
            } else {
                return;
            }
        }
    }

    // @dev Priority queue is sorted by exitable_at | DepositIndex
	function addExitToQueue(uint64 tokenID, uint64 depID, Exit memory etx) private {
	    require(exits[tokenID].exitableTS == 0, "Existing etx found");
        uint256 priority = etx.exitableTS << 128 | depID;
        exitsQueue.insert(priority);
        exits[tokenID] = etx;
        emit ExitStarted(priority);
    }

    // @dev Recovering depositIndex,  tokenID, exitableTS from priority
     function getNextExit() public view returns (uint64 depID, uint64 tokenID, uint256 exitableTS) {
	   uint256 priority = exitsQueue.getMin();
	   depID = uint64(priority);
	   exitableTS = priority >> 128;
	   return (depID, depositIndex[depID], exitableTS);
	 }

	 // test only
	 function kill() public isAuthority {
	    selfdestruct(msg.sender);
	 }

   function () public payable {
      deposit();
   }
}
