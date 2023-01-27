package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/client/chainclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/utxodb"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

const inccounterName = "inc"

func deployInccounter42(e *chainEnv) *isc.ContractAgentID {
	e.deployWasmInccounter(42)
	counterValue := e.getCounterValue()
	require.EqualValues(e.t, 42, counterValue)

	// test calling root.FuncFindContractByName view function using client
	ret, err := e.Chain.Cluster.WaspClient(0).CallView(
		e.Chain.ChainID, root.Contract.Hname(), root.ViewFindContract.Name,
		dict.Dict{
			root.ParamHname: incHname.Bytes(),
		})
	require.NoError(e.t, err)
	recb, err := ret.Get(root.ParamContractRecData)
	require.NoError(e.t, err)
	rec, err := root.ContractRecordFromBytes(recb)
	require.NoError(e.t, err)
	require.EqualValues(e.t, incDescription, rec.Description)

	e.expectCounter(42)
	return isc.NewContractAgentID(e.Chain.ChainID, incHname)
}

// executed in cluster_test.go
func testPostDeployInccounter(t *testing.T, e *chainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounterName, contractID.String())
}

// executed in cluster_test.go
func testPost1Request(t *testing.T, e *chainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounterName, contractID.String())

	myWallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	myClient := e.Chain.SCClient(contractID.Hname(), myWallet)

	tx, err := myClient.PostRequest(incrementFuncName)
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx, 30*time.Second)
	require.NoError(t, err)

	e.expectCounter(43)
}

// executed in cluster_test.go
func testPost3Recursive(t *testing.T, e *chainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounterName, contractID.String())

	myWallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	// fund the contract, so it can post requests to itself

	tx, err := e.NewChainClient().Post1Request(accounts.Contract.Hname(), accounts.FuncTransferAllowanceTo.Hname(), chainclient.PostRequestParams{
		Transfer: isc.NewFungibleBaseTokens(1_500_000),
		Args: map[kv.Key][]byte{
			accounts.ParamAgentID: codec.EncodeAgentID(
				isc.NewContractAgentID(e.Chain.ChainID, incHname),
			),
		},
		Allowance: isc.NewAllowanceBaseTokens(1_000_000),
	})
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx, 30*time.Second)
	require.NoError(t, err)

	myClient := e.Chain.SCClient(contractID.Hname(), myWallet)

	tx, err = myClient.PostRequest(incrementRepeatManyFuncName, chainclient.PostRequestParams{
		// Transfer:  isc.NewFungibleBaseTokens(10 * isc.Million),
		// Allowance: isc.NewAllowanceBaseTokens(9 * isc.Million),
		Args: codec.MakeDict(map[string]interface{}{
			varNumRepeats: 3,
		}),
	})
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx, 30*time.Second)
	require.NoError(t, err)

	e.waitUntilCounterEquals(contractID.Hname(), 43+3, 10*time.Second)
}

// executed in cluster_test.go
func testPost5Requests(t *testing.T, e *chainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounterName, contractID.String())

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	myAgentID := isc.NewAgentID(myAddress)
	myClient := e.Chain.SCClient(contractID.Hname(), myWallet)

	e.checkBalanceOnChain(myAgentID, isc.BaseTokenID, 0)
	onChainBalance := uint64(0)
	for i := 0; i < 5; i++ {
		baseTokesSent := 1 * isc.Million
		tx, err := myClient.PostRequest(incrementFuncName, chainclient.PostRequestParams{
			Transfer: isc.NewFungibleTokens(baseTokesSent, nil),
		})
		require.NoError(t, err)
		receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx, 30*time.Second)
		require.NoError(t, err)
		onChainBalance += baseTokesSent - receipts[0].GasFeeCharged
	}

	e.expectCounter(42 + 5)
	e.checkBalanceOnChain(myAgentID, isc.BaseTokenID, onChainBalance)

	e.checkLedger()
}

// executed in cluster_test.go
func testPost5AsyncRequests(t *testing.T, e *chainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounterName, contractID.String())

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	myAgentID := isc.NewAgentID(myAddress)

	myClient := e.Chain.SCClient(contractID.Hname(), myWallet)

	tx := [5]*iotago.Transaction{}
	onChainBalance := uint64(0)
	baseTokesSent := 1 * isc.Million
	for i := 0; i < 5; i++ {
		tx[i], err = myClient.PostRequest(incrementFuncName, chainclient.PostRequestParams{
			Transfer: isc.NewFungibleTokens(baseTokesSent, nil),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx[i], 30*time.Second)
		require.NoError(t, err)
		onChainBalance += baseTokesSent - receipts[0].GasFeeCharged
	}

	e.expectCounter(42 + 5)
	e.checkBalanceOnChain(myAgentID, isc.BaseTokenID, onChainBalance)

	if !e.Clu.AssertAddressBalances(myAddress,
		isc.NewFungibleBaseTokens(utxodb.FundsFromFaucetAmount-5*baseTokesSent)) {
		t.Fatal()
	}
	e.checkLedger()
}
