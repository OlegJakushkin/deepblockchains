// Copyright 2018 Wolk Inc.
// This file is part of the Wolk Plasma library.
package eventlog

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//{"depositor":"0xa45b77a98e2b840617e2ec6ddfbf71403bdcb683","depositIndex":"0x19","denomination":"0xde0b6b3a7640000","tokenID":"0xe6a1a18a43b04212"}
type DepositEvent struct {
	Depositor    common.Address `json:"depositor"     gencodec:"required"`
	DepositIndex uint64         `json:"depositIndex"  gencodec:"required"`
	Denomination uint64         `json:"denomination"  gencodec:"required"`
	TokenID      uint64         `json:"tokenID"       gencodec:"required"`
}

//go:generate gencodec -type DepositEvent -field-override depositMarshaling -out deposit_json.go
type depositMarshaling struct {
	DepositIndex hexutil.Uint64
	Denomination hexutil.Uint64
	TokenID      hexutil.Uint64
}

//{"exiter":"0xa45b77a98e2b840617e2ec6ddfbf71403bdcb683","depositIndex":"0x0","denomination":"0xde0b6b3a7640000","tokenID":"0xb437230feb2d24db","timestamp":"0x5bb54fc1"}
type StartExitEvent struct {
	Exiter       common.Address `json:"exiter"         gencodec:"required"`
	DepositIndex uint64         `json:"depositIndex"   gencodec:"required"`
	Denomination uint64         `json:"denomination"   gencodec:"required"`
	TokenID      uint64         `json:"tokenID"        gencodec:"required"`
	TS           uint64         `json:"timestamp"      gencodec:"required"`
}

//go:generate gencodec -type StartExitEvent -field-override startExitMarshaling -out startExit_json.go
type startExitMarshaling struct {
	DepositIndex hexutil.Uint64
	Denomination hexutil.Uint64
	TokenID      hexutil.Uint64
	TS           hexutil.Uint64
}

//{"rootHash":"0x82da88c31e874c678d529ad51e43de3a4baf3914","currentDepositIndex":"0x15","blkNum":"0x15"}
type PublishedBlockEvent struct {
	RootHash            common.Hash `json:"rootHash"             gencodec:"required"`
	Blocknumber         uint64      `json:"blknum"               gencodec:"required"`
	CurrentDepositIndex uint64      `json:"currentDepositIndex"  gencodec:"required"`
}

//go:generate gencodec -type PublishedBlockEvent -field-override publishedBlockMarshaling -out publishedBlock_json.go
type publishedBlockMarshaling struct {
	Blocknumber         hexutil.Uint64
	CurrentDepositIndex hexutil.Uint64
}

//{"challenger":"0xbef06cc63c8f81128c26efedd461a9124298092b","tokenID":"0x9af84bc1208918b","timestamp":"0x5bb54fc2"}
type ChallengeEvent struct {
	Challenger common.Address `json:"challenger" gencodec:"required"`
	TokenID    uint64         `json:"tokenID"    gencodec:"required"`
	TS         uint64         `json:"timestamp"  gencodec:"required"`
}

//go:generate gencodec -type ChallengeEvent -field-override challengeMarshaling -out challenge_json.go
type challengeMarshaling struct {
	TokenID hexutil.Uint64
	TS      hexutil.Uint64
}

//{"exiter":"0x74f978a3e049688777e6120d293f24348bde5fa6","depositIndex":"0x4","denomination":"0x3782dace9d900000","tokenID":"0x7c00dfa72e8832ed","timestamp":"0x5bb54ffe"}
type FinalizedExitEvent struct {
	Exiter       common.Address `json:"exiter"         gencodec:"required"`
	DepositIndex uint64         `json:"depositIndex"   gencodec:"required"`
	Denomination uint64         `json:"denomination"   gencodec:"required"`
	TokenID      uint64         `json:"tokenID"        gencodec:"required"`
	TS           uint64         `json:"timestamp"      gencodec:"required"`
}

//go:generate gencodec -type FinalizedExitEvent -field-override finalizedExitEventMarshaling -out finalizedExit_json.go
type finalizedExitEventMarshaling struct {
	DepositIndex hexutil.Uint64
	Denomination hexutil.Uint64
	TokenID      hexutil.Uint64
	TS           hexutil.Uint64
}
