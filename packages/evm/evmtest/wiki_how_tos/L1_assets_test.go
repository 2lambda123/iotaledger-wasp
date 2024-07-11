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
    "github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` L1Assets.sol -o ."
var (
    //go:embed L1Assets.abi
    L1AssetsContractABI string
    //go:embed L1Assets.bin
    L1AssetsContractBytecodeHex string
    L1AssetsContractBytecode = common.FromHex(strings.TrimSpace(L1AssetsContractBytecodeHex))
)

func TestWithdraw(t *testing.T) {
    env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
    privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

    instance := env.DeployContract(privateKey, L1AssetsContractABI, L1AssetsContractBytecode)

    // Define a mock L1 address for withdrawal
    recipientAddress := common.HexToAddress("0x1234567890Abcdef1234567890Abcdef12345678")

    // Mock the allowance that would be returned by ISC.sandbox.getAllowanceFrom
    mockAllowance := isc.NewEmptyAssets()
    mockAllowance.AddBaseTokens(1000)

    // Call the withdraw function and expect an event (assuming "Withdrawn" event)
    var withdrawnEvent struct {
        To common.Address
        Amount *big.Int
    }
    result, err := instance.CallFnExpectEvent(nil, "Withdrawn", &withdrawnEvent, "withdraw", recipientAddress)
    require.NoError(t, err)

    // Verify the event
    assert.Equal(t, recipientAddress, withdrawnEvent.To, "Recipient address should match")
    assert.Equal(t, mockAllowance.BaseTokens(), withdrawnEvent.Amount.Uint64(), "Withdrawn amount should match the allowance")

    // Get the transaction hash from the result
    txHash := result.GetTxHash()
    require.NotNil(t, txHash)

    // Wait for the transaction to be confirmed
    receipt, err := env.Chain.EVM().TransactionReceipt(txHash)
    require.NoError(t, err)
    assert.Equal(t, receipt.Status, uint64(1), "Transaction should be successful")

    // Additional check: Ensure deployer still has funds
    balance, err := env.Chain.EVM().Balance(deployer, nil)
    require.NoError(t, err)
    assert.True(t, balance.Sign() > 0, "Deployer should still have a positive balance")
}
