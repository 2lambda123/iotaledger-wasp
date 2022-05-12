// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmtest

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/contracts/native/inccounter"
	"github.com/iotaledger/wasp/packages/chain/mempool"
	"github.com/iotaledger/wasp/packages/evm/evmtest"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/evm/isccontract"
	"github.com/stretchr/testify/require"
)

func TestDeploy(t *testing.T) {
	initEVM(t)
}

func TestFaucetBalance(t *testing.T) {
	evmChain := initEVM(t)
	bal := evmChain.getBalance(evmChain.faucetAddress())
	require.Zero(t, evmChain.faucetSupply.Cmp(bal))
}

func TestStorageContract(t *testing.T) {
	evmChain := initEVM(t)

	// deploy solidity `storage` contract
	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)
	require.EqualValues(t, 1, evmChain.getBlockNumber())

	// call FuncCallView to call EVM contract's `retrieve` view, get 42
	require.EqualValues(t, 42, storage.retrieve())

	// call FuncSendTransaction with EVM tx that calls `store(43)`
	res, err := storage.store(43)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, res.evmReceipt.Status)
	require.EqualValues(t, 2, evmChain.getBlockNumber())

	// call `retrieve` view, get 43
	require.EqualValues(t, 43, storage.retrieve())
}

func TestERC20Contract(t *testing.T) {
	evmChain := initEVM(t)

	// deploy solidity `erc20` contract
	erc20 := evmChain.deployERC20Contract(evmChain.faucetKey, "TestCoin", "TEST")

	// call `totalSupply` view
	{
		v := erc20.totalSupply()
		// 100 * 10^18
		expected := new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		require.Zero(t, v.Cmp(expected))
	}

	_, recipientAddress := generateEthereumKey(t)
	transferAmount := big.NewInt(1337)

	// call `transfer` => send 1337 TestCoin to recipientAddress
	res, err := erc20.transfer(recipientAddress, transferAmount)
	require.NoError(t, err)

	require.Equal(t, types.ReceiptStatusSuccessful, res.evmReceipt.Status)
	require.Equal(t, 1, len(res.evmReceipt.Logs))

	// call `balanceOf` view => check balance of recipient = 1337 TestCoin
	require.Zero(t, erc20.balanceOf(recipientAddress).Cmp(transferAmount))
}

func TestGetCode(t *testing.T) {
	evmChain := initEVM(t)
	erc20 := evmChain.deployERC20Contract(evmChain.faucetKey, "TestCoin", "TEST")

	// get contract bytecode from EVM emulator
	retrievedBytecode := evmChain.getCode(erc20.address)

	// ensure returned bytecode matches the expected runtime bytecode
	require.True(t, bytes.Equal(retrievedBytecode, evmtest.ERC20ContractRuntimeBytecode), "bytecode retrieved from the chain must match the deployed bytecode")
}

func TestGasCharged(t *testing.T) {
	evmChain := initEVM(t)
	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)

	iotaWallet, _ := evmChain.solo.NewKeyPairWithFunds()

	// call `store(999)` with enough gas
	res, err := storage.store(999, ethCallOptions{iota: iotaCallOptions{wallet: iotaWallet}})
	require.NoError(t, err)
	t.Log("evmChain gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
	require.Greater(t, res.evmReceipt.GasUsed, uint64(0))
	require.Greater(t, res.iscpReceipt.GasBurned, uint64(0))
	require.Greater(t, res.iscpReceipt.GasFeeCharged, uint64(0))
}

func TestGasRatio(t *testing.T) {
	evmChain := initEVM(t)
	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)

	require.Equal(t, evm.DefaultGasRatio, evmChain.getGasRatio())

	res, err := storage.store(43)
	require.NoError(t, err)
	initialGasFee := res.iscpReceipt.GasFeeCharged

	// only the owner can call the setGasRatio endpoint
	newGasRatio := util.Ratio32{A: evm.DefaultGasRatio.A * 10, B: evm.DefaultGasRatio.B}
	newUserWallet, _ := evmChain.solo.NewKeyPairWithFunds()
	err = evmChain.setGasRatio(newGasRatio, iotaCallOptions{wallet: newUserWallet})
	require.True(t, iscp.VMErrorIs(err, vm.ErrUnauthorized))
	require.Equal(t, evm.DefaultGasRatio, evmChain.getGasRatio())

	// current owner is able to set a new gasRatio
	err = evmChain.setGasRatio(newGasRatio, iotaCallOptions{wallet: evmChain.soloChain.OriginatorPrivateKey})
	require.NoError(t, err)
	require.Equal(t, newGasRatio, evmChain.getGasRatio())

	// run an equivalent request and compare the gas fees
	res, err = storage.store(44)
	require.NoError(t, err)
	require.Greater(t, res.iscpReceipt.GasFeeCharged, initialGasFee)
}

