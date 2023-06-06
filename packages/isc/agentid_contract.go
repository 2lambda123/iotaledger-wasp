package isc

import (
	"fmt"
	"io"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/util"
)

// ContractAgentID is an AgentID formed by a ChainID and a contract Hname.
type ContractAgentID struct {
	chainID ChainID
	hname   Hname
}

var _ AgentIDWithL1Address = &ContractAgentID{}

func NewContractAgentID(chainID ChainID, hname Hname) *ContractAgentID {
	return &ContractAgentID{chainID: chainID, hname: hname}
}

func contractAgentIDFromMarshalUtil(mu *marshalutil.MarshalUtil) (AgentID, error) {
	chainID, err := ChainIDFromMarshalUtil(mu)
	if err != nil {
		return nil, err
	}

	h, err := HnameFromMarshalUtil(mu)
	if err != nil {
		return nil, err
	}

	return NewContractAgentID(chainID, h), nil
}

func contractAgentIDFromString(hnamePart, addrPart string) (AgentID, error) {
	chainID, err := ChainIDFromString(addrPart)
	if err != nil {
		return nil, fmt.Errorf("NewAgentIDFromString: %w", err)
	}

	h, err := HnameFromHexString(hnamePart)
	if err != nil {
		return nil, fmt.Errorf("NewAgentIDFromString: %w", err)
	}
	return NewContractAgentID(chainID, h), nil
}

func (a *ContractAgentID) Address() iotago.Address {
	return a.chainID.AsAddress()
}

func (a *ContractAgentID) ChainID() ChainID {
	return a.chainID
}

func (a *ContractAgentID) Hname() Hname {
	return a.hname
}

func (a *ContractAgentID) Kind() AgentIDKind {
	return AgentIDKindContract
}

func (a *ContractAgentID) Bytes() []byte {
	return util.WriterToBytes(a)
}

func (a *ContractAgentID) String() string {
	return a.hname.String() + "@" + a.chainID.String()
}

func (a *ContractAgentID) Equals(other AgentID) bool {
	if other == nil {
		return false
	}
	if other.Kind() != a.Kind() {
		return false
	}
	o := other.(*ContractAgentID)
	return o.chainID.Equals(a.chainID) && o.hname == a.hname
}

// note: local read(), no need to read type byte
func (a *ContractAgentID) read(rr *util.Reader) {
	rr.Read(&a.chainID)
	rr.Read(&a.hname)
}

func (a *ContractAgentID) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteUint8(uint8(a.Kind()))
	ww.Write(&a.chainID)
	ww.Write(&a.hname)
	return ww.Err
}
