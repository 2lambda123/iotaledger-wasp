package solidity_test

import (
	_ "embed"
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotaledger/wasp/packages/solo"
	// import utils
	//"wasp/documentation/tutorial-examples/solidity/utils"
)

//go:generate solc --abi --bin --overwrite ../contracts/GetBalance.sol -o build/
var (
	//go:embed build/GetBalance.abi
	GetBalanceContractABI string
	//go:embed build/GetBalance.bin
	GetBalanceContractBytecodeHex string
	GetBalanceContractBytecode    = common.FromHex(strings.TrimSpace(GetBalanceContractBytecodeHex))
)

func TestGetBalance(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true, Debug: true})
	chain := env.NewChain()

	private_key, user_address := chain.NewEthereumAccountWithL2Funds()

	contract_addr, abi := chain.DeployEVMContract(private_key, GetBalanceContractABI, GetBalanceContractBytecode, &big.Int{})

	_, token_id, _ := chain.NewNativeTokenParams(1000000000).CreateFoundry()

	callArgs, _ := abi.Pack("getBalance", []byte(token_id.ToHex()))

	signer, _ := chain.EVM().Signer()
	gas := uint64(1000000)

	callMsg := ethereum.CallMsg{
		From:     user_address,
		To:       &contract_addr,
		Data:     callArgs,
		Value:    big.NewInt(0),
		GasPrice: chain.EVM().GasPrice(),
	}

	transaction, _ := types.SignNewTx(private_key, signer, &types.LegacyTx{
		Nonce:    1,
		Gas:      gas,
		GasPrice: callMsg.GasPrice,
		To:       callMsg.To,
		Value:    big.NewInt(0),
		Data:     callMsg.Data,
	})

	chain.EVM().SendTransaction(transaction)

	receipt := chain.EVM().TransactionReceipt(transaction.Hash())

	topic := abi.Events["GotBaseBalance"].ID

	var base_balance uint64
	for _, log := range receipt.Logs {
		if slices.Contains(log.Topics, topic) {
			err := abi.UnpackIntoInterface(base_balance, "GotBaseBalance", log.Data)
			t.Log(err)
		}
	}

	t.Log("Value: ", base_balance)

}
