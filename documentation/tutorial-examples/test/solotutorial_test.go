// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/iotaledger/wasp/documentation/tutorial-examples/go/solotutorial"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
	"github.com/stretchr/testify/require"
)

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, solotutorial.ScName, solotutorial.OnDispatch)
	require.NoError(t, ctx.ContractExists(solotutorial.ScName))
}
