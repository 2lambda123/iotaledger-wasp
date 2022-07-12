// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
	"github.com/iotaledger/wasp/tools/gascalibration/storage/go/storage"
	"github.com/stretchr/testify/require"
)

func TestCallF(t *testing.T) {
	wasmlib.ConnectHost(nil)
	ctx := wasmsolo.NewSoloContext(t, storage.ScName, storage.OnLoad)
	require.NoError(t, ctx.Err)

	f := storage.ScFuncs.F(ctx)
	for i := uint32(1); i <= 100; i++ {
		n := i * 10
		f.Params.N().SetValue(n)
		f.Func.Post()
		require.NoError(t, ctx.Err)
		t.Logf("n = %d, gas = %d\n", n, ctx.Gas)
	}
}
