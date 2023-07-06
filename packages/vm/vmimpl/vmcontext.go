package vmimpl

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/buffered"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blob"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/migrations"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/execution"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"github.com/iotaledger/wasp/packages/vm/vmtxbuilder"
)

// vmContext represents state of the chain during one run of the VM while processing
// a batch of requests. vmContext object mutates with each request in the batch.
// The vmContext is created from immutable vm.VMTask object and UTXO state of the
// chain address contained in the statetxbuilder.Builder
type vmContext struct {
	task       *vm.VMTask
	taskResult *vm.VMTaskResult

	chainOwnerID    isc.AgentID
	blockContext    map[isc.Hname]interface{}
	txbuilder       *vmtxbuilder.AnchorTransactionBuilder
	chainInfo       *isc.ChainInfo
	blockGas        blockGas
	reqctx          *requestContext
	anchorOutputSD  uint64
	maintenanceMode bool

	currentStateUpdate *buffered.Mutations
	callStack          []*callContext
}

type blockGas struct {
	burned     uint64
	feeCharged uint64
	// is gas burn enabled
	// TODO: should be in requestGas?
	burnEnabled bool
}

type requestContext struct {
	req               isc.Request
	numPostedOutputs  int
	requestIndex      uint16
	requestEventIndex uint16
	entropy           hashing.HashValue
	evmFailed         *evmFailed
	gas               requestGas
	// SD charged to consume the current request
	sdCharged uint64
	// requests that the sender asked to retry
	unprocessableToRetry []isc.OnLedgerRequest
}

type evmFailed struct {
	tx      *types.Transaction
	receipt *types.Receipt
}

type requestGas struct {
	// max tokens that can be charged for gas fee
	maxTokensToSpendForGasFee uint64
	// final gas budget set for the run
	budgetAdjusted uint64
	// gas already burned
	burned uint64
	// tokens charged
	feeCharged uint64
	// burn history. If disabled, it is nil
	burnLog *gas.BurnLog
}

var _ execution.WaspContext = &vmContext{}

type callContext struct {
	caller             isc.AgentID // calling agent
	contract           isc.Hname   // called contract
	params             isc.Params  // params passed
	allowanceAvailable *isc.Assets // MUTABLE: allowance budget left after TransferAllowedFunds
}

// createVMContext creates a context for the whole batch run
func createVMContext(task *vm.VMTask, taskResult *vm.VMTaskResult) *vmContext {
	// assert consistency. It is a bit redundant double check
	if len(task.Requests) == 0 {
		// should never happen
		panic(errors.New("CreateVMContext.invalid params: must be at least 1 request"))
	}
	prevL1Commitment, err := transaction.L1CommitmentFromAliasOutput(task.AnchorOutput)
	if err != nil {
		// should never happen
		panic(fmt.Errorf("CreateVMContext: can't parse state data as L1Commitment from chain input %w", err))
	}

	taskResult.StateDraft, err = task.Store.NewStateDraft(task.TimeAssumption, prevL1Commitment)
	if err != nil {
		// should never happen
		panic(err)
	}

	vmctx := &vmContext{
		task:            task,
		taskResult:      taskResult,
		blockContext:    make(map[isc.Hname]interface{}),
		maintenanceMode: governance.NewStateAccess(taskResult.StateDraft).MaintenanceStatus(),
	}
	// at the beginning of each block
	l1Commitment, err := transaction.L1CommitmentFromAliasOutput(task.AnchorOutput)
	if err != nil {
		// should never happen
		panic(err)
	}

	var totalL2Funds *isc.Assets
	vmctx.withStateUpdate(func() {
		vmctx.runMigrations(migrations.BaseSchemaVersion, migrations.Migrations)

		// save the anchor tx ID of the current state
		vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
			blocklog.UpdateLatestBlockInfo(
				s,
				vmctx.task.AnchorOutputID.TransactionID(),
				isc.NewAliasOutputWithID(vmctx.task.AnchorOutput, vmctx.task.AnchorOutputID),
				l1Commitment,
			)
		})
		// get the total L2 funds in accounting
		totalL2Funds = vmctx.loadTotalFungibleTokens()
	})

	vmctx.anchorOutputSD = task.AnchorOutput.Amount - totalL2Funds.BaseTokens

	vmctx.txbuilder = vmtxbuilder.NewAnchorTransactionBuilder(
		task.AnchorOutput,
		task.AnchorOutputID,
		vmctx.anchorOutputSD,
		vmtxbuilder.AccountsContractRead{
			NativeTokenOutput:   vmctx.loadNativeTokenOutput,
			FoundryOutput:       vmctx.loadFoundry,
			NFTOutput:           vmctx.loadNFT,
			TotalFungibleTokens: vmctx.loadTotalFungibleTokens,
		},
	)

	return vmctx
}

