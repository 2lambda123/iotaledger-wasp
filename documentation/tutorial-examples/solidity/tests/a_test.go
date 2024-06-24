package solidity_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
	"github.com/stretchr/testify/assert"
)

//go:generate solc --abi --bin --overwrite ../contracts/GetBalance.sol -o ./build/
var (
	//go:embed build/GetBalance.abi
	GetBalanceContractABI string
	//go:embed build/GetBalance.bin
	GetBalanceContractBytecodeHex string
	GetBalanceContractBytecode    = common.FromHex(strings.TrimSpace(GetBalanceContractBytecodeHex))
)

func TestGetBalancea(t *testing.T) {

	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	private_key, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(private_key, GetBalanceContractABI, GetBalanceContractBytecode)

	_, token_id, _ := env.Chain.NewNativeTokenParams(1000000000).CreateFoundry()

	t.Log("token_id: ", token_id.ToHex())

	balance, _ := env.Chain.EVM().Balance(deployer, nil)

	decimals := env.Chain.EVM().BaseToken().Decimals

	var value uint64
	instance.CallFnExpectEvent(nil, "GotBaseBalance", &value, "getBalance", []byte(token_id.ToHex()))

	real := util.BaseTokensDecimalsToEthereumDecimals(value, decimals)

	assert.Equal(t, balance, real)

	sender_agent_id := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// call the contract
	var agent_id []byte
	instance.CallFnExpectEvent(nil, "GotAgentID", &agent_id, "getAgentID", []byte(token_id.ToHex()))

	t.Log("agent_id: ", sender_agent_id)

}
