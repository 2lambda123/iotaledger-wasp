package vmimpl

import (
	"github.com/ethereum/go-ethereum/common"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
)

func EstimateRequiredStorageDeposit(senderAddress iotago.Address, contractHName isc.Hname, par isc.RequestParameters) uint64 {
	contractIdentity := isc.ContractIdentityFromHname(contractHName)
	if contractHName == evm.Contract.Hname() {
		contractIdentity = isc.ContractIdentityFromEVMAddress(common.Address{}) // use empty EVM address as STUB
	}

	out := transaction.BasicOutputFromPostData(
		senderAddress,
		contractIdentity,
		par,
	)

	return parameters.L1().Protocol.RentStructure.MinRent(out)
}

func (reqctx *requestContext) estimateRequiredStorageDeposit(par isc.RequestParameters) uint64 {
	par.AdjustToMinimumStorageDeposit = false

	estimation := EstimateRequiredStorageDeposit(
		reqctx.vm.task.AnchorOutput.AliasID.ToAddress(),
		reqctx.CurrentContractHname(),
		par,
	)

	return estimation
}
