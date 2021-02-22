package solo

import (
	"testing"

	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/stretchr/testify/require"
)

func Test_MustGetInt64Result(t *testing.T) {
	env := New(t, false, false)

	const expectedDecoded = int64(1000)
	dataBytes := codec.EncodeInt64(expectedDecoded)
	actualDecoded := env.MustGetInt64(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetBytesResult(t *testing.T) {
	env := New(t, false, false)

	expectedDecoded := []byte{0, 0, 1}
	actualDecoded := env.MustGetBytes(expectedDecoded)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetStringResult(t *testing.T) {
	env := New(t, false, false)

	const expectedDecoded = "test"
	dataBytes := codec.EncodeString(expectedDecoded)
	actualDecoded := env.MustGetString(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetAgentIDResult(t *testing.T) {
	env := New(t, false, false)

	keyPair := env.NewSignatureScheme()
	agentID := coretypes.NewAgentIDFromSigScheme(keyPair)

	expectedDecoded := agentID
	dataBytes := codec.EncodeAgentID(expectedDecoded)
	actualDecoded := env.MustGetAgentID(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetAddressResult(t *testing.T) {
	env := New(t, false, false)

	keyPair := env.NewSignatureScheme()
	address := keyPair.Address()

	expectedDecoded := address
	dataBytes := codec.EncodeAddress(expectedDecoded)
	actualDecoded := env.MustGetAddress(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetChainIDResult(t *testing.T) {
	env := New(t, false, false)

	chain := env.NewChain(nil, "dummyChain")
	chainID := chain.ChainID

	expectedDecoded := chainID
	dataBytes := codec.EncodeChainID(expectedDecoded)
	actualDecoded := env.MustGetChainID(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetColorResult(t *testing.T) {
	env := New(t, false, false)

	chain := env.NewChain(nil, "dummyChain")
	chainColor := chain.ChainColor

	expectedDecoded := chainColor
	dataBytes := codec.EncodeColor(expectedDecoded)
	actualDecoded := env.MustGetColor(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetContractIDResult(t *testing.T) {
	env := New(t, false, false)

	chain := env.NewChain(nil, "dummyChain")
	expectedDecoded := coretypes.NewContractID(chain.ChainID, coretypes.Hn(accounts.Interface.Name))

	dataBytes := codec.EncodeContractID(expectedDecoded)
	actualDecoded := env.MustGetContractID(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHashResult(t *testing.T) {
	env := New(t, false, false)

	chain := env.NewChain(nil, "dummyChain")
	expectedDecoded := chain.State.Hash()

	dataBytes := codec.EncodeHashValue(expectedDecoded)
	actualDecoded := env.MustGetHash(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHnameResult(t *testing.T) {
	env := New(t, false, false)

	expectedDecoded := coretypes.Hn(accounts.Interface.Name)

	dataBytes := codec.EncodeHname(expectedDecoded)
	actualDecoded := env.MustGetHname(dataBytes)

	require.Equal(t, expectedDecoded, actualDecoded)
}
