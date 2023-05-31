package vmcontext

import (
	"math"
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/util/panicutil"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/vmcontext/vmexceptions"
)

// creditToAccount deposits transfer from request to chain account of of the called contract
// It adds new tokens to the chain ledger. It is used when new tokens arrive with a request
func (vmctx *VMContext) creditToAccount(agentID isc.AgentID, ftokens *isc.Assets) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.CreditToAccount(s, agentID, ftokens)
	})
}

func (vmctx *VMContext) creditNFTToAccount(agentID isc.AgentID, nft *isc.NFT) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.CreditNFTToAccount(s, agentID, nft)
	})
}

// debitFromAccount subtracts tokens from account if it is enough of it.
// should be called only when posting request
func (vmctx *VMContext) debitFromAccount(agentID isc.AgentID, transfer *isc.Assets) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.DebitFromAccount(s, agentID, transfer)
	})
}

// debitNFTFromAccount removes a NFT from account.
// should be called only when posting request
func (vmctx *VMContext) debitNFTFromAccount(agentID isc.AgentID, nftID iotago.NFTID) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.DebitNFTFromAccount(s, agentID, nftID)
	})
}

func (vmctx *VMContext) mustMoveBetweenAccounts(fromAgentID, toAgentID isc.AgentID, assets *isc.Assets) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.MustMoveBetweenAccounts(s, fromAgentID, toAgentID, assets)
	})
}

func (vmctx *VMContext) findContractByHname(contractHname isc.Hname) (ret *root.ContractRecord) {
	vmctx.callCore(root.Contract, func(s kv.KVStore) {
		ret = root.FindContract(s, contractHname)
	})
	return ret
}

func (vmctx *VMContext) getChainInfo() *isc.ChainInfo {
	var ret *isc.ChainInfo
	vmctx.callCore(governance.Contract, func(s kv.KVStore) {
		ret = governance.MustGetChainInfo(s, vmctx.ChainID())
	})
	return ret
}

func (vmctx *VMContext) GetBaseTokensBalance(agentID isc.AgentID) uint64 {
	var ret uint64
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.GetBaseTokensBalance(s, agentID)
	})
	return ret
}

func (vmctx *VMContext) HasEnoughForAllowance(agentID isc.AgentID, allowance *isc.Assets) bool {
	var ret bool
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.HasEnoughForAllowance(s, agentID, allowance)
	})
	return ret
}

func (vmctx *VMContext) GetNativeTokenBalance(agentID isc.AgentID, nativeTokenID iotago.NativeTokenID) *big.Int {
	var ret *big.Int
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.GetNativeTokenBalance(s, agentID, nativeTokenID)
	})
	return ret
}

func (vmctx *VMContext) GetNativeTokenBalanceTotal(nativeTokenID iotago.NativeTokenID) *big.Int {
	var ret *big.Int
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.GetNativeTokenBalanceTotal(s, nativeTokenID)
	})
	return ret
}

func (vmctx *VMContext) GetNativeTokens(agentID isc.AgentID) iotago.NativeTokens {
	var ret iotago.NativeTokens
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.GetNativeTokens(s, agentID)
	})
	return ret
}

func (vmctx *VMContext) GetAccountNFTs(agentID isc.AgentID) (ret []iotago.NFTID) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.GetAccountNFTs(s, agentID)
	})
	return ret
}

func (vmctx *VMContext) GetNFTData(nftID iotago.NFTID) (ret *isc.NFT) {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		ret = accounts.MustGetNFTData(s, nftID)
	})
	return ret
}

func (vmctx *VMContext) GetSenderTokenBalanceForFees() uint64 {
	sender := vmctx.req.SenderAccount()
	if sender == nil {
		return 0
	}
	return vmctx.GetBaseTokensBalance(sender)
}

func (vmctx *VMContext) requestLookupKey() blocklog.RequestLookupKey {
	return blocklog.NewRequestLookupKey(vmctx.task.StateDraft.BlockIndex(), vmctx.requestIndex)
}

