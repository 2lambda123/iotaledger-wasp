package solo

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/stretchr/testify/require"
)

const (
	CouldNotConvertDataInto = "Could not convert data into "
	DataDoesNotExist        = "Data does not exist."
)

// MustGetInt64 converts input data into int64. Panics when either no data is provided or cannot be converted.
func (env *Solo) MustGetInt64(data []byte) int64 {
	result, exists, err := codec.DecodeInt64(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" int64")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetBytes returns the input as is. Panics if no input is provided.
func (env *Solo) MustGetBytes(data []byte) []byte {
	require.NotNil(env.T, data, CouldNotConvertDataInto+" bytes")
	return data
}

// MustGetString converts input data into int64. Panics when either no data is provided or cannot be converted.
func (env *Solo) MustGetString(data []byte) string {
	result, exists, err := codec.DecodeString(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" string")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetAgentID converts input data into an AgentID. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetAgentID(data []byte) coretypes.AgentID {
	result, exists, err := codec.DecodeAgentID(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" AgentID")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetAddress converts input data into an Address. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetAddress(data []byte) address.Address {
	result, exists, err := codec.DecodeAddress(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" Address")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetChainID converts input data into a ChainID. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetChainID(data []byte) coretypes.ChainID {
	result, exists, err := codec.DecodeChainID(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" ChainID")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetColor converts input data into a Color. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetColor(data []byte) balance.Color {
	result, exists, err := codec.DecodeColor(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" Color")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetContractID converts input data into a ContractID. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetContractID(data []byte) coretypes.ContractID {
	result, exists, err := codec.DecodeContractID(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" ContractID")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetHash converts input data into a HashValue. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetHash(data []byte) hashing.HashValue {
	result, exists, err := codec.DecodeHashValue(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" HashValue")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}

// MustGetHname converts input data into an Hname. Panics if no input is provided or cannot be converted.
func (env *Solo) MustGetHname(data []byte) coretypes.Hname {
	result, exists, err := codec.DecodeHname(data)
	require.NoError(env.T, err, CouldNotConvertDataInto+" Hname")
	require.True(env.T, exists, DataDoesNotExist)
	return result
}
