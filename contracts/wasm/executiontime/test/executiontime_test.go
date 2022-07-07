// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/iotaledger/wasp/contracts/wasm/executiontime/go/executiontime"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
	"github.com/stretchr/testify/require"
)

func TestCallF(t *testing.T) {
	wasmlib.ConnectHost(nil)
	ctx := wasmsolo.NewSoloContext(t, executiontime.ScName, executiontime.OnLoad)
	require.NoError(t, ctx.ContractExists(executiontime.ScName))

	f := executiontime.ScFuncs.F(ctx)
	for i := uint32(1000); i <= 60000; i += 1000 {
		f.Params.N().SetValue(i)
		f.Func.Post()
		require.NoError(t, ctx.Err)
		t.Logf("n = %d, gas = %d\n", i, ctx.Gas)
	}
}
