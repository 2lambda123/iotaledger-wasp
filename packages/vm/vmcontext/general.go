package vmcontext

import (
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
)

func (vmctx *VMContext) ChainID() *iscp.ChainID {
	return (*iscp.ChainID)(&vmctx.task.AnchorOutput.AliasID)
}

func (vmctx *VMContext) ChainOwnerID() *iscp.AgentID {
	return vmctx.chainOwnerID
}

func (vmctx *VMContext) ContractCreator() *iscp.AgentID {
	rec, ok := vmctx.findContractByHname(vmctx.CurrentContractHname())
	if !ok {
		panic("can't find current contract")
	}
	return rec.Creator
}

func (vmctx *VMContext) CurrentContractHname() iscp.Hname {
	return vmctx.getCallContext().contract
}

func (vmctx *VMContext) MyAgentID() *iscp.AgentID {
	return iscp.NewAgentID(vmctx.ChainID().AsAddress(), vmctx.CurrentContractHname())
}

func (vmctx *VMContext) Caller() *iscp.AgentID {
	return vmctx.getCallContext().caller
}

func (vmctx *VMContext) Timestamp() int64 {
	return vmctx.virtualState.Timestamp().UnixNano()
}

func (vmctx *VMContext) Entropy() hashing.HashValue {
	return vmctx.entropy
}

func (vmctx *VMContext) Request() iscp.Request {
	return vmctx.req.Request()
}

// TODO dust provision
func (vmctx *VMContext) Send(target iotago.Address, assets *iscp.Assets, metadata *iscp.SendMetadata, options ...*iscp.SendOptions) {
	if assets == nil {
		panic("post request assets can't be nil")
	}
	var sendOptions *iscp.SendOptions
	if len(options) > 0 {
		sendOptions = options[0]
	}
	// debit the assets from the on-chain account
	vmctx.debitFromAccount(vmctx.AccountID(), assets)

	// instruct tx builder about changed totals on-chain
	vmctx.txbuilder.SubDeltaIotas(assets.Iotas)
	bi := new(big.Int)
	for _, nt := range assets.Tokens {
		bi.Neg(nt.Amount)
		vmctx.txbuilder.AddDeltaNativeToken(nt.ID, bi)
	}

	vmctx.txbuilder.AddPostedRequest(iscp.PostRequestData{
		TargetAddress:  target,
		SenderContract: vmctx.CurrentContractHname(),
		Assets:         assets,
		Metadata:       metadata,
		SendOptions:    sendOptions,
	})
}

var _ iscp.StateAnchor = &VMContext{}

func (vmctx *VMContext) StateController() iotago.Address {
	return vmctx.task.AnchorOutput.StateController
}

func (vmctx *VMContext) GovernanceController() iotago.Address {
	return vmctx.task.AnchorOutput.GovernanceController
}

func (vmctx *VMContext) StateIndex() uint32 {
	return vmctx.task.AnchorOutput.StateIndex
}

func (vmctx *VMContext) StateHash() hashing.HashValue {
	h, err := hashing.HashValueFromBytes(vmctx.task.AnchorOutput.StateMetadata)
	if err != nil {
		panic(err)
	}
	return h
}

func (vmctx *VMContext) OutputID() iotago.UTXOInput {
	return vmctx.task.AnchorOutputID
}