func (vmctx *vmContext) withStateUpdate(f func()) {
	if vmctx.currentStateUpdate != nil {
		panic("expected currentStateUpdate == nil")
	}
	defer func() { vmctx.currentStateUpdate = nil }()

	vmctx.currentStateUpdate = buffered.NewMutations()
	f()
	vmctx.currentStateUpdate.ApplyTo(vmctx.taskResult.StateDraft)
}

// extractBlock does the closing actions on the block
// return nil for normal block and rotation address for rotation block
func (vmctx *vmContext) extractBlock(
	numRequests, numSuccess, numOffLedger uint16,
	unprocessable []isc.OnLedgerRequest,
) (uint32, *state.L1Commitment, time.Time, iotago.Address) {
	vmctx.GasBurnEnable(false)
	var rotationAddr iotago.Address
	vmctx.withStateUpdate(func() {
		rotationAddr = vmctx.saveBlockInfo(numRequests, numSuccess, numOffLedger)
		vmctx.closeBlockContexts()
		vmctx.saveInternalUTXOs(unprocessable)
	})

	block := vmctx.task.Store.ExtractBlock(vmctx.taskResult.StateDraft)

	l1Commitment := block.L1Commitment()

	blockIndex := vmctx.taskResult.StateDraft.BlockIndex()
	timestamp := vmctx.taskResult.StateDraft.Timestamp()

	return blockIndex, l1Commitment, timestamp, rotationAddr
}

func (vmctx *vmContext) checkRotationAddress() (ret iotago.Address) {
	vmctx.callCore(governance.Contract, func(s kv.KVStore) {
		ret = governance.GetRotationAddress(s)
	})
	return
}

// saveBlockInfo is in the blocklog partition context. Returns rotation address if this block is a rotation block
func (vmctx *vmContext) saveBlockInfo(numRequests, numSuccess, numOffLedger uint16) iotago.Address {
	if rotationAddress := vmctx.checkRotationAddress(); rotationAddress != nil {
		// block was marked fake by the governance contract because it is a committee rotation.
		// There was only on request in the block
		// We skip saving block information in order to avoid inconsistencies
		return rotationAddress
	}

	blockInfo := &blocklog.BlockInfo{
		SchemaVersion:         blocklog.BlockInfoLatestSchemaVersion,
		Timestamp:             vmctx.taskResult.StateDraft.Timestamp(),
		TotalRequests:         numRequests,
		NumSuccessfulRequests: numSuccess,
		NumOffLedgerRequests:  numOffLedger,
		PreviousAliasOutput:   isc.NewAliasOutputWithID(vmctx.task.AnchorOutput, vmctx.task.AnchorOutputID),
		GasBurned:             vmctx.blockGas.burned,
		GasFeeCharged:         vmctx.blockGas.feeCharged,
	}

	vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		blocklog.SaveNextBlockInfo(s, blockInfo)
		blocklog.Prune(s, blockInfo.BlockIndex(), vmctx.chainInfo.BlockKeepAmount)
	})
	vmctx.task.Log.Debugf("saved blockinfo: %s", blockInfo)
	return nil
}

// openBlockContexts calls the block context open function for all subscribed core contracts
func (vmctx *vmContext) openBlockContexts() {
	if vmctx.blockGas.burnEnabled {
		panic("expected gasBurnEnabled == false")
	}

	vmctx.withStateUpdate(func() {
		vmctx.loadChainConfig()

		var subs []root.BlockContextSubscription
		vmctx.callCore(root.Contract, func(s kv.KVStore) {
			subs = root.GetBlockContextSubscriptions(s)
		})
		for _, sub := range subs {
			vmctx.callProgram(sub.Contract, sub.OpenFunc, nil, nil, &isc.NilAgentID{})
		}
	})
}

// closeBlockContexts closes block contexts in deterministic FIFO sequence
func (vmctx *vmContext) closeBlockContexts() {
	if vmctx.blockGas.burnEnabled {
		panic("expected gasBurnEnabled == false")
	}
	var subs []root.BlockContextSubscription
	vmctx.callCore(root.Contract, func(s kv.KVStore) {
		subs = root.GetBlockContextSubscriptions(s)
	})
	for i := len(subs) - 1; i >= 0; i-- {
		vmctx.callProgram(subs[i].Contract, subs[i].CloseFunc, nil, nil, &isc.NilAgentID{})
	}
}

