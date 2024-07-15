package wiki_how_tos_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	//"github.com/iotaledger/wasp/packages/isc"
	//"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"

	//"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
	"github.com/stretchr/testify/require"
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` L1Assets.sol -o ."
var (
	//go:embed L1Assets.abi
	L1AssetsContractABI string
	//go:embed L1Assets.bin
	L1AssetsContractBytecodeHex string
	L1AssetsContractBytecode    = common.FromHex(strings.TrimSpace(L1AssetsContractBytecodeHex))
)

func TestWithdraw(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	_, receiver := env.Chain.Env.NewKeyPair()

	// Deploy L1Assets contract
	instance := env.DeployContract(privateKey, L1AssetsContractABI, L1AssetsContractBytecode)

	require.Zero(t, env.Chain.Env.L1BaseTokens(receiver))
	senderInitialBalance := env.Chain.L2BaseTokens(isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer))

	// transfer 1 mil from ethAddress L2 to receiver L1
	transfer := 1 * isc.Million

    // create a new native token on L1
	foundry, tokenID, err := env.Chain.NewNativeTokenParams(100000000000000).CreateFoundry()
	require.NoError(t, err)
	// the token id in bytes, used to call the contract
	nativeTokenIDBytes := isc.NativeTokenIDToBytes(tokenID)

    // mint some native tokens to the chain originator
	err = env.Chain.MintTokens(foundry, 10000000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

//     // Create ISCAssets with native tokens
//     assets := isc.ISCAssets{
//         NativeTokens: []isc.NativeToken{{ID: nativeTokenIDBytes, Amount: 1000}},
// }

	// Allow the L1Assets contract to withdraw the funds
	_, err = instance.CallFn(nil, "allow", deployer, nativeTokenIDBytes)
	require.NoError(t, err)

	// Withdraw funds to receiver using the withdraw function of L1Assets contract
	_, err = instance.CallFn(nil, "withdraw", iscmagic.WrapL1Address(receiver), transfer)
	require.NoError(t, err)
	require.GreaterOrEqual(t, env.Chain.Env.L1BaseTokens(receiver), transfer-500)

	// Verify balances
	require.LessOrEqual(t, env.Chain.L2BaseTokens(isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)), senderInitialBalance-transfer)
}
