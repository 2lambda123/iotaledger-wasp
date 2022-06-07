package evmtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallFibonacci(t *testing.T) {
	evn := initEVM(t)
	ethKey, _ := evn.soloChain.NewEthereumAccountWithL2Funds()
	testCore := evn.deployTestCoreContract(ethKey)
	require.EqualValues(t, 1, evn.getBlockNumber())

	res, err := testCore.fibonacci(3)
	require.NoError(t, err)
	t.Log("evm gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
}
