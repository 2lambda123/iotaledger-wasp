package solidity_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate solc --abi --bin --overwrite ../contracts/GetBalance.sol -o ./build/
var (
	//go:embed build/GetBalance.abi
	GetBalanceContractABI string
	//go:embed build/GetBalance.bin
	GetBalanceContractBytecodeHex string
	GetBalanceContractBytecode    = common.FromHex(strings.TrimSpace(GetBalanceContractBytecodeHex))
)

func TestGetBalancea(t *testing.T) {

	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// create a new native token on L1
	foundry, tokenId, err := env.Chain.NewNativeTokenParams(100000000000000).CreateFoundry()
	require.NoError(t, err)
	// the token id in bytes, used to call the contract
	nativeTokenIdBytes := isc.NativeTokenIDToBytes(tokenId)

	// mint some native tokens to the chain originator
	err = env.Chain.MintTokens(foundry, 10000000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the agentId of the contract deployer
	senderAgentId := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// test 1
	// get the actual base balance of the contract deployer
	// and compare it with the balance returned by the contract
	balance, _ := env.Chain.EVM().Balance(deployer, nil)
	decimals := env.Chain.EVM().BaseToken().Decimals
	var value uint64
	instance.CallFnExpectMultipleEvents(nil, "GotBaseBalance", &value, "getBalance", nativeTokenIdBytes)
	real := util.BaseTokensDecimalsToEthereumDecimals(value, decimals)
	assert.Equal(t, balance, real)

	// test 2
	// get the agnetId of the contract deployer
	// and compare it with the agentId returned by the contract
	var agent_id []byte
	instance.CallFnExpectMultipleEvents(nil, "GotAgentID", &agent_id, "getBalance", nativeTokenIdBytes)
	assert.Equal(t, senderAgentId.Bytes(), agent_id)

	// test 3
	// get the native token balance of the contract deployer
	// It should be 0, because the contract deployer has not received any native tokens yet
	nativeBalance := new(big.Int)
	instance.CallFnExpectMultipleEvents(nil, "GotNativeTokenBalance", &nativeBalance, "getBalance", nativeTokenIdBytes)
	assert.Equal(t, int64(0), nativeBalance.Int64())

	// test 4
	// send some native tokens to the contract deployer
	// and check if the balance returned by the contract is correct
	err = env.Chain.SendFromL2ToL2AccountNativeTokens(tokenId, senderAgentId, 100000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)
	instance.CallFnExpectMultipleEvents(nil, "GotNativeTokenBalance", &nativeBalance, "getBalance", nativeTokenIdBytes)
	assert.Equal(t, int64(100000), nativeBalance.Int64())

	// test 5
	// mint an NFToken to the contract deployer
	// and check if the balance returned by the contract is correct
	mockMetaData := []byte("sesa")
	nfti, info, err := env.Chain.Env.MintNFTL1(env.Chain.OriginatorPrivateKey, env.Chain.OriginatorAddress, mockMetaData)
	require.NoError(t, err)
	env.Chain.MustDepositNFT(nfti, env.Chain.OriginatorAgentID, env.Chain.OriginatorPrivateKey)

	transfer := isc.NewEmptyAssets()
	transfer.AddNFTs(info.NFTID)

	// send the NFT to the contract deployer
	err = env.Chain.SendFromL2ToL2Account(transfer, senderAgentId, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the NFT balance of the contract deployer
	nftBalance := new(big.Int)
	instance.CallFnExpectMultipleEvents(nil, "GotNFTIDs", &nftBalance, "getBalance", nativeTokenIdBytes)
	assert.Equal(t, int64(1), nftBalance.Int64())

}
