// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package vmcontext

import (
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/sandbox"
)

type contractSandbox struct {
	sandbox.SandboxBase
}

func NewSandbox(vmctx *VMContext) isc.Sandbox {
	ret := &contractSandbox{}
	ret.Ctx = vmctx
	return ret
}

// Call calls an entry point of contract, passes parameters and funds
func (s *contractSandbox) Call(target, entryPoint isc.Hname, params dict.Dict, transfer *isc.Assets) dict.Dict {
	s.Ctx.GasBurn(gas.BurnCodeCallContract)
	return s.Ctx.Call(target, entryPoint, params, transfer)
}

// DeployContract deploys contract by the binary hash
// and calls "init" endpoint (constructor) with provided parameters
func (s *contractSandbox) DeployContract(programHash hashing.HashValue, name, description string, initParams dict.Dict) {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeDeployContract)
	s.Ctx.(*VMContext).DeployContract(programHash, name, description, initParams)
}

func (s *contractSandbox) Event(msg string) {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeEmitEventFixed)
	s.Log().Infof("event::%s -> '%s'", s.Ctx.(*VMContext).CurrentContractHname(), msg)
	s.Ctx.(*VMContext).MustSaveEvent(s.Ctx.(*VMContext).CurrentContractHname(), msg)
}

func (s *contractSandbox) GetEntropy() hashing.HashValue {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.(*VMContext).Entropy()
}

func (s *contractSandbox) AllowanceAvailable() *isc.Assets {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeGetAllowance)
	return s.Ctx.(*VMContext).AllowanceAvailable()
}

func (s *contractSandbox) TransferAllowedFunds(target isc.AgentID, transfer ...*isc.Assets) *isc.Assets {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeTransferAllowance)
	return s.Ctx.(*VMContext).TransferAllowedFunds(target, transfer...)
}

func (s *contractSandbox) TransferAllowedFundsForceCreateTarget(target isc.AgentID, transfer ...*isc.Assets) *isc.Assets {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeTransferAllowance)
	return s.Ctx.(*VMContext).TransferAllowedFunds(target, transfer...)
}

func (s *contractSandbox) Request() isc.Calldata {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.(*VMContext).Request()
}

func (s *contractSandbox) Send(par isc.RequestParameters) {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeSendL1Request, uint64(s.Ctx.(*VMContext).NumPostedOutputs))
	s.Ctx.(*VMContext).Send(par)
}

func (s *contractSandbox) EstimateRequiredStorageDeposit(par isc.RequestParameters) uint64 {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeEstimateStorageDepositCost)
	return s.Ctx.(*VMContext).EstimateRequiredStorageDeposit(par)
}

func (s *contractSandbox) State() kv.KVStore {
	return s.Ctx.(*VMContext).State()
}

func (s *contractSandbox) StateAnchor() *isc.StateAnchor {
	s.Ctx.(*VMContext).GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.(*VMContext).StateAnchor()
}

func (s *contractSandbox) RegisterError(messageFormat string) *isc.VMErrorTemplate {
	return s.Ctx.(*VMContext).RegisterError(messageFormat)
}

func (s *contractSandbox) EVMTracer() *isc.EVMTracer {
	return s.Ctx.(*VMContext).task.EVMTracer
}

// helper methods

func (s *contractSandbox) RequireCallerAnyOf(agentIDs []isc.AgentID) {
	ok := false
	for _, agentID := range agentIDs {
		if s.Caller().Equals(agentID) {
			ok = true
		}
	}
	if !ok {
		panic(vm.ErrUnauthorized)
	}
}

func (s *contractSandbox) RequireCaller(agentID isc.AgentID) {
	if !s.Caller().Equals(agentID) {
		panic(vm.ErrUnauthorized)
	}
}

func (s *contractSandbox) RequireCallerIsChainOwner() {
	s.RequireCaller(s.ChainOwnerID())
}

func (s *contractSandbox) Privileged() isc.Privileged {
	return s
}

// privileged methods:

func (s *contractSandbox) TryLoadContract(programHash hashing.HashValue) error {
	return s.Ctx.(*VMContext).TryLoadContract(programHash)
}

func (s *contractSandbox) CreateNewFoundry(scheme iotago.TokenScheme, metadata []byte) (uint32, uint64) {
	return s.Ctx.(*VMContext).CreateNewFoundry(scheme, metadata)
}

func (s *contractSandbox) DestroyFoundry(sn uint32) uint64 {
	return s.Ctx.(*VMContext).DestroyFoundry(sn)
}

func (s *contractSandbox) ModifyFoundrySupply(sn uint32, delta *big.Int) int64 {
	return s.Ctx.(*VMContext).ModifyFoundrySupply(sn, delta)
}

func (s *contractSandbox) SetBlockContext(bctx interface{}) {
	s.Ctx.(*VMContext).SetBlockContext(bctx)
}

func (s *contractSandbox) BlockContext() interface{} {
	return s.Ctx.(*VMContext).BlockContext()
}

func (s *contractSandbox) GasBurnEnable(enable bool) {
	s.Ctx.GasBurnEnable(enable)
}

func (s *contractSandbox) MustMoveBetweenAccounts(fromAgentID, toAgentID isc.AgentID, assets *isc.Assets) {
	s.Ctx.(*VMContext).mustMoveBetweenAccounts(fromAgentID, toAgentID, assets)
	s.checkRemainingTokens(fromAgentID)
}

func (s *contractSandbox) DebitFromAccount(agentID isc.AgentID, tokens *isc.Assets) {
	s.Ctx.(*VMContext).debitFromAccount(agentID, tokens)
	s.checkRemainingTokens(agentID)
}

func (s *contractSandbox) checkRemainingTokens(debitedAccount isc.AgentID) {
	// assert that remaining tokens in the sender's account are enough to pay for the gas budget
	if debitedAccount.Equals(s.Request().SenderAccount()) && !s.HasInAccount(
		debitedAccount,
		s.totalGasTokens(),
	) {
		panic(vm.ErrNotEnoughTokensLeftForGas)
	}
}

func (s *contractSandbox) CreditToAccount(agentID isc.AgentID, tokens *isc.Assets) {
	s.Ctx.(*VMContext).creditToAccount(agentID, tokens)
}

func (s *contractSandbox) totalGasTokens() *isc.Assets {
	if s.Ctx.(*VMContext).task.EstimateGasMode {
		return isc.NewEmptyAssets()
	}
	amount := s.Ctx.(*VMContext).gasMaxTokensToSpendForGasFee
	return isc.NewAssetsBaseTokens(amount)
}

func (s *contractSandbox) CallOnBehalfOf(caller isc.AgentID, target, entryPoint isc.Hname, params dict.Dict, transfer *isc.Assets) dict.Dict {
	s.Ctx.GasBurn(gas.BurnCodeCallContract)
	return s.Ctx.(*VMContext).CallOnBehalfOf(caller, target, entryPoint, params, transfer)
}

func (s *contractSandbox) RetryUnprocessable(req isc.Request, blockIndex uint32, outputIndex uint16) {
	s.Ctx.(*VMContext).RetryUnprocessable(req, blockIndex, outputIndex)
}