// tests that the gas limits are correctly enforced based on the iotas sent
func TestGasLimit(t *testing.T) {
	evmChain := initEVM(t)
	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)

	// set a gas ratio such that evmChain gas cost in iotas is larger than dust cost
	err := evmChain.setGasRatio(util.Ratio32{A: 10, B: 1}, iotaCallOptions{wallet: evmChain.soloChain.OriginatorPrivateKey})
	require.NoError(t, err)

	// estimate gas by sending a valid tx
	iotaWallet1, _ := evmChain.solo.NewKeyPairWithFunds()
	result, err := storage.store(123, ethCallOptions{iota: iotaCallOptions{wallet: iotaWallet1}})
	require.NoError(t, err)
	gas := result.iscpReceipt.GasBurned
	fee := result.iscpReceipt.GasFeeCharged
	t.Logf("gas: %d, fee: %d", gas, fee)

	// send again with same gas limit but not enough iotas
	iotaWallet2, _ := evmChain.solo.NewKeyPairWithFunds()
	_, err = storage.store(124, ethCallOptions{iota: iotaCallOptions{
		wallet: iotaWallet2,
		before: func(req *solo.CallParams) {
			req.WithGasBudget(gas).AddIotas(fee * 9 / 10)
		},
	}})
	require.Error(t, err)
	require.Regexp(t, `\bgas\b`, err.Error())

	// send again with gas limit not enough for transaction
	iotaWallet3, _ := evmChain.solo.NewKeyPairWithFunds()
	_, err = storage.store(125, ethCallOptions{iota: iotaCallOptions{
		wallet: iotaWallet3,
		before: func(req *solo.CallParams) {
			req.WithGasBudget(gas / 2).AddIotas(fee)
		},
	}})
	require.Error(t, err)
	require.Regexp(t, `\bgas\b`, err.Error())
}

// ensure the amount of iotas sent impacts the amount of gas used
func TestLoop(t *testing.T) {
	evmChain := initEVM(t)
	loop := evmChain.deployLoopContract(evmChain.faucetKey)

	gasPerToken := evmChain.soloChain.GetGasFeePolicy().GasPerToken

	for _, iscGasBudget := range []uint64{200000, 400000} {
		iotasSent := iscGasBudget / gasPerToken
		iotaWallet, iotaAddress := evmChain.solo.NewKeyPairWithFunds()
		loop.loop(ethCallOptions{
			iota: iotaCallOptions{
				wallet: iotaWallet,
				before: func(req *solo.CallParams) {
					req.WithGasBudget(iscGasBudget).AddIotas(iotasSent)
				},
			},
		})
		// gas fee is charged regardless of result
		require.LessOrEqual(t,
			evmChain.soloChain.L2Iotas(iscp.NewAgentID(iotaAddress)),
			iotasSent-iscGasBudget,
		)
	}
}

func TestPrePaidFees(t *testing.T) {
	evmChain := initEVM(t)
	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)

	iotaWallet, iotaAddress := evmChain.solo.NewKeyPairWithFunds()

	// test sending off-ledger request without depositing funds first
	txdata, _, _ := storage.buildEthTxData(nil, "store", uint32(999))
	offledgerRequest := evmChain.buildSoloRequest(evm.FuncSendTransaction.Name, evm.FieldTransactionData, txdata).
		WithMaxAffordableGasBudget()
	_, err := evmChain.soloChain.PostRequestOffLedger(offledgerRequest, iotaWallet)
	require.Error(t, err)
	require.EqualValues(t, 42, storage.retrieve())

	// deposit funds
	initialBalance := evmChain.solo.L1Iotas(iotaAddress)
	err = evmChain.soloChain.DepositIotasToL2(initialBalance/2, iotaWallet)
	require.NoError(t, err)

	// send offledger request again and check that is works
	_, err = evmChain.soloChain.PostRequestOffLedger(offledgerRequest, iotaWallet)
	require.NoError(t, err)
	require.EqualValues(t, 999, int(storage.retrieve()))
}

