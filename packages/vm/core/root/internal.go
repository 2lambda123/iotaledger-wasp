package root

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
)

func (s *StateWriter) SetInitialState(v isc.SchemaVersion, contracts []*coreutil.ContractInfo) {
	s.SetSchemaVersion(v)

	contractRegistry := s.GetContractRegistry()
	if contractRegistry.Len() != 0 {
		panic("contract registry must be empty on chain start")
	}

	// forbid deployment of custom contracts by default
	s.SetDeployPermissionsEnabled(true)

	for _, c := range contracts {
		s.StoreContractRecord(ContractRecordFromContractInfo(c))
	}
}

var errContractAlreadyExists = coreerrors.Register("contract with hname %08x already exists")

func (s *StateWriter) StoreContractRecord(rec *ContractRecord) {
	hname := isc.Hn(rec.Name)
	// storing contract record in the registry
	contractRegistry := s.GetContractRegistry()
	if contractRegistry.HasAt(hname.Bytes()) {
		panic(errContractAlreadyExists.Create(hname))
	}
	contractRegistry.SetAt(hname.Bytes(), rec.Bytes())
}

func (s *StateWriter) SetDeployPermissionsEnabled(enabled bool) {
	s.state.Set(varDeployPermissionsEnabled, codec.Bool.Encode(enabled))
}

func (s *StateReader) GetDeployPermissionsEnabled() bool {
	return lo.Must(codec.Bool.Decode(s.state.Get(varDeployPermissionsEnabled)))
}

func (s *StateWriter) GetDeployPermissions() *collections.Map {
	return collections.NewMap(s.state, varDeployPermissions)
}

func (s *StateReader) GetDeployPermissions() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, varDeployPermissions)
}

func (s *StateWriter) GetContractRegistry() *collections.Map {
	return collections.NewMap(s.state, varContractRegistry)
}

func (s *StateReader) GetContractRegistry() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, varContractRegistry)
}

// FindContract is an internal utility function which finds a contract in the KVStore
// It is called from within the 'root' contract as well as VMContext and viewcontext objects
// It is not directly exposed to the sandbox
// If contract is not found by the given hname, nil is returned
func (s *StateReader) FindContract(hname isc.Hname) *ContractRecord {
	contractRegistry := s.GetContractRegistry()
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

// decodeContractRegistry encodes the whole contract registry from the map into a Go map.
func decodeContractRegistry(contractRegistry *collections.ImmutableMap) (map[isc.Hname]*ContractRecord, error) {
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
