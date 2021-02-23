package sbtests

import (
	"testing"

	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
)

func TestGetSet(t *testing.T) { run2(t, testGetSet) }
func testGetSet(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncSetInt,
		sbtestsc.ParamIntParamName, "ppp",
		sbtestsc.ParamIntParamValue, 314,
	)
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncGetInt,
		sbtestsc.ParamIntParamName, "ppp")
	require.NoError(t, err)

	retInt := chain.Env.MustGetInt64(ret["ppp"])
	require.EqualValues(t, 314, retInt)
}

func TestCallRecursive(t *testing.T) { run2(t, testCallRecursive) }
func testCallRecursive(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	cID, _ := setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncCallOnChain,
		sbtestsc.ParamIntParamValue, 31,
		sbtestsc.ParamHnameContract, cID.Hname(),
		sbtestsc.ParamHnameEP, coretypes.Hn(sbtestsc.FuncRunRecursion),
	)
	ret, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	ret, err = chain.CallView(sbtestsc.Interface.Name, sbtestsc.FuncGetCounter)
	require.NoError(t, err)

	r := chain.Env.MustGetInt64(ret[sbtestsc.VarCounter])
	require.EqualValues(t, 32, r)
}

const n = 10

func fibo(n int64) int64 {
	if n == 0 || n == 1 {
		return n
	}
	return fibo(n-1) + fibo(n-2)
}

func TestCallFibonacci(t *testing.T) { run2(t, testCallFibonacci) }
func testCallFibonacci(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncGetFibonacci,
		sbtestsc.ParamIntParamValue, n,
	)
	require.NoError(t, err)

	val := chain.Env.MustGetInt64(ret[sbtestsc.ParamIntParamValue])
	require.EqualValues(t, fibo(n), val)
}

func TestCallFibonacciIndirect(t *testing.T) { run2(t, testCallFibonacciIndirect) }
func testCallFibonacciIndirect(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	cID, _ := setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncCallOnChain,
		sbtestsc.ParamIntParamValue, n,
		sbtestsc.ParamHnameContract, cID.Hname(),
		sbtestsc.ParamHnameEP, coretypes.Hn(sbtestsc.FuncGetFibonacci),
	)
	ret, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)
	r := chain.Env.MustGetInt64(ret[sbtestsc.ParamIntParamValue])
	require.EqualValues(t, fibo(n), r)

	ret, err = chain.CallView(sbtestsc.Interface.Name, sbtestsc.FuncGetCounter)
	require.NoError(t, err)

	r = chain.Env.MustGetInt64(ret[sbtestsc.VarCounter])
	require.EqualValues(t, 1, r)
}