func (vmctx *VMContext) eventLookupKey() blocklog.EventLookupKey {
	return blocklog.NewEventLookupKey(vmctx.task.StateDraft.BlockIndex(), vmctx.requestIndex, vmctx.requestEventIndex)
}

func (vmctx *VMContext) writeReceiptToBlockLog(vmError *isc.VMError) *blocklog.RequestReceipt {
	receipt := &blocklog.RequestReceipt{
		Request:       vmctx.req,
		GasBudget:     vmctx.gasBudgetAdjusted,
		GasBurned:     vmctx.gasBurned,
		GasFeeCharged: vmctx.gasFeeCharged,
		SDCharged:     vmctx.sdCharged,
	}

	if vmError != nil {
		b := vmError.Bytes()
		if len(b) > isc.VMErrorMessageLimit {
			vmError = coreerrors.ErrErrorMessageTooLong
		}
		receipt.Error = vmError.AsUnresolvedError()
	}

	vmctx.Debugf("writeReceiptToBlockLog - reqID:%s err: %v", vmctx.req.ID(), vmError)

	receipt.GasBurnLog = vmctx.gasBurnLog

	if vmctx.task.EnableGasBurnLogging {
		vmctx.gasBurnLog = gas.NewGasBurnLog()
	}
	key := vmctx.requestLookupKey()
	var err error
	vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		err = blocklog.SaveRequestReceipt(s, receipt, key)
	})
	if err != nil {
		panic(err)
	}
	return receipt
}

func (vmctx *VMContext) storeUnprocessable(lastInternalAssetUTXOIndex uint16) {
	if len(vmctx.unprocessable) == 0 {
		return
	}
	blockIndex := vmctx.task.AnchorOutput.StateIndex + 1

	vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		for _, r := range vmctx.unprocessable {
			txsnapshot := vmctx.createTxBuilderSnapshot()
			err := panicutil.CatchPanic(func() {
				position := vmctx.txbuilder.ConsumeUnprocessable(r)
				outputIndex := position + int(lastInternalAssetUTXOIndex)
				if blocklog.HasUnprocessable(s, r.ID()) {
					panic("already in unprocessable list")
				}
				// save the unprocessable requests and respective output indices onto the state so they can be retried later
				blocklog.SaveUnprocessable(s, r, blockIndex, uint16(outputIndex))
			})
			if err != nil {
				// protocol exception triggered. Rollback
				vmctx.restoreTxBuilderSnapshot(txsnapshot)
			}
		}
	})
}

func (vmctx *VMContext) MustSaveEvent(hContract isc.Hname, topic string, payload []byte) {
	if vmctx.requestEventIndex == math.MaxUint16 {
		panic(vm.ErrTooManyEvents)
	}
	vmctx.Debugf("MustSaveEvent/%s: topic: '%s'", hContract.String(), topic)

	event := &isc.Event{
		ContractID: hContract,
		Topic:      topic,
		Payload:    payload,
		Timestamp:  uint64(vmctx.Timestamp().UnixNano()),
	}
	eventKey := vmctx.eventLookupKey().Bytes()
	vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		blocklog.SaveEvent(s, eventKey, event)
	})
	vmctx.requestEventIndex++
}

// updateOffLedgerRequestNonce updates stored nonce for off ledger requests
func (vmctx *VMContext) updateOffLedgerRequestNonce() {
	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		accounts.IncrementNonce(s, vmctx.req.SenderAccount())
	})
}

// adjustL2BaseTokensIfNeeded adjust L2 ledger for base tokens if the L1 changed because of storage deposit changes
func (vmctx *VMContext) adjustL2BaseTokensIfNeeded(adjustment int64, account isc.AgentID) {
	if adjustment == 0 {
		return
	}
	err := panicutil.CatchPanicReturnError(func() {
		vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
			accounts.AdjustAccountBaseTokens(s, account, adjustment)
		})
	}, accounts.ErrNotEnoughFunds)
	if err != nil {
		panic(vmexceptions.ErrNotEnoughFundsForInternalStorageDeposit)
	}
}
