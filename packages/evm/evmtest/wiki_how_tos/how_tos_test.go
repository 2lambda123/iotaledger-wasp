package wiki_how_tos_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` GetBalance.sol TakeAllowance.sol Allowance.sol -o ."
var (
	//go:embed GetBalance.abi
	GetBalanceContractABI string
	//go:embed GetBalance.bin
	GetBalanceContractBytecodeHex string
	GetBalanceContractBytecode    = common.FromHex(strings.TrimSpace(GetBalanceContractBytecodeHex))

	//go:embed MintToken.abi
	MintTokenContractABI string
	//go:embed MintToken.bin
	MintTokenContractBytecodeHex string
	MintTokenContractBytecode    = common.FromHex(strings.TrimSpace(MintTokenContractBytecodeHex))
)

func TestBaseBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	balance, _ := env.Chain.EVM().Balance(deployer, nil)
	decimals := env.Chain.EVM().BaseToken().Decimals
	var value uint64
	instance.CallFnExpectEvent(nil, "GotBaseBalance", &value, "getBalanceBaseTokens")
	realBalance := util.BaseTokensDecimalsToEthereumDecimals(value, decimals)
	assert.Equal(t, balance, realBalance)
}

func TestNativeBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// create a new native token on L1
	foundry, tokenID, err := env.Chain.NewNativeTokenParams(100000000000000).CreateFoundry()
	require.NoError(t, err)
	// the token id in bytes, used to call the contract
	nativeTokenIDBytes := isc.NativeTokenIDToBytes(tokenID)

	// mint some native tokens to the chain originator
	err = env.Chain.MintTokens(foundry, 10000000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// send some native tokens to the contract deployer
	// and check if the balance returned by the contract is correct
	err = env.Chain.SendFromL2ToL2AccountNativeTokens(tokenID, senderAgentID, 100000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	nativeBalance := new(big.Int)
	instance.CallFnExpectEvent(nil, "GotNativeTokenBalance", &nativeBalance, "getBalanceNativeTokens", nativeTokenIDBytes)
	assert.Equal(t, int64(100000), nativeBalance.Int64())
}

func TestNFTBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// mint an NFToken to the contract deployer
	// and check if the balance returned by the contract is correct
	mockMetaData := []byte("sesa")
	nfti, info, err := env.Chain.Env.MintNFTL1(env.Chain.OriginatorPrivateKey, env.Chain.OriginatorAddress, mockMetaData)
	require.NoError(t, err)
	env.Chain.MustDepositNFT(nfti, env.Chain.OriginatorAgentID, env.Chain.OriginatorPrivateKey)

	transfer := isc.NewEmptyAssets()
	transfer.AddNFTs(info.NFTID)

	// send the NFT to the contract deployer
	err = env.Chain.SendFromL2ToL2Account(transfer, senderAgentID, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the NFT balance of the contract deployer
	nftBalance := new(big.Int)
	instance.CallFnExpectEvent(nil, "GotNFTIDs", &nftBalance, "getBalanceNFTs")
	assert.Equal(t, int64(1), nftBalance.Int64())
}

func TestAgentID(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// get the agnetId of the contract deployer
	// and compare it with the agentId returned by the contract
	var agentID []byte
	instance.CallFnExpectEvent(nil, "GotAgentID", &agentID, "getAgentID")
	assert.Equal(t, senderAgentID.Bytes(), agentID)
}

func TestMintNativeToken(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, _ := env.Chain.NewEthereumAccountWithL2Funds()

	// Deploy the contract
	var tokenDecimals uint8 = 6
	maximumSupply := new(big.Int)
	maximumSupply.SetString("100000", 10)
	var storageDeposit uint64 = 3
	multiplier := new(big.Int)
	multiplier.Exp(big.NewInt(10), big.NewInt(12), nil) // 10^12

	amount := new(big.Int).Mul(big.NewInt(int64(storageDeposit)), multiplier)
	amountStr := amount.String()

	value := new(big.Int)
	value.SetString(amountStr, 10)

	instance2, err := env.Chain.DeployEVMContract(
		privateKey,
		MintTokenContractABI,
		MintTokenContractBytecode,
		value,
		"Test Token",
		"TT",
		tokenDecimals,
		maximumSupply,
		storageDeposit,
	)

	t.Log(instance2, err)
}
