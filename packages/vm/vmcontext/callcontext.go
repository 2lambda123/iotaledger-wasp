package vmcontext

import (
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// pushCallContextAndMoveAssets moves funds from the caller to the target before pushing new context to the stack
func (vmctx *VMContext) pushCallContextAndMoveAssets(contract iscp.Hname, params dict.Dict, transfer *iscp.Assets) error {
	if transfer != nil {
		targetAccount := iscp.NewAgentID(vmctx.ChainID().AsAddress(), contract)
		targetAccount = vmctx.adjustAccount(targetAccount)
		var sourceAccount *iscp.AgentID
		if len(vmctx.callStack) == 0 {
			sourceAccount = vmctx.req.Request().SenderAccount()
		} else {
			sourceAccount = vmctx.AccountID()
		}
		vmctx.mustMoveBetweenAccounts(sourceAccount, targetAccount, transfer)
	}
	vmctx.pushCallContext(contract, params, transfer)
	return nil
}

const traceStack = false

func (vmctx *VMContext) pushCallContext(contract iscp.Hname, params dict.Dict, transfer *iscp.Assets) {
	if traceStack {
		vmctx.Log().Debugf("+++++++++++ PUSH %d, stack depth = %d", contract, len(vmctx.callStack))
	}
	var caller *iscp.AgentID
	if len(vmctx.callStack) == 0 {
		// request context
		caller = vmctx.req.Request().SenderAccount()
	} else {
		caller = vmctx.MyAgentID()
	}
	if traceStack {
		vmctx.Log().Debugf("+++++++++++ PUSH %d, stack depth = %d caller = %s", contract, len(vmctx.callStack), caller.String())
	}
	vmctx.callStack = append(vmctx.callStack, &callContext{
		caller:   caller,
		contract: contract,
		params:   params.Clone(),
		transfer: transfer,
	})
}

func (vmctx *VMContext) popCallContext() {
	if traceStack {
		vmctx.Log().Debugf("+++++++++++ POP @ depth %d", len(vmctx.callStack))
	}
	vmctx.callStack[len(vmctx.callStack)-1] = nil // for GC
	vmctx.callStack = vmctx.callStack[:len(vmctx.callStack)-1]
}

func (vmctx *VMContext) getCallContext() *callContext {
	if len(vmctx.callStack) == 0 {
		panic("getCallContext: stack is empty")
	}
	return vmctx.callStack[len(vmctx.callStack)-1]
}
