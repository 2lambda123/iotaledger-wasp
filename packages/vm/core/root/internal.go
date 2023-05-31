package root

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
)

func GetContractRegistry(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, StateVarContractRegistry)
}

func GetContractRegistryR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, StateVarContractRegistry)
}

// FindContract is an internal utility function which finds a contract in the KVStore
// It is called from within the 'root' contract as well as VMContext and viewcontext objects
// It is not directly exposed to the sandbox
// If contract is not found by the given hname, nil is returned
func FindContract(state kv.KVStoreReader, hname isc.Hname) *ContractRecord {
	contractRegistry := GetContractRegistryR(state)
	retBin := contractRegistry.GetAt(hname.Bytes())
	if retBin != nil {
		ret, err := ContractRecordFromBytes(retBin)
		if err != nil {
			panic(fmt.Errorf("FindContract: %w", err))
		}
		return ret
	}
	if hname == Contract.Hname() {
		// it happens during bootstrap
		return ContractRecordFromContractInfo(Contract)
	}
	return nil
}

func ContractExists(state kv.KVStoreReader, hname isc.Hname) bool {
	return GetContractRegistryR(state).HasAt(hname.Bytes())
}

// DecodeContractRegistry encodes the whole contract registry from the map into a Go map.
func DecodeContractRegistry(contractRegistry *collections.ImmutableMap) (map[isc.Hname]*ContractRecord, error) {
	ret := make(map[isc.Hname]*ContractRecord)
	var err error
	contractRegistry.Iterate(func(k []byte, v []byte) bool {
		var deploymentHash isc.Hname
		deploymentHash, err = isc.HnameFromBytes(k)
		if err != nil {
			return false
		}

		cr, err2 := ContractRecordFromBytes(v)
		if err2 != nil {
			return false
		}

		ret[deploymentHash] = cr
		return true
	})
	return ret, err
}

type BlockContextSubscription struct {
	Contract  isc.Hname
	OpenFunc  isc.Hname
	CloseFunc isc.Hname
}

func (s *BlockContextSubscription) Encode() []byte {
	b := make([]byte, 0, 12)
	b = append(b, codec.EncodeHname(s.Contract)...)
	b = append(b, codec.EncodeHname(s.OpenFunc)...)
	b = append(b, codec.EncodeHname(s.CloseFunc)...)
	return b
}

func mustDecodeBlockContextSubscription(b []byte) (s BlockContextSubscription) {
	if len(b) != 12 {
		panic("invalid length")
	}
	s.Contract = codec.MustDecodeHname(b[0:4])
	s.OpenFunc = codec.MustDecodeHname(b[4:8])
	s.CloseFunc = codec.MustDecodeHname(b[8:12])
	return
}

func getBlockContextSubscriptions(state kv.KVStore) *collections.Array {
	return collections.NewArray(state, StateVarBlockContextSubscriptions)
}

func getBlockContextSubscriptionsR(state kv.KVStoreReader) *collections.ArrayReadOnly {
	return collections.NewArrayReadOnly(state, StateVarBlockContextSubscriptions)
}

func SubscribeBlockContext(state kv.KVStore, contract, openFunc, closeFunc isc.Hname) {
	s := BlockContextSubscription{
		Contract:  contract,
		OpenFunc:  openFunc,
		CloseFunc: closeFunc,
	}
	getBlockContextSubscriptions(state).Push(s.Encode())
}

// GetBlockContextSubscriptions returns all contracts that are subscribed to block context,
// in deterministic order
func GetBlockContextSubscriptions(state kv.KVStoreReader) []BlockContextSubscription {
	subscriptions := getBlockContextSubscriptionsR(state)
	ret := make([]BlockContextSubscription, 0, subscriptions.Len())
	for i := range ret {
		ret = append(ret, mustDecodeBlockContextSubscription(subscriptions.GetAt(uint32(i))))
	}
	return ret
}
