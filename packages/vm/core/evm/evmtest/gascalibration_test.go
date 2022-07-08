package evmtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const factor = 10

func TestGasUsageMemoryContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTest := env.deployGasTestMemoryContract(ethKey)

	for i := uint32(1); i <= 100; i++ {
		n := i * factor
		res, err := gasTest.f(n)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", n, res.iscpReceipt.GasBurned)
	}
}

func TestGasUsageStorageContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTest := env.deployGasTestStorageContract(ethKey)

	for i := uint32(1); i <= 100; i++ {
		n := i * factor
		res, err := gasTest.f(n)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", n, res.iscpReceipt.GasBurned)
	}
}

func TestGasUsageExecutionTimeContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTestContract := env.deployGasTestExecutionTimeContract(ethKey)

	for i := uint32(1); i <= 100; i++ {
		n := i * factor
		res, err := gasTestContract.f(n)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", n, res.iscpReceipt.GasBurned)
	}
}
