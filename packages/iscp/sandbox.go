// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscp

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// SandboxBase is the common interface of Sandbox and SandboxView
type SandboxBase interface {
	// AccountID returns the agentID of the current contract
	AccountID() *AgentID
	// Params returns the parameters of the current call
	Params() dict.Dict
	// Balance returns number of iotas in the balance of the smart contract
	BalanceIotas() uint64
	// Assets returns all assets: iotas and native tokens
	Assets() *Assets
	// ChainID returns the chain ID
	ChainID() *ChainID
	// ChainOwnerID returns the AgentID of the current owner of the chain
	ChainOwnerID() *AgentID
	// Contract returns the Hname of the contract in the current chain
	Contract() Hname
	// ContractCreator returns the agentID that deployed the contract
	ContractCreator() *AgentID
	// GetTimestamp returns the timestamp of the current state
	GetTimestamp() int64
	// Log returns a logger that ouputs on the local machine. It includes Panicf method
	Log() LogInterface
	// Utils provides access to common necessary functionality
	Utils() Utils
	// Gas returns interface for gas related functions
	Gas() Gas
}

// Sandbox is an interface given to the processor to access the VMContext
// and virtual state, transaction builder and request parameters through it.
type Sandbox interface {
	SandboxBase

	// State k/v store of the current call (in the context of the smart contract)
	State() kv.KVStore
	// Request return the request in the context of which the smart contract is called
	Request() Request

	// Call calls the entry point of the contract with parameters and transfer.
	// If the entry point is full entry point, transfer tokens are moved between caller's and
	// target contract's accounts (if enough). If the entry point is view, 'transfer' has no effect
	Call(target, entryPoint Hname, params dict.Dict, transfer *Assets) (dict.Dict, error)
	// Caller is the agentID of the caller.
	Caller() *AgentID
	// DeployContract deploys contract on the same chain. 'initParams' are passed to the 'init' entry point
	DeployContract(programHash hashing.HashValue, name string, description string, initParams dict.Dict) error
	// Event publishes "vmmsg" message through Publisher on nanomsg. It also logs locally, but it is not the same thing
	Event(msg string)
	// GetEntropy 32 random bytes based on the hash of the current state transaction
	GetEntropy() hashing.HashValue // 32 bytes of deterministic and unpredictably random data
	// IncomingTransfer return colored balances transferred by the call. They are already accounted into the Balances()
	IncomingTransfer() *Assets
	// Send one generic method for sending assets with ledgerstate.ExtendedLockedOutput
	// replaces TransferToAddress and Post1Request
	Send(target iotago.Address, assets *Assets, metadata *SendMetadata, options ...*SendOptions)
	// Internal for use in native hardcoded contracts
	BlockContext(construct func(sandbox Sandbox) interface{}, onClose func(interface{})) interface{}
	// properties of the anchor output
	StateAnchor() StateAnchor
}

type Gas interface {
	Burn(uint64)
	Budget() uint64
}

// properties of the anchor output/transaction in the current context
type StateAnchor interface {
	StateController() iotago.Address
	GovernanceController() iotago.Address
	StateIndex() uint32
	StateHash() hashing.HashValue
	OutputID() iotago.UTXOInput
}

type SendOptions struct {
}

// RequestMetadata represents content of the data payload of the output
type SendMetadata struct {
	TargetContract Hname
	EntryPoint     Hname
	Args           dict.Dict
	Transfer       Assets
}

// PostRequestData is a parameters for a cross-chain request
type PostRequestData struct {
	TargetAddress  iotago.Address
	SenderContract Hname
	Assets         *Assets
	Metadata       *SendMetadata
	SendOptions    *SendOptions
	GasBudget      uint64
}
