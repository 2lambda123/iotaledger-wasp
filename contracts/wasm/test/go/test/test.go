// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"strconv"

	"github.com/iotaledger/wasp/packages/vm/wasmlib/go/wasmlib"
)

func funcAddTestStruct(ctx wasmlib.ScFuncContext, f *AddTestStructContext) {
	testStructs := f.State.TestStructs()
	length := testStructs.Length()
	description := f.Params.Description().Value()
	testStruct := &TestStruct{
		Id:          length,
		Description: description,
	}
	testStructs.GetTestStruct(length).SetValue(testStruct)

	ctx.Event("teststruct.Added " + strconv.Itoa(int(length)) + " " + description)
}

func funcClearAll(ctx wasmlib.ScFuncContext, f *ClearAllContext) {
	f.State.TestStructs().Clear()
	ctx.Log("teststructs.Cleared " + strconv.Itoa(int(f.State.TestStructs().Length())))

	for i := int32(1); i < 10; i++ {
		description := f.State.TestStructs().GetTestStruct(i).Value().Description
		id := f.State.TestStructs().GetTestStruct(i).Value().Id
		ctx.Log("teststructs.Cleared.Log " + strconv.Itoa(int(id)) + " " + description)
	}
}
