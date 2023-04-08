// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/coregovernance"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
)

func setupGovernance(t *testing.T) *wasmsolo.SoloContext {
	ctx := setup(t)
	ctx = ctx.SoloContextForCore(t, coregovernance.ScName, coregovernance.OnDispatch)
	require.NoError(t, ctx.Err)
	return ctx
}

func TestRotateStateController(t *testing.T) {
	t.Skip("Chain.runRequestsNolock() hasn't been implemented yet")
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	user := ctx.NewSoloAgent("user")
	fadd := coregovernance.ScFuncs.AddAllowedStateControllerAddress(ctx)
	fadd.Params.Address().SetValue(user.ScAgentID().Address())
	fadd.Func.Post()
	require.NoError(t, ctx.Err)

	frot := coregovernance.ScFuncs.RotateStateController(ctx)
	frot.Params.Address().SetValue(user.ScAgentID().Address())
	frot.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestAddAllowedStateControllerAddress(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	user := ctx.NewSoloAgent("user")
	f := coregovernance.ScFuncs.AddAllowedStateControllerAddress(ctx)
	f.Params.Address().SetValue(user.ScAgentID().Address())
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestRemoveAllowedStateControllerAddress(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	user := ctx.NewSoloAgent("user")
	f := coregovernance.ScFuncs.RemoveAllowedStateControllerAddress(ctx)
	f.Params.Address().SetValue(user.ScAgentID().Address())
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestClaimChainOwnership(t *testing.T) {
	t.SkipNow()
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	user := ctx.NewSoloAgent("user")
	fdele := coregovernance.ScFuncs.DelegateChainOwnership(ctx)
	fdele.Params.ChainOwner().SetValue(user.ScAgentID())
	fdele.Func.Post()
	require.NoError(t, ctx.Err)

	fclaim := coregovernance.ScFuncs.ClaimChainOwnership(ctx)
	fclaim.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestDelegateChainOwnership(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	user := ctx.NewSoloAgent("user")
	f := coregovernance.ScFuncs.DelegateChainOwnership(ctx)
	f.Params.ChainOwner().SetValue(user.ScAgentID())
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestSetEVMGasRatioAndGetEVMGasRatio(t *testing.T) {
	t.Skip()

	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	gasRatio := "1:2"
	fSet := coregovernance.ScFuncs.SetEVMGasRatio(ctx)
	r, err := util.Ratio32FromString(gasRatio)
	require.NoError(t, err)
	fSet.Params.GasRatio().SetValue(r.Bytes())
	fSet.Func.Post()
	require.NoError(t, ctx.Err)

	fGet := coregovernance.ScFuncs.GetEVMGasRatio(ctx)
	fGet.Func.Call()
	require.NoError(t, ctx.Err)
	require.Equal(t, gasRatio, fGet.Results.GasRatio().String())
}

func TestSetFeePolicy(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	gfp0 := gas.DefaultFeePolicy()
	gfp0.GasPerToken = util.Ratio32{A: 1, B: 10}
	f := coregovernance.ScFuncs.SetFeePolicy(ctx)
	f.Params.FeePolicy().SetValue(gfp0.Bytes())
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestSetGasLimitsAndGetGasLimits(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	gasLimit := &gas.Limits{
		MaxGasPerBlock:         9_000_000_000,
		MinGasPerRequest:       1_000,
		MaxGasPerRequest:       5_000_000,
		MaxGasExternalViewCall: 5_000_000,
	}

	fSet := coregovernance.ScFuncs.SetGasLimits(ctx)
	fSet.Params.GasLimits().SetValue(gasLimit.Bytes())
	fSet.Func.Post()
	require.NoError(t, ctx.Err)

	fGet := coregovernance.ScFuncs.GetGasLimits(ctx)
	fGet.Func.Call()
	require.NoError(t, ctx.Err)
	retGasLimitBytes := fGet.Results.GasLimits().Value()
	retGasLimit, err := gas.LimitsFromBytes(retGasLimitBytes)
	require.NoError(t, err)
	require.Equal(t, gasLimit, retGasLimit)
}

func TestAddCandidateNode(t *testing.T) {
	t.SkipNow()
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.AddCandidateNode(ctx)
	f.Params.PubKey().SetValue(nil)
	f.Params.Certificate().SetValue(nil)
	f.Params.AccessAPI().SetValue("")
	f.Params.AccessOnly().SetValue(false)
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestRevokeAccessNode(t *testing.T) {
	t.SkipNow()
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.RevokeAccessNode(ctx)
	f.Params.PubKey().SetValue(nil)
	f.Params.Certificate().SetValue(nil)
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestSetMaintenanceStatusAndSetOnAndOff(t *testing.T) {
	t.Skip()

	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	fOn := coregovernance.ScFuncs.SetMaintenanceOn(ctx)
	fOn.Func.Post()
	require.NoError(t, ctx.Err)

	fStatus := coregovernance.ScFuncs.GetMaintenanceStatus(ctx)
	fStatus.Func.Call()
	require.NoError(t, ctx.Err)
	status := fStatus.Results.Status().Value()
	require.True(t, status)

	fOff := coregovernance.ScFuncs.SetMaintenanceOff(ctx)
	fOff.Func.Post()
	require.NoError(t, ctx.Err)

	fStatus.Func.Call()
	require.NoError(t, ctx.Err)
	status = fStatus.Results.Status().Value()
	require.False(t, status)
}

func TestSetCustomMetadata(t *testing.T) {
	// TODO
}

func TestGetAllowedStateControllerAddresses(t *testing.T) {
	// TODO
}

func TestChangeAccessNodes(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.ChangeAccessNodes(ctx)
	f.Params.Actions()
	f.Func.Post()
	require.NoError(t, ctx.Err)
}

func TestGetChainOwner(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.GetChainOwner(ctx)
	f.Func.Call()
	require.NoError(t, ctx.Err)
	assert.Equal(t, ctx.ChainOwnerID(), f.Results.ChainOwner().Value())
}

func TestGetChainNodes(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	// TODO first set up nodes / candidates so we have something to test f.Results for
	f := coregovernance.ScFuncs.GetChainNodes(ctx)
	f.Func.Call()
	require.NoError(t, ctx.Err)
}

func TestGetCustomMetadata(t *testing.T) {
	// TODO
}

func TestGetFeePolicy(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.GetFeePolicy(ctx)
	f.Func.Call()
	require.NoError(t, ctx.Err)
	fpBin := f.Results.FeePolicy().Value()
	gfp, err := gas.FeePolicyFromBytes(fpBin)
	require.NoError(t, err)
	require.Equal(t, gas.DefaultGasPerToken, gfp.GasPerToken)
	require.Equal(t, uint8(0), gfp.ValidatorFeeShare) // default fee share is 0
}

func TestGetChainInfo(t *testing.T) {
	ctx := setupGovernance(t)
	require.NoError(t, ctx.Err)

	f := coregovernance.ScFuncs.GetChainInfo(ctx)
	f.Func.Call()
	require.NoError(t, ctx.Err)
	assert.Equal(t, ctx.ChainOwnerID().String(), f.Results.ChainOwnerID().Value().String())
	gfp, err := gas.FeePolicyFromBytes(f.Results.FeePolicy().Value())
	require.NoError(t, err)
	assert.Equal(t, ctx.Chain.GetGasFeePolicy(), gfp)
}
