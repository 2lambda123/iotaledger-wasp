package solidity_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/solo"
	// import the solidity package
)

// compile the solidity contract

//go:generate solc --abi --bin --overwrite NativeTokenBalance.sol -o .
var (
	//go:embed NativeTokenBalance.abi
	NativeTokenBalanceContractABI string
	//go:embed NativeTokenBalance.bin
	nativeTokenBalanceContractBytecodeHex string
	NativeTokenBalanceContractBytecode    = common.FromHex(strings.TrimSpace(nativeTokenBalanceContractBytecodeHex))
)

func TestNativeTokenBalance(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()

	chainID, chainOwnerID, _ := chain.GetInfo()

	t.Log("chainID: ", chainID.String())
	t.Log("chain owner ID: ", chainOwnerID.String())

	private_key, user_address := chain.NewEthereumAccountWithL2Funds()

	t.Log("Address of the userWallet is: ", private_key)
	t.Log("Address of the userWallet1 is: ", user_address)

	contract_addr, abi := chain.DeployEVMContract(private_key, NativeTokenBalanceContractABI, NativeTokenBalanceContractBytecode, &big.Int{})

	callArgs, _ := abi.Pack("simpleFunction")
	callMsg := ethereum.CallMsg{
		To:   &contract_addr,
		Data: callArgs,
	}

	result, _ := chain.EVM().CallContract(callMsg, nil)

	t.Log("result: ", result)

}