func TestISCContract(t *testing.T) {
	// deploy the evmChain contract, which starts an EVM chain and automatically
	// deploys the isc.sol EVM contract at address 0x1074
	evmChain := initEVM(t)

	// deploy the isc-test.sol EVM contract
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// call the ISCTest.getChainId() view function of isc-test.sol which in turn:
	//  calls the ISC.getChainId() view function of isc.sol at 0x1074, which:
	//   returns the ChainID of the underlying ISC chain
	chainID := iscTest.getChainID()

	require.True(t, evmChain.soloChain.ChainID.Equals(chainID))
}

func TestISCChainOwnerID(t *testing.T) {
	evmChain := initEVM(t)

	ret := new(isccontract.ISCAgentID)
	evmChain.ISCContract(evmChain.faucetKey).callView(nil, "getChainOwnerID", nil, &ret)

	chainOwnerID := evmChain.soloChain.OriginatorAgentID
	require.True(t, chainOwnerID.Equals(ret.MustUnwrap()))
}

func TestISCTimestamp(t *testing.T) {
	evmChain := initEVM(t)

	var ret int64
	evmChain.ISCContract(evmChain.faucetKey).callView(nil, "getTimestampUnixSeconds", nil, &ret)

	require.EqualValues(t, evmChain.soloChain.GetLatestBlockInfo().Timestamp.Unix(), ret)
}

func TestISCGetParam(t *testing.T) {
	evmChain := initEVM(t)

	key := string(evm.FieldCallMsg) // callView sends an ISC request including this parameter

	var has bool
	evmChain.ISCContract(evmChain.faucetKey).callView(nil, "hasParam", []interface{}{key}, &has)
	require.True(t, has)

	var ret []byte
	evmChain.ISCContract(evmChain.faucetKey).callView(nil, "getParam", []interface{}{key}, &ret)
	require.NotEmpty(t, ret)
}

func TestISCCallView(t *testing.T) {
	evmChain := initEVM(t)

	ret := new(isccontract.ISCDict)
	evmChain.ISCContract(evmChain.faucetKey).callView(nil, "callView", []interface{}{
		accounts.Contract.Hname(),
		accounts.ViewBalance.Hname(),
		&isccontract.ISCDict{Items: []isccontract.ISCDictItem{{
			Key:   []byte(accounts.ParamAgentID),
			Value: evmChain.soloChain.OriginatorAgentID.Bytes(),
		}}},
	}, &ret)

	require.NotEmpty(t, ret.Unwrap())
}

func TestISCLogPanic(t *testing.T) {
	evmChain := initEVM(t)

	_, err := evmChain.ISCContract(evmChain.faucetKey).callFn(
		[]ethCallOptions{{iota: iotaCallOptions{
			before: func(req *solo.CallParams) {
				req.AddIotas(10000).WithMaxAffordableGasBudget()
			},
		}}},
		"logPanic",
		"Hi from EVM!",
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Hi from EVM!")
}

func TestISCNFTData(t *testing.T) {
	evmChain := initEVM(t)

	// mint an NFT and send it to the chain
	issuerWallet, issuerAddress := evmChain.solo.NewKeyPairWithFunds()
	metadata := []byte("foobar")
	nftInfo, err := evmChain.solo.MintNFTL1(issuerWallet, issuerAddress, []byte("foobar"))
	require.NoError(t, err)
	_, err = evmChain.soloChain.PostRequestSync(
		solo.NewCallParams(accounts.Contract.Name, accounts.FuncDeposit.Name).
			AddIotas(100000).
			WithNFT(&iscp.NFT{
				ID:       nftInfo.NFTID,
				Issuer:   issuerAddress,
				Metadata: metadata,
			}).
			WithMaxAffordableGasBudget().
			WithSender(nftInfo.NFTID.ToAddress()),
		issuerWallet,
	)
	require.NoError(t, err)

	// call getNFTData from EVM
	ret := new(isccontract.ISCNFT)
	evmChain.ISCContract(evmChain.faucetKey).callView(
		nil,
		"getNFTData",
		[]interface{}{isccontract.WrapIotaNFTID(nftInfo.NFTID)},
		&ret,
	)

	require.EqualValues(t, nftInfo.NFTID, ret.MustUnwrap().ID)
	require.True(t, issuerAddress.Equal(ret.MustUnwrap().Issuer))
	require.EqualValues(t, metadata, ret.MustUnwrap().Metadata)
}

func TestISCTriggerEvent(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// call ISCTest.triggerEvent(string) function of isc-test.sol which in turn:
	//  calls the ISC.iscpTriggerEvent(string) function of isc.sol at 0x1074, which:
	//   triggers an ISC event with the given string parameter
	res, err := iscTest.triggerEvent("Hi from EVM!")
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, res.evmReceipt.Status)
	ev, err := evmChain.soloChain.GetEventsForBlock(evmChain.soloChain.GetLatestBlockInfo().BlockIndex)
	require.NoError(t, err)
	require.Len(t, ev, 1)
	require.Contains(t, ev[0], "Hi from EVM!")
}