// saveInternalUTXOs relies on the order of the outputs in the anchor tx. If that order changes, this will be broken.
// Anchor Transaction outputs order must be:
// 0. Anchor Output
// 1. NativeTokens
// 2. Foundries
// 3. NFTs
// 4. produced outputs
// 5. unprocessable requests
func (vmctx *vmContext) saveInternalUTXOs(unprocessable []isc.OnLedgerRequest) {
	// create a mock AO, with a nil statecommitment, just to calculate changes in the minimum SD
	mockAO := vmctx.txbuilder.CreateAnchorOutput(vmctx.StateMetadata(state.L1CommitmentNil))
	newMinSD := parameters.L1().Protocol.RentStructure.MinRent(mockAO)
	oldMinSD := vmctx.anchorOutputSD
	changeInSD := int64(oldMinSD) - int64(newMinSD)

	if changeInSD != 0 {
		vmctx.task.Log.Debugf("adjusting commonAccount because AO SD cost changed, old:%d new:%d", oldMinSD, newMinSD)
		// update the commonAccount with the change in SD cost
		vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
			accounts.AdjustAccountBaseTokens(s, accounts.CommonAccount(), changeInSD)
		})
	}

	nativeTokenIDsToBeUpdated, nativeTokensToBeRemoved := vmctx.txbuilder.NativeTokenRecordsToBeUpdated()
	// IMPORTANT: do not iterate by this map, order of the slice above must be respected
	nativeTokensMap := vmctx.txbuilder.NativeTokenOutputsByTokenIDs(nativeTokenIDsToBeUpdated)

	foundryIDsToBeUpdated, foundriesToBeRemoved := vmctx.txbuilder.FoundriesToBeUpdated()
	// IMPORTANT: do not iterate by this map, order of the slice above must be respected
	foundryOutputsMap := vmctx.txbuilder.FoundryOutputsBySN(foundryIDsToBeUpdated)

	NFTOutputsToBeAdded, NFTOutputsToBeRemoved := vmctx.txbuilder.NFTOutputsToBeUpdated()

	blockIndex := vmctx.task.AnchorOutput.StateIndex + 1
	outputIndex := uint16(1)

	vmctx.callCore(accounts.Contract, func(s kv.KVStore) {
		// update native token outputs
		for _, ntID := range nativeTokenIDsToBeUpdated {
			vmctx.task.Log.Debugf("saving NT %s, outputIndex: %d", ntID, outputIndex)
			accounts.SaveNativeTokenOutput(s, nativeTokensMap[ntID], blockIndex, outputIndex)
			outputIndex++
		}
		for _, id := range nativeTokensToBeRemoved {
			vmctx.task.Log.Debugf("deleting NT %s", id)
			accounts.DeleteNativeTokenOutput(s, id)
		}

		// update foundry UTXOs
		for _, foundryID := range foundryIDsToBeUpdated {
			vmctx.task.Log.Debugf("saving foundry %d, outputIndex: %d", foundryID, outputIndex)
			accounts.SaveFoundryOutput(s, foundryOutputsMap[foundryID], blockIndex, outputIndex)
			outputIndex++
		}
		for _, sn := range foundriesToBeRemoved {
			vmctx.task.Log.Debugf("deleting foundry %d", sn)
			accounts.DeleteFoundryOutput(s, sn)
		}

		// update NFT Outputs
		for _, out := range NFTOutputsToBeAdded {
			vmctx.task.Log.Debugf("saving NFT %s, outputIndex: %d", out.NFTID, outputIndex)
			accounts.SaveNFTOutput(s, out, blockIndex, outputIndex)
			outputIndex++
		}
		for _, out := range NFTOutputsToBeRemoved {
			vmctx.task.Log.Debugf("deleting NFT %s", out.NFTID)
			accounts.DeleteNFTOutput(s, out.NFTID)
		}
	})

	// add unprocessable requests
	vmctx.storeUnprocessable(unprocessable, outputIndex)
}

func (vmctx *vmContext) removeUnprocessable(reqID isc.RequestID) {
	vmctx.withStateUpdate(func() {
		vmctx.callCore(blocklog.Contract, func(s kv.KVStore) {
			blocklog.RemoveUnprocessable(s, reqID)
		})
	})
}

func (vmctx *vmContext) assertConsistentGasTotals() {
	var sumGasBurned, sumGasFeeCharged uint64

	for _, r := range vmctx.taskResult.RequestResults {
		sumGasBurned += r.Receipt.GasBurned
		sumGasFeeCharged += r.Receipt.GasFeeCharged
	}
	if vmctx.blockGas.burned != sumGasBurned {
		panic("vmctx.gasBurnedTotal != sumGasBurned")
	}
	if vmctx.blockGas.feeCharged != sumGasFeeCharged {
		panic("vmctx.gasFeeChargedTotal != sumGasFeeCharged")
	}
}

func (vmctx *vmContext) LocateProgram(programHash hashing.HashValue) (vmtype string, binary []byte, err error) {
	vmctx.callCore(blob.Contract, func(s kv.KVStore) {
		vmtype, binary, err = blob.LocateProgram(s, programHash)
	})
	return vmtype, binary, err
}

func (vmctx *vmContext) Processors() *processors.Cache {
	return vmctx.task.Processors
}
