package main

import (
	"fmt"
	"io"
	"os"

	"github.com/wasmerio/wasmer-go/wasmer"
)

var i32 = []wasmer.ValueKind{wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32}

func functype(nrParams, nrResults int) *wasmer.FunctionType {
	params := wasmer.NewValueTypes(i32[:nrParams]...)
	results := wasmer.NewValueTypes(i32[:nrResults]...)
	return wasmer.NewFunctionType(params, results)
}

func main() {
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
	file, err := os.Open("pkg/helloworld_bg.wasm")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	wasmData, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	module, err := wasmer.NewModule(store, wasmData)
	if err != nil {
		panic(err)
	}
	instance, err := wasmer.NewInstance(module, linker)
	if err != nil {
		panic(err)
	}
	_ = instance
	fmt.Println("no errors")
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
