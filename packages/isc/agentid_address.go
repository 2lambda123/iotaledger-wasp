package isc

import (
	"io"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
)

// AddressAgentID is an AgentID backed by a non-alias address.
type AddressAgentID struct {
	a iotago.Address
}

var _ AgentIDWithL1Address = &AddressAgentID{}

func addressAgentIDFromMarshalUtil(mu *marshalutil.MarshalUtil) (AgentID, error) {
	var addr iotago.Address
	var err error
	if addr, err = AddressFromMarshalUtil(mu); err != nil {
		return nil, err
	}
	return NewAgentID(addr), nil
}

func addressAgentIDFromString(s string) (AgentID, error) {
	_, addr, err := iotago.ParseBech32(s)
	if err != nil {
		return nil, err
	}
	return NewAgentID(addr), nil
}

func (a *AddressAgentID) Address() iotago.Address {
	return a.a
}

func (a *AddressAgentID) Kind() AgentIDKind {
	return AgentIDKindAddress
}

func (a *AddressAgentID) Bytes() []byte {
	return util.WriterToBytes(a)
}

func (a *AddressAgentID) String() string {
	return a.a.Bech32(parameters.L1().Protocol.Bech32HRP)
}

func (a *AddressAgentID) Equals(other AgentID) bool {
	if other == nil {
		return false
	}
	if other.Kind() != a.Kind() {
		return false
	}
	return other.(*AddressAgentID).a.Equal(a.a)
}

// note: local read(), no need to read type byte
func (a *AddressAgentID) read(rr *util.Reader) {
	a.a = rr.ReadAddress()
}

func (a *AddressAgentID) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteUint8(uint8(a.Kind()))
	ww.WriteAddress(a.a)
	return ww.Err
}
