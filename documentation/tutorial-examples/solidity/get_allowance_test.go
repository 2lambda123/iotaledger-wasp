package solidity_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/solo"
)

// compile the solidity contract

//go:generate solc --abi --bin --overwrite GetAllowance.sol -o .
var (
	//go:embed GetAllowance.abi
	GetAllowanceContarctABI string
	//go:embed GetAllowance.bin
	GetAllowanceContractBytecodeHex string
	GetAllowanceeContractBytecode   = common.FromHex(strings.TrimSpace(GetAllowanceContractBytecodeHex))
)

func TestGetAllowance(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()

	chainID, chainOwnerID, _ := chain.GetInfo()

	t.Log("chain", chainID.String())
	t.Log("chain owner ID: ", chainOwnerID.String())

	private_key, user_address := chain.NewEthereumAccountWithL2Funds()
	t.Log("Private key of the userWallet is: ", private_key)
	t.Log("Address of the userWallet is: ", user_address)

	contract_addr, abi := chain.DeployEVMContract(private_key, GetAllowanceContarctABI, GetAllowanceeContractBytecode, &big.Int{})

	t.Log("contract_addr: ", contract_addr, abi)

	callArgs, _ := abi.Pack("getAllowanceFrom", user_address)
	callMsg := ethereum.CallMsg{
		To:   &contract_addr,
		Data: callArgs,
	}

	result, _ := chain.EVM().CallContract(callMsg, nil)

	t.Log("result: ", result)

}
