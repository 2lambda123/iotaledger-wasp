// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/iotaledger/wasp/contracts/wasm/storage/go/storage"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
	"github.com/stretchr/testify/require"
)

func TestCallF(t *testing.T) {
	wasmlib.ConnectHost(nil)
	ctx := wasmsolo.NewSoloContext(t, storage.ScName, storage.OnLoad)
	require.NoError(t, ctx.ContractExists(storage.ScName))

	f := storage.ScFuncs.F(ctx)
	for i := uint32(100); i <= 5000; i += 100 {
		f.Params.N().SetValue(i)
		f.Func.Post()
		require.NoError(t, ctx.Err)
		t.Logf("n = %d, gas = %d\n", i, ctx.Gas)
	}
}
