package accounts

import (
	"github.com/samber/lo"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
)

func baseTokensKey(accountKey kv.Key) kv.Key {
	return prefixBaseTokens + accountKey
}

func getBaseTokens(state kv.KVStoreReader, accountKey kv.Key) iotago.BaseToken {
	return iotago.BaseToken(lo.Must(codec.Uint64.Decode(state.Get(baseTokensKey(accountKey)), 0)))
}

func setBaseTokens(state kv.KVStore, accountKey kv.Key, n iotago.BaseToken) {
	state.Set(baseTokensKey(accountKey), codec.Uint64.Encode(uint64(n)))
}

func AdjustAccountBaseTokens(state kv.KVStore, account isc.AgentID, adjustment int64, chainID isc.ChainID) {
	switch {
	case adjustment > 0:
		CreditToAccount(state, account, isc.NewFungibleTokens(iotago.BaseToken(adjustment), nil), chainID)
	case adjustment < 0:
		DebitFromAccount(state, account, isc.NewFungibleTokens(iotago.BaseToken(-adjustment), nil), chainID)
	}
}

func GetBaseTokensBalance(state kv.KVStoreReader, agentID isc.AgentID, chainID isc.ChainID) iotago.BaseToken {
	return getBaseTokens(state, accountKey(agentID, chainID))
}
