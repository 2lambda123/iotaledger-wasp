package test

import (
	"math"
	"testing"
	"time"

	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/stretchr/testify/require"
)

//func deposit(t *testing.T, ctx *wasmsolo.SoloContext, user, target *wasmsolo.SoloAgent, amount int64) {
//	ctxAcc := ctx.SoloContextForCore(t, coreaccounts.ScName, coreaccounts.OnDispatch)
//	f := coreaccounts.ScFuncs.Deposit(ctxAcc.Sign(user))
//	if target != nil {
//		f.Params.AgentID().SetValue(target.ScAgentID())
//	}
//	f.Func.TransferIotas(amount).Post()
//	require.NoError(t, ctxAcc.Err)
//}

func setupTestSandboxSC(t *testing.T, chain *solo.Chain, user *cryptolib.KeyPair, runWasm bool) isc.AgentID {
	err := deployContract(chain, user, runWasm)
	require.NoError(t, err)

	deployed := isc.NewContractAgentID(chain.ChainID, testcore.HScName)
	req := solo.NewCallParams(testcore.ScName, testcore.FuncDoNothing.Name).
		WithGasBudget(100_000)
	_, err = chain.PostRequestSync(req, user)
	require.NoError(t, err)
	t.Logf("deployed test_sandbox'%s': %s", testcore.ScName, testcore.HScName)
	return deployed
}

