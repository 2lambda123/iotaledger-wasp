package test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test2Chains(t *testing.T) {
	// run2(t, test2Chains)
	test2Chains(t, true)
}

func test2Chains(t *testing.T, w bool) {
	corecontracts.PrintWellKnownHnames()

	timeLayout := "04:05.000000000"
	l := testlogger.NewNamedLogger(t.Name(), timeLayout)
	l = testlogger.WithLevel(l, zapcore.ErrorLevel, true)
	env := solo.New(t, &solo.InitOptions{
		AutoAdjustStorageDeposit: true,
		// Debug:                    false,
		PrintStackTrace: true,
		Log:             l,
	})
	chain1 := wasmsolo.StartChain(t, "chain1", env)
	chain2 := wasmsolo.StartChain(t, "chain2", env)
	chain1.CheckAccountLedger()
	chain2.CheckAccountLedger()

	println("0------------------------chain1_------------------------")
	println(chain1.DumpAccounts())
	println("0------------------------chain2_------------------------")
	println(chain2.DumpAccounts())
	println("0------------------------------------------------------")

	ctx2 := deployTestCoreOnChain(t, w, chain2, nil)
	require.NoError(t, ctx2.Err)
	ctx1 := deployTestCoreOnChain(t, w, chain1, nil)
	require.NoError(t, ctx1.Err)

	userWallet, userAddress := env.NewKeyPairWithFunds()
	userAgentID := isc.NewAgentID(userAddress)
	chain1ContractAgentID := ctx1.Account().AgentID()
	chain2ContractAgentID := ctx2.Account().AgentID()
	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount)
	fmt.Println("userAgentID: ", userAgentID)
	fmt.Println("chain1.CommonAccount(): ", chain1.CommonAccount())
	fmt.Println("chain2.CommonAccount(): ", chain2.CommonAccount())

	chain1CommonAccountBaseTokens := chain1.L2BaseTokens(chain1.CommonAccount())
	chain2CommonAccountBaseTokens := chain2.L2BaseTokens(chain2.CommonAccount())

	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens)
	chain1.AssertL2BaseTokens(chain1ContractAgentID, wasmsolo.L2FundsContract)
	chain1.AssertL2TotalBaseTokens(chain1CommonAccountBaseTokens + chain1.L2BaseTokens(chain1.OriginatorAgentID) + wasmsolo.L2FundsContract)

	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens)
	chain2.AssertL2BaseTokens(chain2ContractAgentID, wasmsolo.L2FundsContract)
	chain2.AssertL2TotalBaseTokens(chain2CommonAccountBaseTokens + chain2.L2BaseTokens(chain2.OriginatorAgentID) + wasmsolo.L2FundsContract)

	chain1TotalBaseTokens := chain1.L2TotalBaseTokens()
	chain2TotalBaseTokens := chain2.L2TotalBaseTokens()

	chain1.WaitForRequestsMark()
	chain2.WaitForRequestsMark()

	println("1------------------------chain1_------------------------")
	println(chain1.DumpAccounts())
	println("1------------------------chain2_------------------------")
	println(chain2.DumpAccounts())
	println("1------------------------------------------------------")

	// send base tokens to chain2ContractAgentID (that is an entity of chain2) on chain1
	const baseTokensCreditedToScOnChain1 = 5 * isc.Million
	const baseTokensToSend = baseTokensCreditedToScOnChain1 + 1*isc.Million
	req := solo.NewCallParams(
		accounts.Contract.Name, accounts.FuncTransferAllowanceTo.Name,
		accounts.ParamAgentID, chain2ContractAgentID,
	).
		AddBaseTokens(baseTokensToSend).
		AddAllowanceBaseTokens(baseTokensCreditedToScOnChain1).
		WithGasBudget(math.MaxUint64)

	_, err := chain1.PostRequestSync(req, userWallet)
	require.NoError(t, err)

	receipt1 := chain1.LastReceipt()

	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount-baseTokensToSend)
	chain1.AssertL2BaseTokens(userAgentID, baseTokensToSend-baseTokensCreditedToScOnChain1-receipt1.GasFeeCharged)
	chain1.AssertL2BaseTokens(chain2ContractAgentID, baseTokensCreditedToScOnChain1)
	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens+receipt1.GasFeeCharged)
	chain1CommonAccountBaseTokens += receipt1.GasFeeCharged
	chain1.AssertL2TotalBaseTokens(chain1TotalBaseTokens + baseTokensToSend)
	chain1TotalBaseTokens += baseTokensToSend

	chain2.AssertL2BaseTokens(userAgentID, 0)
	chain2.AssertL2BaseTokens(chain2ContractAgentID, wasmsolo.L2FundsContract)
	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens)
	chain2.AssertL2TotalBaseTokens(chain2TotalBaseTokens)

	println("2------------------------chain1_------------------------")
	println(chain1.DumpAccounts())
	println("2------------------------chain2_------------------------")
	println(chain2.DumpAccounts())
	println("2------------------------------------------------------")

	// make chain2 send a call to chain1 to withdraw base tokens
	baseTokensToWithdrawalFromChain1 := baseTokensCreditedToScOnChain1 // try to withdraw all base tokens deposited to chain1 on behalf of chain2's contract
	// reqAllowance is the allowance provided to the "withdraw from chain" contract (chain2) that needs to be enough to
	// pay the gas fees of withdraw func on chain1
	reqAllowance := accounts.ConstDepositFeeTmp + 1*isc.Million
	// allowance + x, where x will be used for the gas costs of `FuncWithdrawFromChain` on chain2
	baseTokensToSend2 := reqAllowance + 1*isc.Million

	fmt.Println("baseTokensToWithdrawalFromChain1: ", baseTokensToWithdrawalFromChain1)
	fmt.Println("reqAllowance: ", reqAllowance)
	fmt.Println("baseTokensToSend2: ", baseTokensToSend2)
	req = solo.NewCallParams(testcore.ScName, testcore.FuncWithdrawFromChain,
		testcore.ParamChainID, chain1.ChainID,
		testcore.ParamBaseTokensWithdrawal, baseTokensToWithdrawalFromChain1).
		AddBaseTokens(baseTokensToSend2).
		WithAllowance(isc.NewAssetsBaseTokens(reqAllowance)).
		WithGasBudget(math.MaxUint64)

	fmt.Println("0 chain1.L2Assets(chain1ContractAgentID).BaseTokens: ", chain1.L2Assets(chain1ContractAgentID).BaseTokens)
	fmt.Println("0 chain1.L2Assets(chain2ContractAgentID).BaseTokens: ", chain1.L2Assets(chain2ContractAgentID).BaseTokens)
	fmt.Println("0 chain2.L2Assets(chain1ContractAgentID).BaseTokens: ", chain2.L2Assets(chain1ContractAgentID).BaseTokens)
	fmt.Println("0 chain2.L2Assets(chain2ContractAgentID).BaseTokens: ", chain2.L2Assets(chain2ContractAgentID).BaseTokens)
	fmt.Println("==============")
	_, err = chain2.PostRequestSync(req, userWallet)

	require.NoError(t, err)
	chain2SendWithdrawalReceipt := chain2.LastReceipt()

	require.True(t, chain1.WaitForRequestsThrough(2, 10*time.Second))
	require.True(t, chain2.WaitForRequestsThrough(2, 10*time.Second))
	fmt.Println("1 chain1.L2Assets(chain1ContractAgentID).BaseTokens: ", chain1.L2Assets(chain1ContractAgentID).BaseTokens)
	fmt.Println("1 chain1.L2Assets(chain2ContractAgentID).BaseTokens: ", chain1.L2Assets(chain2ContractAgentID).BaseTokens)
	fmt.Println("1 chain2.L2Assets(chain1ContractAgentID).BaseTokens: ", chain2.L2Assets(chain1ContractAgentID).BaseTokens)
	fmt.Println("1 chain2.L2Assets(chain2ContractAgentID).BaseTokens: ", chain2.L2Assets(chain2ContractAgentID).BaseTokens)
	println("3------------------------chain1_------------------------")
	println(chain1.DumpAccounts())
	println("3------------------------chain2_------------------------")
	println(chain2.DumpAccounts())
	println("3------------------------------------------------------")

	chain2DepositReceipt := chain2.LastReceipt()

	chain1WithdrawalReceipt := chain1.LastReceipt()

	require.Equal(t, chain1WithdrawalReceipt.DeserializedRequest().CallTarget().Contract, accounts.Contract.Hname())
	require.Equal(t, chain1WithdrawalReceipt.DeserializedRequest().CallTarget().EntryPoint, accounts.FuncWithdraw.Hname())
	require.Nil(t, chain1WithdrawalReceipt.Error)

	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount-baseTokensToSend-baseTokensToSend2)

	chain1.AssertL2BaseTokens(userAgentID, baseTokensToSend-baseTokensCreditedToScOnChain1-receipt1.GasFeeCharged)
	// amount of base tokens sent from chain2 to chain1 in order to call the "withdrawal" request
	chain1.AssertL2BaseTokens(chain2ContractAgentID, reqAllowance-chain1WithdrawalReceipt.GasFeeCharged)
	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens+chain1WithdrawalReceipt.GasFeeCharged)
	chain1.AssertL2TotalBaseTokens(chain1TotalBaseTokens + reqAllowance - baseTokensToWithdrawalFromChain1)

	chain2.AssertL2BaseTokens(userAgentID, baseTokensToSend2-reqAllowance-chain2SendWithdrawalReceipt.GasFeeCharged)
	chain2.AssertL2BaseTokens(chain2ContractAgentID, baseTokensToWithdrawalFromChain1-accounts.ConstDepositFeeTmp)
	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens+chain2SendWithdrawalReceipt.GasFeeCharged+chain2DepositReceipt.GasFeeCharged)
	println(chain2.DumpAccounts())
	chain2.AssertL2TotalBaseTokens(chain2TotalBaseTokens + baseTokensToSend2 - reqAllowance + baseTokensCreditedToScOnChain1)
}