func TestISCTriggerEventThenFail(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// test that triggerEvent() followed by revert() does not actually trigger the event
	_, err := iscTest.triggerEventFail("Hi from EVM!", ethCallOptions{iota: iotaCallOptions{
		before: func(req *solo.CallParams) {
			req.AddIotas(10000).WithMaxAffordableGasBudget()
		},
	}})
	require.Error(t, err)
	ev, err := evmChain.soloChain.GetEventsForBlock(evmChain.soloChain.GetLatestBlockInfo().BlockIndex)
	require.NoError(t, err)
	require.Len(t, ev, 0)
}

func TestISCEntropy(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// call the ISCTest.emitEntropy() function of isc-test.sol which in turn:
	//  calls ISC.iscpEntropy() function of isc.sol at 0x1074, which:
	//   returns the entropy value from the sandbox
	//  emits an EVM event (aka log) with the entropy value
	var entropy hashing.HashValue
	iscTest.callFnExpectEvent(nil, "EntropyEvent", &entropy, "emitEntropy")

	require.NotEqualValues(t, hashing.NilHash, entropy)
}

func TestISCGetRequestID(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	reqID := new(iscp.RequestID)
	iscTest.callFnExpectEvent(nil, "RequestIDEvent", &reqID, "emitRequestID")

	require.EqualValues(t, evmChain.soloChain.LastReceipt().Request.ID(), *reqID)
}

func TestISCGetCaller(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	agentID := new(isccontract.ISCAgentID)
	iscTest.callFnExpectEvent(nil, "GetCallerEvent", &agentID, "emitGetCaller")

	originatorAgentID := iscp.NewAgentID(evmChain.soloChain.OriginatorAddress)
	require.True(t, originatorAgentID.Equals(agentID.MustUnwrap()))
}

func TestISCGetSenderAccount(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	sender := new(isccontract.ISCAgentID)
	iscTest.callFnExpectEvent(nil, "SenderAccountEvent", &sender, "emitSenderAccount")

	require.EqualValues(t, isccontract.WrapISCAgentID(evmChain.soloChain.LastReceipt().Request.SenderAccount()), *sender)
}

func TestISCGetSenderAddress(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	sender := new(isccontract.IotaAddress)
	iscTest.callFnExpectEvent(nil, "SenderAddressEvent", &sender, "emitSenderAddress")

	require.EqualValues(t, isccontract.WrapIotaAddress(evmChain.soloChain.LastReceipt().Request.SenderAddress()), *sender)
}

func TestISCGetAllowanceIotas(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	var iotas uint64
	iscTest.callFnExpectEvent([]ethCallOptions{{iota: iotaCallOptions{
		before: func(req *solo.CallParams) {
			req.AddAllowanceIotas(42)
		},
	}}}, "AllowanceIotasEvent", &iotas, "emitAllowanceIotas")

	require.EqualValues(t, 42, iotas)
}

func TestISCGetAllowanceAvailableIotas(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	var iotasAvailable uint64
	iscTest.callFnExpectEvent([]ethCallOptions{{iota: iotaCallOptions{
		before: func(cp *solo.CallParams) {
			cp.AddAllowanceIotas(42)
		},
	}}}, "AllowanceAvailableIotasEvent", &iotasAvailable, "emitAllowanceAvailableIotas")

	require.EqualValues(t, 42, iotasAvailable)
}

func TestRevert(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// mint some native tokens
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gas
	sn, _, err := evmChain.soloChain.NewFoundryParams(10000).
		WithUser(evmChain.soloChain.OriginatorPrivateKey).
		CreateFoundry()
	require.NoError(t, err)
	err = evmChain.soloChain.MintTokens(sn, 10000, evmChain.soloChain.OriginatorPrivateKey)
	require.NoError(t, err)

	err = iscTest.callFnExpectError([]ethCallOptions{{
		iota: iotaCallOptions{
			before: func(req *solo.CallParams) {
				req.AddIotas(200000).
					WithMaxAffordableGasBudget()
			},
		},
	}}, "emitRevertVMError")

	t.Log(err.Error())
	require.Error(t, err)
	require.EqualValues(t, err.Error(), "PostRequestSync failed: execution reverted: contractId: 07cb02c1, errorId: 62505")
}

