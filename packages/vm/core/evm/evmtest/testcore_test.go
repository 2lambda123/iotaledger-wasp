package evmtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const N = 7

func fibo(n uint32) uint32 {
	if n == 0 || n == 1 {
		return n
	}
	return fibo(n-1) + fibo(n-2)
}

func TestCallFibonacci(t *testing.T) {
	evn := initEVM(t)
	ethKey, _ := evn.soloChain.NewEthereumAccountWithL2Funds()
	testCore := evn.deployTestCoreContract(ethKey)
	require.EqualValues(t, 1, evn.getBlockNumber())

	res, err := testCore.fibonacci(N)
	require.NoError(t, err)
	t.Log("evm gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
}

func TestCallFibonacciIndirect(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	testCore := env.deployTestCoreContract(ethKey)
	require.EqualValues(t, 1, env.getBlockNumber())

	res, err := testCore.fibonacciIndirect(N)
	require.NoError(t, err)
	t.Log("evm gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
}

func TestCallFibonacciLoop(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	testCore := env.deployTestCoreContract(ethKey)
	require.EqualValues(t, 1, env.getBlockNumber())

	var v uint32
	res, err := testCore.fibonacciLoop(N, &v)
	require.NoError(t, err)
	require.EqualValues(t, fibo(N), v)
	t.Log("evm gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
}

func TestCallLoop(t *testing.T) {
	env := initEVM(t)
	ethKey, _ := env.soloChain.NewEthereumAccountWithL2Funds()
	testCore := env.deployTestCoreContract(ethKey)
	require.EqualValues(t, 1, env.getBlockNumber())

	res, err := testCore.loop(1024)
	require.NoError(t, err)
	t.Log("evm gas used:", res.evmReceipt.GasUsed)
	t.Log("iscp gas used:", res.iscpReceipt.GasBurned)
	t.Log("iscp gas fee:", res.iscpReceipt.GasFeeCharged)
}
