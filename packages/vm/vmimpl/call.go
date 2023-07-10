package vmimpl

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/execution"
	"github.com/iotaledger/wasp/packages/vm/sandbox"
)

// Call implements sandbox logic of the call between contracts on-chain
func (vmctx *vmContext) Call(targetContract, epCode isc.Hname, params dict.Dict, allowance *isc.Assets) dict.Dict {
	vmctx.Debugf("Call: targetContract: %s entry point: %s", targetContract, epCode)
	return vmctx.callProgram(targetContract, epCode, params, allowance)
}

func (vmctx *vmContext) callProgram(targetContract, epCode isc.Hname, params dict.Dict, allowance *isc.Assets, caller ...isc.AgentID) dict.Dict {
	contractRecord := vmctx.getOrCreateContractRecord(targetContract)
	ep := execution.GetEntryPointByProgHash(vmctx, targetContract, epCode, contractRecord.ProgramHash)

	vmctx.pushCallContext(targetContract, params, allowance, caller...)
	defer vmctx.popCallContext()

	// distinguishing between two types of entry points. Passing different types of sandboxes
	if ep.IsView() {
		return ep.Call(sandbox.NewSandboxView(vmctx))
	}
	// prevent calling 'init' not from root contract
	if epCode == isc.EntryPointInit {
		if !vmctx.callerIsRoot() {
			panic(fmt.Errorf("%v: target=(%s, %s)",
				vm.ErrRepeatingInitCall, vmctx.reqctx.req.CallTarget().Contract, epCode))
		}
	}
	return ep.Call(NewSandbox(vmctx))
}

func (vmctx *vmContext) callerIsRoot() bool {
	caller, ok := vmctx.Caller().(*isc.ContractAgentID)
	if !ok {
		return false
	}
	if !caller.ChainID().Equals(vmctx.ChainID()) {
		return false
	}
	return caller.Hname() == root.Contract.Hname()
}

const traceStack = false

func (vmctx *vmContext) pushCallContext(contract isc.Hname, params dict.Dict, allowance *isc.Assets, optCaller ...isc.AgentID) {
	var caller isc.AgentID
	if len(optCaller) != 0 {
		caller = optCaller[0]
	} else {
		caller = vmctx.getToBeCaller()
	}
	ctx := &callContext{
		caller:   caller,
		contract: contract,
		params: isc.Params{
			Dict:      params,
			KVDecoder: kvdecoder.New(params, vmctx.task.Log),
		},
		allowanceAvailable: allowance.Clone(), // we have to clone it because it will be mutated by TransferAllowedFunds
	}
	if traceStack {
		vmctx.Debugf("+++++++++++ PUSH %d, stack depth = %d caller = %s", contract, len(vmctx.callStack), ctx.caller)
	}
	vmctx.callStack = append(vmctx.callStack, ctx)
}

func (vmctx *vmContext) popCallContext() {
	if traceStack {
		vmctx.Debugf("+++++++++++ POP @ depth %d", len(vmctx.callStack))
	}
	vmctx.callStack[len(vmctx.callStack)-1] = nil // for GC
	vmctx.callStack = vmctx.callStack[:len(vmctx.callStack)-1]
}

func (vmctx *vmContext) getToBeCaller() isc.AgentID {
	if len(vmctx.callStack) > 0 {
		return vmctx.MyAgentID()
	}
	if vmctx.reqctx == nil {
		// e.g. saving the anchor ID
		return vmctx.ChainOwnerID()
	}
	return vmctx.reqctx.req.SenderAccount()
}

func (vmctx *vmContext) getCallContext() *callContext {
	if len(vmctx.callStack) == 0 {
		panic("getCallContext: stack is empty")
	}
	return vmctx.callStack[len(vmctx.callStack)-1]
}

func (vmctx *vmContext) callCore(c *coreutil.ContractInfo, f func(s kv.KVStore)) {
	vmctx.pushCallContext(c.Hname(), nil, nil)
	defer vmctx.popCallContext()

	f(vmctx.State())
}