# Example smart contract

This example native smart contract stores a string in the state upon request and allows to retrive the store string using a view.

The code below creates a new smart contract and initialize a processor with the required entry points and initializer method. The initializer method sets initial state.

```go
package helloworld

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// This is the contract object
var Contract = coreutil.NewContract("helloworld", "Say hello, a PoC contract")

// This is the contract processor
var Processor = Contract.Processor(initialize,
	FuncSetGreeting.WithHandler(setGreeting),
	FuncGetGreeting.WithHandler(getGreeting),
)

var (
	FuncSetGreeting = coreutil.Func("setGreeting")
	FuncGetGreeting = coreutil.ViewFunc("getGreeting")
)

const VarGreeting = "greeting"

// initialize smart contract
// sets the initial greeting to "Hello World"
func initialize(ctx iscp.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("helloworld.init in %s", ctx.Contract().String())
	greeting := "Hello World"
	ctx.State().Set(VarGreeting, codec.EncodeString(greeting))
	ctx.Event(fmt.Sprintf("helloworld.init.successs. greeting = %s", greeting))
	return nil, nil
}

// setGreeting is a full entry point. It sets the greeting state
func setGreeting(ctx iscp.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("helloworld.setGreeting in %s", ctx.Contract().String())
	params := ctx.Params()
	val, _, err := codec.DecodeString(params.MustGet(VarGreeting))
	if err != nil {
		return nil, fmt.Errorf("helloworld.setGreeting: %v", err)
	}
	ctx.State().Set(VarGreeting, codec.EncodeString(val))
	ctx.Event(fmt.Sprintf("helloworld.setGreeting.success. greeting = %s", val))
	return nil, nil
}

// getGreeting is a view entry point. It takes an iscp.SandboxView as parameter
func getGreeting(ctx iscp.SandboxView) (dict.Dict, error) {
	ctx.Log().Debugf("helloworld.getGreeting in %s", ctx.Contract().String())
	state := ctx.State()
	val, _, _ := codec.DecodeString(state.MustGet(VarGreeting))
	ret := dict.New()
	ret.Set(VarGreeting, codec.EncodeString(val))
	return ret, nil
}
```

To Learn how to interact with this smart contract from `solo` and `wasp-cli` proceed to the next tutorial.
