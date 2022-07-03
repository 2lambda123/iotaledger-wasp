package evmtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGasUsageMemoryContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTest := env.deployGasTestMemoryContract(ethKey)

	for i := uint32(1000); i <= 100000; i += 1000 {
		res, err := gasTest.f(i)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", i, res.iscpReceipt.GasBurned)
	}
}

func TestGasUsageStorageContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTest := env.deployGasTestStorageContract(ethKey)

	for i := uint32(100); i <= 5000; i += 100 {
		res, err := gasTest.f(i)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", i, res.iscpReceipt.GasBurned)
	}
}

func TestGasUsageExecutionTimeContract(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	gasTestContract := env.deployGasTestExecutionTimeContract(ethKey)

	for i := uint32(1000); i <= 60000; i += 1000 {
		res, err := gasTestContract.f(i)
		require.NoError(t, err)
		t.Logf("n = %d, gas used: %d", i, res.iscpReceipt.GasBurned)
	}
}