func TestSend(t *testing.T) {
	evmChain := initEVM(t, inccounter.Processor)
	err := evmChain.soloChain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash)
	require.NoError(t, err)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gas

	iscTest.callFnExpectEvent([]ethCallOptions{{iota: iotaCallOptions{
		before: func(req *solo.CallParams) {
			req.AddIotas(200000).
				WithMaxAffordableGasBudget()
		},
	}}}, "SendEvent", nil, "emitSend")
}

func TestSendAsNFT(t *testing.T) {
	evmChain := initEVM(t, inccounter.Processor)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	err := evmChain.soloChain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash)
	require.NoError(t, err)

	// mint an NFT and send to chain
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gas
	issuerWallet, issuerAddress := evmChain.solo.NewKeyPairWithFunds()
	metadata := []byte("foobar")
	nftInfo, err := evmChain.solo.MintNFTL1(issuerWallet, issuerAddress, metadata)
	require.NoError(t, err)

	_, err = iscTest.callFn([]ethCallOptions{{
		iota: iotaCallOptions{
			wallet: issuerWallet,
			before: func(cp *solo.CallParams) {
				cp.AddIotas(100000).
					WithNFT(&iscp.NFT{
						ID:       nftInfo.NFTID,
						Issuer:   issuerAddress,
						Metadata: metadata,
					}).
					AddAllowanceNFTs(nftInfo.NFTID).
					WithMaxAffordableGasBudget()
			},
		},
	}}, "callSendAsNFT", isccontract.WrapIotaNFTID(nftInfo.NFTID))
	require.NoError(t, err)
}

func TestISCGetAllowanceAvailableNativeTokens(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// mint some native tokens
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gaz
	sn, tokenID, err := evmChain.soloChain.NewFoundryParams(10000).
		WithUser(evmChain.soloChain.OriginatorPrivateKey).
		CreateFoundry()
	require.NoError(t, err)
	err = evmChain.soloChain.MintTokens(sn, 10000, evmChain.soloChain.OriginatorPrivateKey)
	require.NoError(t, err)

	nt := new(isccontract.IotaNativeToken)
	iscTest.callFnExpectEvent([]ethCallOptions{{
		iota: iotaCallOptions{
			before: func(cp *solo.CallParams) {
				cp.AddAllowanceNativeTokens(&tokenID, 42)
			},
		},
	}}, "AllowanceAvailableNativeTokenEvent", &nt, "emitAllowanceAvailableNativeTokens")

	require.EqualValues(t, tokenID[:], nt.ID.Data)
	require.EqualValues(t, 42, nt.Amount.Uint64())
}

func TestISCGetAllowanceNFTs(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// mint an NFT and send to chain
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gas
	issuerWallet, issuerAddress := evmChain.solo.NewKeyPairWithFunds()
	metadata := []byte("foobar")
	nftInfo, err := evmChain.solo.MintNFTL1(issuerWallet, issuerAddress, metadata)
	require.NoError(t, err)

	nft := new(isccontract.ISCNFT)
	iscTest.callFnExpectEvent([]ethCallOptions{{
		iota: iotaCallOptions{
			wallet: issuerWallet,
			before: func(cp *solo.CallParams) {
				cp.AddIotas(10000).
					WithNFT(&iscp.NFT{
						ID:       nftInfo.NFTID,
						Issuer:   issuerAddress,
						Metadata: metadata,
					}).
					WithAllowance(&iscp.Allowance{
						Assets: &iscp.FungibleTokens{Iotas: 1000},
						NFTs:   []iotago.NFTID{nftInfo.NFTID},
					}).
					WithGasBudget(1000000)
			},
		},
	}}, "AllowanceNFTEvent", &nft, "emitAllowanceNFTs")

	require.EqualValues(t, nftInfo.NFTID, nft.MustUnwrap().ID)
	require.True(t, issuerAddress.Equal(nft.MustUnwrap().Issuer))
	require.EqualValues(t, metadata, nft.MustUnwrap().Metadata)
}

