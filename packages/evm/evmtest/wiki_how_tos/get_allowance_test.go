package wiki_how_tos_test

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
//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` GetAllowance.sol -o ."

var (
	//go:embed GetAllowance.abi
	GetAllowanceContractABI string
	//go:embed GetAllowance.bin
	GetAllowanceContractBytecodeHex string
	GetAllowanceContractBytecode    = common.FromHex(strings.TrimSpace(GetAllowanceContractBytecodeHex))
)

func TestGetAllowance(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()

	chainID, chainOwnerID, _ := chain.GetInfo()

	t.Log("chain", chainID.String())
	t.Log("chain owner ID:", chainOwnerID.String())

	privateKey, userAddress := chain.NewEthereumAccountWithL2Funds()

	contractAddr, abi := chain.DeployEVMContract(privateKey, GetAllowanceContractABI, GetAllowanceContractBytecode, &big.Int{})

	t.Log("contract address:", contractAddr)
	t.Log("contract ABI:", abi)

	callArgs, _ := abi.Pack("getAllowanceFrom", userAddress)
	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: callArgs,
	}

	result, _ := chain.EVM().CallContract(callMsg, nil)

	t.Log("result:", result)
}
