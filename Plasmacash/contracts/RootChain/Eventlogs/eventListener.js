var fs = require('fs');
var os = require("os");
var Web3 = require('web3');
var _ = require('underscore');

var plasma_enpointUrl = 'http://localhost:8502'
var rootchain_enpointUrl = 'ws://localhost:8545'
var rootcontract_addr = '0xa611dd32bb2cc893bc57693bfa423c52658367ca'
var abi = [{"anonymous":false,"inputs":[{"indexed":false,"name":"exitableTS","type":"uint256"},{"indexed":false,"name":"cuurrentTS","type":"uint256"}],"name":"ExitTime","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"exiter","type":"address"},{"indexed":true,"name":"depositIndex","type":"uint64"},{"indexed":false,"name":"denomination","type":"uint64"},{"indexed":true,"name":"tokenID","type":"uint64"},{"indexed":true,"name":"timestamp","type":"uint256"}],"name":"FinalizedExit","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"rootHash","type":"bytes32"},{"indexed":true,"name":"blknum","type":"uint64"},{"indexed":true,"name":"currentDepositIndex","type":"uint64"}],"name":"PublishedBlock","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"priority","type":"uint256"}],"name":"ExitStarted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"exiter","type":"address"},{"indexed":true,"name":"depositIndex","type":"uint64"},{"indexed":false,"name":"denomination","type":"uint64"},{"indexed":true,"name":"tokenID","type":"uint64"},{"indexed":true,"name":"timestamp","type":"uint256"}],"name":"StartExit","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"depositor","type":"address"},{"indexed":true,"name":"depositIndex","type":"uint64"},{"indexed":false,"name":"denomination","type":"uint64"},{"indexed":true,"name":"tokenID","type":"uint64"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"priority","type":"uint256"}],"name":"ExitCompleted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"challenger","type":"address"},{"indexed":true,"name":"tokenID","type":"uint64"},{"indexed":true,"name":"timestamp","type":"uint256"}],"name":"Challenge","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"depID","type":"uint64"},{"indexed":false,"name":"tokenID","type":"uint64"},{"indexed":false,"name":"exitableTS","type":"uint256"}],"name":"CurrtentExit","type":"event"}];

//Connect to Layer1 - Rootchain
var web3 = new Web3();
web3.setProvider(new Web3.providers.WebsocketProvider(rootchain_enpointUrl));

//Connect to Layer2 - Plasma
var web3_plasma = new Web3();
web3_plasma.setProvider(new Web3.providers.HttpProvider(plasma_enpointUrl));
web3_plasma.extend({
		property: 'plasma',
		methods: [{
				name: 'eventHandler',
        call: 'plasma_eventHandler',
        params: 1,
        inputFormatter: [ null ]
    }]
	});

var contract = new web3.eth.Contract(abi, rootcontract_addr);

function writeLog(filename) {
	var logger = fs.createWriteStream(filename, {
  		flags: 'a' // 'a' means appending (old data will be preserved)
	});

	var subscription = web3.eth.subscribe('logs', {
	    address: rootcontract_addr,
	}, function(error, result){
	    if (error)
	      console.log(result);
	})
	.on("data", function(log){
		var eventJsonInterface = _.find(
			contract._jsonInterface,
			o => o.signature === log.topics[0] && o.type === 'event',
		)
		var decodedevent = web3.eth.abi.decodeLog(eventJsonInterface.inputs, log.data, log.topics.slice(1));
		decodedevent = _.omit(decodedevent,_.range(decodedevent.__length__),'__length__');

		var eventLog = {'event' : eventJsonInterface.name};
		_.each(decodedevent, function(value, key) {
			 decodedevent[key] = web3.utils.toHex(value);
		});

		eventLog['log'] = decodedevent;
		logger.write(JSON.stringify(eventLog) + os.EOL);
		console.log("EventHandler call : " + JSON.stringify(eventLog))
    //TODO: sign her prior forwarding to eventHandler
		var eventHandler_resp = web3_plasma.plasma.eventHandler(eventLog)
		//TODO: handle eventHandler_resp
	})
	.on("changed", function(log){
		var eventJsonInterface = _.find(
			contract._jsonInterface,
			o => o.signature === log.topics[0] && o.type === 'event',
		)
		var decodedevent = web3.eth.abi.decodeLog(eventJsonInterface.inputs, log.data, log.topics.slice(1));
		decodedevent = _.omit(decodedevent,_.range(decodedevent.__length__),'__length__');

		var eventLog = {'event' : eventJsonInterface.name};
		_.each(decodedevent, function(value, key) {
			decodedevent[key] = web3.utils.toHex(value);
		});
		eventLog['log'] = decodedevent;
		eventLog['removed'] = true;
		console.log(JSON.stringify(eventLog));
		logger.write(JSON.stringify(eventLog) + os.EOL);
	});
}

module.exports = {
    writeLog
}

require('make-runnable');
