// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/iotaledger/wasp/contracts/wasm/test/go/test"
	"github.com/iotaledger/wasp/packages/vm/wasmsolo"
	"github.com/stretchr/testify/require"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, test.ScName, test.OnLoad)
	require.NoError(t, ctx.ContractExists(test.ScName))
}

func TestClear(t *testing.T) {
	t.SkipNow()
	//*wasmsolo.GoDebug = true
	ctx := wasmsolo.NewSoloContext(t, test.ScName, test.OnLoad)
	require.NoError(t, ctx.Err)

	f := test.ScFuncs.AddTestStruct(ctx)
	f.Params.Description().SetValue("Eric")
	f.Func.TransferIotas(1).Post()
	require.NoError(t, ctx.Err)

	c := test.ScFuncs.ClearAll(ctx)
	c.Func.TransferIotas(1).Post()
	require.NoError(t, ctx.Err)

	f = test.ScFuncs.AddTestStruct(ctx)
	f.Params.Description().SetValue("Hop")
	f.Func.TransferIotas(1).Post()
	require.NoError(t, ctx.Err)
}

var i32 = []wasmer.ValueKind{wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32}

func functype(nrParams, nrResults int) *wasmer.FunctionType {
	params := wasmer.NewValueTypes(i32[:nrParams]...)
	results := wasmer.NewValueTypes(i32[:nrResults]...)
	return wasmer.NewFunctionType(params, results)
}

func TestWasmer(t *testing.T) {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	linker := wasmer.NewImportObject()
	funcs := map[string]wasmer.IntoExtern{
		"hostGetBytes":    wasmer.NewFunction(store, functype(5, 1), exportHostGetBytes).IntoExtern(),
		"hostGetKeyID":    wasmer.NewFunction(store, functype(2, 1), exportHostGetKeyID).IntoExtern(),
		"hostGetObjectID": wasmer.NewFunction(store, functype(3, 1), exportHostGetObjectID).IntoExtern(),
		"hostSetBytes":    wasmer.NewFunction(store, functype(5, 0), exportHostSetBytes).IntoExtern(),
	}
	linker.Register("WasmLib", funcs)
	file, err := os.Open("test_bg.wasm")
	require.NoError(t, err)
	defer file.Close()
	wasmData, err := io.ReadAll(file)
	require.NoError(t, err)
	module, err := wasmer.NewModule(store, wasmData)
	require.NoError(t, err)
	require.NotNil(t, module)
	_, err = wasmer.NewInstance(module, linker)
	require.NoError(t, err)
}

func exportHostGetBytes(args []wasmer.Value) ([]wasmer.Value, error) {
	fmt.Println("exportHostGetBytes")
	return nil, nil
}

func exportHostGetKeyID(args []wasmer.Value) ([]wasmer.Value, error) {
	fmt.Println("exportHostGetKeyID")
	return nil, nil
}

func exportHostGetObjectID(args []wasmer.Value) ([]wasmer.Value, error) {
	fmt.Println("exportHostGetObjectID")
	return nil, nil
}

func exportHostSetBytes(args []wasmer.Value) ([]wasmer.Value, error) {
	fmt.Println("exportHostSetBytes")
	return nil, nil
}
