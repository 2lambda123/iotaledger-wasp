// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmhost

import (
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type WasmerVM struct {
	WasmVMBase
	engine     *wasmer.Engine
	instance   *wasmer.Instance
	linker     *wasmer.ImportObject
	memory     *wasmer.Memory
	module     *wasmer.Module
	store      *wasmer.Store
	lastBudget uint64
	instances  uint32
}

var _ WasmVM = new(WasmerVM)

var wasmerI32Params = []wasmer.ValueKind{wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32}

func NewWasmerVM() WasmVM {
	vm := &WasmerVM{}
	opmap := map[wasmer.Opcode]uint32{
		// TODO Add gas fees
	}
	config := wasmer.NewConfig().PushMeteringMiddleware(gas.MaxGasPerRequest, opmap)
	engine := wasmer.NewEngineWithConfig(config)
	vm.store = wasmer.NewStore(engine)
	return vm
}

func (vm *WasmerVM) GasBudget(budget uint64) {
	// save budget so we can later determine how much the VM burned
	vm.lastBudget = budget

	// new budget for VM, top up to desired budget
	vm.instance.SetRemainingPoints(budget)

	// consume 0 fuel to determine remaining budget
	remainingBudget := vm.instance.GetRemainingPoints()
	if remainingBudget != budget {
		panic("GasBudget.determine: can't set fuel")
	}
}

// GasBurned will return the gas burned since the last time GasBudget() was called
func (vm *WasmerVM) GasBurned() uint64 {
	// consume 0 fuel to determine remaining budget
	remainingBudget := vm.instance.GetRemainingPoints()
	burned := vm.lastBudget - remainingBudget
	return burned
}

//TODO
func (vm *WasmerVM) Interrupt() {
	panic("implement me")
}

func (vm *WasmerVM) LinkHost() (err error) {
	vm.linker = wasmer.NewImportObject()

	funcs := map[string]wasmer.IntoExtern{
		FuncHostStateGet: vm.importFunc(4, 1, vm.exportHostStateGet),
		FuncHostStateSet: vm.importFunc(4, 0, vm.exportHostStateSet),
	}
	vm.linker.Register(ModuleWasmLib, funcs)

	funcs = map[string]wasmer.IntoExtern{
		FuncAbort: vm.importFunc(4, 0, vm.exportAbort),
	}
	vm.linker.Register(ModuleEnv, funcs)

	// TinyGo Wasm implementation uses this one to write panic message to console
	funcs = map[string]wasmer.IntoExtern{
		FuncFdWrite: vm.importFunc(4, 1, vm.exportFdWrite),
	}
	vm.linker.Register(ModuleWasi1, funcs)
	vm.linker.Register(ModuleWasi2, funcs)
	return nil
}

func (vm *WasmerVM) LoadWasm(wasmData []byte) error {
	var err error
	vm.module, err = wasmer.NewModule(vm.store, wasmData)
	if err != nil {
		return err
	}
	vm.instance, err = wasmer.NewInstance(vm.module, vm.linker)
	if err != nil {
		return err
	}
	vm.memory, err = vm.instance.Exports.GetMemory("memory")
	return err
}

func (vm *WasmerVM) NewInstance(wc *WasmContext) WasmVM {
	var err error
	vm.instances++
	if (vm.instances & 0xff) == 0 {
		err := vm.LinkHost()
		if err != nil {
			panic(err)
		}
	}
	vmInstance := &WasmerVM{
		engine: vm.engine,
		module: vm.module,
		linker: vm.linker,
		store:  vm.store,
	}
	vmInstance.wc = wc
	vmInstance.timeoutStarted = true // DisableWasmTimeout
	vmInstance.instance, err = wasmer.NewInstance(vm.module, vm.linker)
	if err != nil {
		panic(err)
	}
	vm.memory, err = vm.instance.Exports.GetMemory("memory")
	if err != nil {
		panic(err)
	}
	return vmInstance
}

func (vm *WasmerVM) RunFunction(functionName string, args ...interface{}) error {
	export, err := vm.instance.Exports.GetFunction(functionName)
	if err != nil {
		return err
	}
	return vm.Run(func() error {
		_, err = export(args...)
		return err
	})
}

func (vm *WasmerVM) RunScFunction(index int32) error {
	export, err := vm.instance.Exports.GetFunction("on_call")
	if err != nil {
		return err
	}
	err = vm.Run(func() error {
		_, err = export(index)
		return err
	})
	return err
}

func (vm *WasmerVM) UnsafeMemory() []byte {
	return vm.memory.Data()
}

func (vm *WasmerVM) exportAbort(args []wasmer.Value) ([]wasmer.Value, error) {
	errMsg := args[0].I32()
	fileName := args[1].I32()
	line := args[2].I32()
	col := args[3].I32()
	vm.HostAbort(errMsg, fileName, line, col)
	return nil, nil
}

func (vm *WasmerVM) exportFdWrite(args []wasmer.Value) ([]wasmer.Value, error) {
	fd := args[0].I32()
	iovs := args[1].I32()
	size := args[2].I32()
	written := args[3].I32()
	ret := vm.HostFdWrite(fd, iovs, size, written)
	return []wasmer.Value{wasmer.NewI32(ret)}, nil
}

func (vm *WasmerVM) exportHostStateGet(args []wasmer.Value) ([]wasmer.Value, error) {
	keyRef := args[0].I32()
	keyLen := args[1].I32()
	valRef := args[2].I32()
	valLen := args[3].I32()
	ret := vm.HostStateGet(keyRef, keyLen, valRef, valLen)
	return []wasmer.Value{wasmer.NewI32(ret)}, nil
}

func (vm *WasmerVM) exportHostStateSet(args []wasmer.Value) ([]wasmer.Value, error) {
	keyRef := args[0].I32()
	keyLen := args[1].I32()
	valRef := args[2].I32()
	valLen := args[3].I32()
	vm.HostStateSet(keyRef, keyLen, valRef, valLen)
	return nil, nil
}

func (vm *WasmerVM) importFunc(nrParams, nrResults int, function func([]wasmer.Value) ([]wasmer.Value, error)) *wasmer.Extern {
	params := wasmer.NewValueTypes(wasmerI32Params[:nrParams]...)
	results := wasmer.NewValueTypes(wasmerI32Params[:nrResults]...)
	funcType := wasmer.NewFunctionType(params, results)
	return wasmer.NewFunction(vm.store, funcType, function).IntoExtern()
}
