pragma solidity ^0.4.24;
// abigen --abi plasmacash.abi --pkg plasmachain --type PlasmaCash --out plasmacash.go

import './TxHash.sol';
import "./PriorityQueue.sol";
import "./SparseMerkle.sol";

/**
 * @title  PlasmaCash-MVP by Wolk
 * @author Michael Chung (michael@wolk.com)
 */

contract PlasmaCash {

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


    using TxHash for bytes;
    address public authority;
    uint64 public currentDepositIndex;
    uint64 public currentBlkNum;
    mapping(uint64 => bytes32) public childChain;
    mapping(uint64 => uint64) public depositIndex;
    mapping(uint64 => uint64) public depositBalance;
    mapping(uint64 => exitTX) public exits;
    PriorityQueue exitsQueue;
    SparseMerkle  smt;

    struct exitTX {
        uint64  txblk1;
        uint64  txblk2;
        address exitor;
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
        //require(currentBlkNum + 1 == _blknum);
        //currentBlkNum += 1;
        childChain[_blknum] = _blkRoot;
        currentBlkNum = _blknum;
        emit PublishedBlock(_blkRoot, _blknum, currentDepositIndex);
    }

    // @dev Allows anyone to deposit eth into the Plasma chain, Reject tokendeposit for now.
    function deposit() public payable {
        require (msg.value < (2**64 - 2) ); //18.446744073709551615Eth
        uint64 depositAmount = uint64(msg.value % (2 ** 64)) ;
        uint64 tokenID = uint64(uint256(smt.Keccak256(msg.sender, currentDepositIndex, depositAmount)) % (2 ** 64));
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
        require(uint64(uint256(smt.Keccak256(msg.sender, _depositIndex, denomination)) % (2 ** 64)) == tokenID);
        require(denomination > 0);
        exitTX memory exitx = exitTX(0, currentBlkNum, msg.sender, block.timestamp + 60, msg.value);
        addExitToQueue(tokenID, _depositIndex, exitx);
        emit StartExit(msg.sender, _depositIndex, denomination, tokenID, block.timestamp);
    }

    // @ dev Takes in the transaction transfering ownership to the current owner and the proofs necesary to prove there inclusion
    function startExit(uint64 tokenID, bytes txBytes1, bytes txBytes2, bytes proof1, bytes proof2, uint64 blk1, uint64 blk2) public {

        require(exits[tokenID].exitableTS == 0);
        TxHash.TX memory tx2 = txBytes2.getTx();
        TxHash.TX memory tx1 = txBytes1.getTx();
        require(tx2.Recipient == msg.sender, "unauth exit");
        require(tx1.TokenId == tokenID && tx2.TokenId == tokenID, "tokenID mismatch");
        require(txBytes1.verifyTX(), "tx1 sig failure");
        require(txBytes2.verifyTX(), "tx2 sig failure");
        require(blk1 < blk2, "Invlid Exit order");
        require(blk1 == tx2.PrevBlock, "potential challengeBetween");

        //checkMembership(leaf, root, tokenID, proof);
        require(smt.checkMembership(smt.Keccak256(txBytes1),childChain[blk1],tx1.TokenId, proof1), "tx1 non member");
        require(smt.checkMembership(smt.Keccak256(txBytes2),childChain[blk2],tx2.TokenId, proof2), "tx2 non member");

        //StartExit. bond(currently not required)
        exitTX memory exitx = exitTX(blk1, blk2, msg.sender, block.timestamp + 60, msg.value);
        addExitToQueue(tokenID, tx2.DepositIndex, exitx);
        emit StartExit(msg.sender, tx2.DepositIndex, tx2.Denomination, tokenID, block.timestamp);

    }


    // @ dev Submit proof that the exiting uid has been spent on the child chain either between the prevTx and tx or after the exit has been triggered. proving that the owner of the exiting uid is illegitimate
    function challenge(uint64 tokenID, bytes txBytes, bytes proof, uint64 blk) public {

        exitTX memory exitx = exits[tokenID];
        require((blk > exitx.txblk2) || (exitx.txblk1 < blk && blk < exitx.txblk2), "invalid txn bound");
        TxHash.TX memory tx = txBytes.getTx();
        require(tx.TokenId == tokenID, "tokenID mismatch");
        require(txBytes.verifyTX(), "sig failure");
        require(smt.checkMembership(smt.Keccak256(txBytes),childChain[blk],tx.TokenId, proof), "non member");

        //valid challenge, exitx removed from queue
        delete exits[tokenID];
        Challenge(msg.sender, tx.TokenId, block.timestamp);
        //TODO: claim bond attached to exit
    }


    function finalizeExits() public {
        uint64 depID;
        uint64 tokenID;
        uint256 exitableTS;
        uint256 currenTS = block.timestamp;
        (depID, tokenID, exitableTS) = getNextExit();
        exitTX memory currentExit;

        emit CurrtentExit(depID, tokenID, exitableTS);
        emit ExitTime(exitableTS, currenTS);

        //TODO: modify getNextExit to create_at
        while (exitableTS < currenTS) {
            currentExit = exits[tokenID];

            //Guard against Invalid Exit cancelled by challenge;
            if (currentExit.exitor != 0x0) {
                uint64 denomination = depositBalance[tokenID];
                delete depositBalance[tokenID];
                //TODO: refund bond attached to valid exit
                currentExit.exitor.transfer(denomination);
                delete depositIndex[depID];
                emit FinalizedExit(currentExit.exitor, depID, denomination, tokenID, currenTS);
            }

            emit ExitCompleted(exitsQueue.getMin());
            exitsQueue.delMin();
            delete exits[tokenID];

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
	function addExitToQueue(uint64 tokenID, uint64 depID, exitTX memory etx) private {
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
