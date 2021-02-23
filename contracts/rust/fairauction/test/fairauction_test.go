// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/contracts/common"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

var auctioneer signaturescheme.SignatureScheme
var tokenColor balance.Color

func setupTest(t *testing.T) *solo.Chain {
	chain := common.StartChainAndDeployWasmContractByName(t, ScName)

	// set up auctioneer account and mint some tokens to auction off
	auctioneer = chain.Env.NewSignatureSchemeWithFunds()
	newColor, err := chain.Env.MintTokens(auctioneer, 10)
	require.NoError(t, err)
	chain.Env.AssertAddressBalance(auctioneer.Address(), balance.ColorIOTA, solo.Saldo-10)
	chain.Env.AssertAddressBalance(auctioneer.Address(), newColor, 10)
	tokenColor = newColor

	// start auction
	req := solo.NewCallParams(ScName, FuncStartAuction,
		ParamColor, tokenColor,
		ParamMinimumBid, 500,
		ParamDescription, "Cool tokens for sale!",
	).WithTransfers(map[balance.Color]int64{
		balance.ColorIOTA: 25, // deposit, must be >=minimum*margin
		tokenColor:        10, // the tokens to auction
	})
	_, err = chain.PostRequestSync(req, auctioneer)
	require.NoError(t, err)
	return chain
}

func TestDeploy(t *testing.T) {
	chain := common.StartChainAndDeployWasmContractByName(t, ScName)
	_, err := chain.FindContract(ScName)
	require.NoError(t, err)
}

func TestFaStartAuction(t *testing.T) {
	chain := setupTest(t)

	// note 1 iota should be stuck in the delayed finalize_auction
	chain.AssertAccountBalance(common.ContractAccount, balance.ColorIOTA, 25-1)
	chain.AssertAccountBalance(common.ContractAccount, tokenColor, 10)

	// auctioneer sent 25 deposit + 10 tokenColor + used 1 for request
	chain.Env.AssertAddressBalance(auctioneer.Address(), balance.ColorIOTA, solo.Saldo-35-1)
	// 1 used for request was sent back to auctioneer's account on chain
	account := coretypes.NewAgentIDFromSigScheme(auctioneer)
	chain.AssertAccountBalance(account, balance.ColorIOTA, 1)
}

func TestFaAuctionInfo(t *testing.T) {
	chain := setupTest(t)

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)

	expectedAgent := coretypes.NewAgentIDFromSigScheme(auctioneer)
	actualAgentID := chain.Env.MustGetAgentID(res[VarCreator])
	require.EqualValues(t, expectedAgent, actualAgentID)

	const expectedBidders = int64(0)
	actualBidders := chain.Env.MustGetInt64(res[VarBidders])
	require.EqualValues(t, expectedBidders, actualBidders)
}

func TestFaNoBids(t *testing.T) {
	chain := setupTest(t)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)

	const expectedBidders = int64(0)
	actualBidders := chain.Env.MustGetInt64(res[VarBidders])
	require.EqualValues(t, expectedBidders, actualBidders)
}

func TestFaOneBidTooLow(t *testing.T) {
	chain := setupTest(t)

	req := solo.NewCallParams(ScName, FuncPlaceBid,
		ParamColor, tokenColor,
	).WithTransfer(balance.ColorIOTA, 100)
	_, err := chain.PostRequestSync(req, auctioneer)
	require.Error(t, err)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)

	const expectedHighestBid = int64(-1)
	actualHighestBid := chain.Env.MustGetInt64(res[VarHighestBid])
	require.EqualValues(t, expectedHighestBid, actualHighestBid)

	const expectedBidders = int64(0)
	actualBidders := chain.Env.MustGetInt64(res[VarBidders])
	require.EqualValues(t, expectedBidders, actualBidders)
}

func TestFaOneBid(t *testing.T) {
	chain := setupTest(t)

	bidder := chain.Env.NewSignatureSchemeWithFunds()
	req := solo.NewCallParams(ScName, FuncPlaceBid,
		ParamColor, tokenColor,
	).WithTransfer(balance.ColorIOTA, 500)
	_, err := chain.PostRequestSync(req, bidder)
	require.NoError(t, err)

	// wait for finalize_auction
	chain.Env.AdvanceClockBy(61 * time.Minute)
	chain.WaitForEmptyBacklog()

	res, err := chain.CallView(
		ScName, ViewGetInfo,
		ParamColor, tokenColor,
	)
	require.NoError(t, err)

	const expectedBidders = int64(1)
	actualBidders := chain.Env.MustGetInt64(res[VarBidders])
	require.EqualValues(t, expectedBidders, actualBidders)

	expectedAgent := coretypes.NewAgentIDFromSigScheme(bidder)
	actualAgentID := chain.Env.MustGetAgentID(res[VarHighestBidder])
	require.EqualValues(t, expectedAgent, actualAgentID)
}