func TestISCGetAllowanceAvailableNFTs(t *testing.T) {
	evmChain := initEVM(t)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	// mint an NFT and send to chain
	evmChain.soloChain.MustDepositIotasToL2(10_000, nil) // for gas
	issuerWallet, issuerAddress := evmChain.solo.NewKeyPairWithFunds()
	metadata := []byte("foobar")
	nftInfo, err := evmChain.solo.MintNFTL1(issuerWallet, issuerAddress, metadata)
	require.NoError(t, err)

	nft := new(isccontract.ISCNFT)
	iscTest.callFnExpectEvent([]ethCallOptions{{
		iota: iotaCallOptions{
			wallet: issuerWallet,
			before: func(cp *solo.CallParams) {
				cp.AddIotas(10000).
					WithNFT(&iscp.NFT{
						ID:       nftInfo.NFTID,
						Issuer:   issuerAddress,
						Metadata: metadata,
					}).
					WithAllowance(&iscp.Allowance{
						Assets: &iscp.FungibleTokens{Iotas: 1000},
						NFTs:   []iotago.NFTID{nftInfo.NFTID},
					}).
					WithGasBudget(1000000)
			},
		},
	}}, "AllowanceAvailableNFTEvent", &nft, "emitAllowanceAvailableNFTs")

	require.EqualValues(t, nftInfo.NFTID, nft.MustUnwrap().ID)
	require.True(t, issuerAddress.Equal(nft.MustUnwrap().Issuer))
	require.EqualValues(t, metadata, nft.MustUnwrap().Metadata)
}

func TestISCCall(t *testing.T) {
	evmChain := initEVM(t, inccounter.Processor)
	err := evmChain.soloChain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash)
	require.NoError(t, err)
	iscTest := evmChain.deployISCTestContract(evmChain.faucetKey)

	res, err := iscTest.callFn(nil, "callInccounter")
	require.NoError(evmChain.solo.T, err)
	require.Equal(evmChain.solo.T, types.ReceiptStatusSuccessful, res.evmReceipt.Status)

	r, err := evmChain.soloChain.CallView(
		inccounter.Contract.Name,
		inccounter.FuncGetCounter.Name,
	)
	require.NoError(evmChain.solo.T, err)
	require.EqualValues(t, 42, codec.MustDecodeInt64(r.MustGet(inccounter.VarCounter)))
}

func TestBlockTime(t *testing.T) {
	t.SkipNow() // TODO: skipping because it fails randomly

	evmChain := initEVM(t)

	// deposit funds to cover for dust, gas, etc
	_, err := evmChain.soloChain.PostRequestSync(
		solo.NewCallParams(accounts.Contract.Name, accounts.FuncTransferAllowanceTo.Name,
			accounts.ParamAgentID, iscp.NewContractAgentID(evmChain.soloChain.ChainID, evm.Contract.Hname()),
		).
			AddIotas(200000).
			AddAllowanceIotas(100000).
			WithMaxAffordableGasBudget(),
		nil,
	)
	require.NoError(t, err)

	evmChain.setBlockTime(60)

	storage := evmChain.deployStorageContract(evmChain.faucetKey, 42)
	require.EqualValues(t, 42, storage.retrieve())
	require.EqualValues(t, 0, evmChain.getBlockNumber())

	res, err := storage.store(43)
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, res.evmReceipt.Status)

	require.EqualValues(t, 43, storage.retrieve())
	require.EqualValues(t, 0, evmChain.getBlockNumber())

	// there is 1 timelocked request
	mempoolInfo := evmChain.soloChain.MempoolInfo()
	require.EqualValues(t, 1, mempoolInfo.InBufCounter-mempoolInfo.OutPoolCounter)

	// first block gets minted
	evmChain.solo.AdvanceClockBy(61*time.Second, 1)
	evmChain.soloChain.WaitUntil(func(mstats mempool.MempoolInfo) bool {
		return mstats.OutPoolCounter == mempoolInfo.InBufCounter
	})
	require.EqualValues(t, 1, evmChain.getBlockNumber())
	block := evmChain.getBlockByNumber(1)
	require.EqualValues(t, 2, len(block.Transactions()))

	// there is 1 timelocked request
	mempoolInfo = evmChain.soloChain.MempoolInfo()
	require.EqualValues(t, 1, mempoolInfo.InBufCounter-mempoolInfo.OutPoolCounter)

	// second (empty) block gets minted
	evmChain.solo.AdvanceClockBy(61*time.Second, 1)
	evmChain.soloChain.WaitUntil(func(mstats mempool.MempoolInfo) bool {
		return mstats.OutPoolCounter == mempoolInfo.InBufCounter
	})
	require.EqualValues(t, 2, evmChain.getBlockNumber())
	block = evmChain.getBlockByNumber(2)
	require.EqualValues(t, 0, len(block.Transactions()))
}