func Test2Chains(t *testing.T, w bool) {
	corecontracts.PrintWellKnownHnames()

	env := solo.New(t, &solo.InitOptions{
		AutoAdjustStorageDeposit: true,
		Debug:                    true,
		PrintStackTrace:          true,
	}).
		WithNativeContract(testcore.Processor)
	chain1 := env.NewChain()
	chain2, _ := env.NewChainExt(nil, 0, "chain2")
	chain1.CheckAccountLedger()
	chain2.CheckAccountLedger()

	setupTestSandboxSC(t, chain1, nil, w)
	contractAgentID := setupTestSandboxSC(t, chain2, nil, w)

	userWallet, userAddress := env.NewKeyPairWithFunds()
	userAgentID := isc.NewAgentID(userAddress)
	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount)

	chain1CommonAccountBaseTokens := chain1.L2BaseTokens(chain1.CommonAccount())
	chain2CommonAccountBaseTokens := chain2.L2BaseTokens(chain2.CommonAccount())

	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens)
	chain1.AssertL2TotalBaseTokens(chain1CommonAccountBaseTokens + chain1.L2BaseTokens(chain1.OriginatorAgentID))

	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens)
	chain2.AssertL2TotalBaseTokens(chain2CommonAccountBaseTokens + chain2.L2BaseTokens(chain2.OriginatorAgentID))

	chain1TotalBaseTokens := chain1.L2TotalBaseTokens()
	chain2TotalBaseTokens := chain2.L2TotalBaseTokens()

	chain1.WaitForRequestsMark()
	chain2.WaitForRequestsMark()

	// send base tokens to contractAgentID (that is an entity of chain2) on chain1
	const baseTokensToSend = 11 * isc.Million
	const baseTokensCreditedToScOnChain1 = 10 * isc.Million
	req := solo.NewCallParams(
		accounts.Contract.Name, accounts.FuncTransferAllowanceTo.Name,
		accounts.ParamAgentID, contractAgentID,
	).
		AddBaseTokens(baseTokensToSend).
		AddAllowanceBaseTokens(baseTokensCreditedToScOnChain1).
		WithGasBudget(math.MaxUint64)

	_, err := chain1.PostRequestSync(req, userWallet)
	require.NoError(t, err)

	receipt1 := chain1.LastReceipt()

	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount-baseTokensToSend)
	chain1.AssertL2BaseTokens(userAgentID, baseTokensToSend-baseTokensCreditedToScOnChain1-receipt1.GasFeeCharged)
	chain1.AssertL2BaseTokens(contractAgentID, baseTokensCreditedToScOnChain1)
	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens+receipt1.GasFeeCharged)
	chain1CommonAccountBaseTokens += receipt1.GasFeeCharged
	chain1.AssertL2TotalBaseTokens(chain1TotalBaseTokens + baseTokensToSend)
	chain1TotalBaseTokens += baseTokensToSend

	chain2.AssertL2BaseTokens(userAgentID, 0)
	chain2.AssertL2BaseTokens(contractAgentID, 0)
	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens)
	chain2.AssertL2TotalBaseTokens(chain2TotalBaseTokens)

	println("----chain1------------------------------------------")
	println(chain1.DumpAccounts())
	println("-----chain2-----------------------------------------")
	println(chain2.DumpAccounts())
	println("----------------------------------------------")

	// make chain2 send a call to chain1 to withdraw base tokens
	baseTokensToWithdrawalFromChain1 := baseTokensCreditedToScOnChain1 // try to withdraw all base tokens deposited to chain1 on behalf of chain2's contract
	// reqAllowance is the allowance provided to the "withdraw from chain" contract (chain2) that needs to be enough to
	// pay the gas fees of withdraw func on chain1
	reqAllowance := accounts.ConstDepositFeeTmp + 1*isc.Million
	// allowance + x, where x will be used for the gas costs of `FuncWithdrawFromChain` on chain2
	baseTokensToSend2 := reqAllowance + 1*isc.Million

	req = solo.NewCallParams(testcore.ScName, testcore.FuncWithdrawFromChain.Name,
		testcore.ParamChainID, chain1.ChainID,
		testcore.ParamBaseTokensToWithdrawal, baseTokensToWithdrawalFromChain1).
		AddBaseTokens(baseTokensToSend2).
		WithAllowance(isc.NewAssetsBaseTokens(reqAllowance)).
		WithGasBudget(math.MaxUint64)

	_, err = chain2.PostRequestSync(req, userWallet)
	require.NoError(t, err)
	chain2SendWithdrawalReceipt := chain2.LastReceipt()

	require.True(t, chain1.WaitForRequestsThrough(2, 10*time.Second))
	require.True(t, chain2.WaitForRequestsThrough(2, 10*time.Second))

	println("----chain1------------------------------------------")
	println(chain1.DumpAccounts())
	println("-----chain2-----------------------------------------")
	println(chain2.DumpAccounts())
	println("----------------------------------------------")

	chain2DepositReceipt := chain2.LastReceipt()

	chain1WithdrawalReceipt := chain1.LastReceipt()

	require.Equal(t, chain1WithdrawalReceipt.DeserializedRequest().CallTarget().Contract, accounts.Contract.Hname())
	require.Equal(t, chain1WithdrawalReceipt.DeserializedRequest().CallTarget().EntryPoint, accounts.FuncWithdraw.Hname())
	require.Nil(t, chain1WithdrawalReceipt.Error)

	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount-baseTokensToSend-baseTokensToSend2)

	chain1.AssertL2BaseTokens(userAgentID, baseTokensToSend-baseTokensCreditedToScOnChain1-receipt1.GasFeeCharged)
	chain1.AssertL2BaseTokens(contractAgentID, reqAllowance-chain1WithdrawalReceipt.GasFeeCharged) // amount of base tokens sent from chain2 to chain1 in order to call the "withdrawal" request
	chain1.AssertL2BaseTokens(chain1.CommonAccount(), chain1CommonAccountBaseTokens+chain1WithdrawalReceipt.GasFeeCharged)
	chain1.AssertL2TotalBaseTokens(chain1TotalBaseTokens + reqAllowance - baseTokensToWithdrawalFromChain1)

	chain2.AssertL2BaseTokens(userAgentID, baseTokensToSend2-reqAllowance-chain2SendWithdrawalReceipt.GasFeeCharged)
	chain2.AssertL2BaseTokens(contractAgentID, baseTokensToWithdrawalFromChain1-accounts.ConstDepositFeeTmp)
	chain2.AssertL2BaseTokens(chain2.CommonAccount(), chain2CommonAccountBaseTokens+chain2SendWithdrawalReceipt.GasFeeCharged+chain2DepositReceipt.GasFeeCharged)
	println(chain2.DumpAccounts())
	chain2.AssertL2TotalBaseTokens(chain2TotalBaseTokens + baseTokensToSend2 - reqAllowance + baseTokensCreditedToScOnChain1)
}
