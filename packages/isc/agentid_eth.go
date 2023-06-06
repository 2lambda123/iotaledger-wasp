package isc

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/util"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
)

// EthereumAddressAgentID is an AgentID formed by an Ethereum address
type EthereumAddressAgentID struct {
	eth common.Address
}

var _ AgentID = &EthereumAddressAgentID{}

func NewEthereumAddressAgentID(eth common.Address) *EthereumAddressAgentID {
	return &EthereumAddressAgentID{eth: eth}
}

func ethAgentIDFromMarshalUtil(mu *marshalutil.MarshalUtil) (AgentID, error) {
	var ethBytes []byte
	var err error
	if ethBytes, err = mu.ReadBytes(common.AddressLength); err != nil {
		return nil, err
	}
	var eth common.Address
	eth.SetBytes(ethBytes)
	return NewEthereumAddressAgentID(eth), nil
}

func ethAgentIDFromString(s string) (AgentID, error) {
	eth := common.HexToAddress(s)
	return NewEthereumAddressAgentID(eth), nil
}

func (a *EthereumAddressAgentID) EthAddress() common.Address {
	return a.eth
}

func (a *EthereumAddressAgentID) Kind() AgentIDKind {
	return AgentIDKindEthereumAddress
}

func (a *EthereumAddressAgentID) Bytes() []byte {
	return util.WriterToBytes(a)
}

func (a *EthereumAddressAgentID) String() string {
	return a.eth.String() // includes "0x"
}

func (a *EthereumAddressAgentID) Equals(other AgentID) bool {
	if other == nil {
		return false
	}
	if other.Kind() != a.Kind() {
		return false
	}
	return other.(*EthereumAddressAgentID).eth == a.eth
}

// note: local read(), no need to read type byte
func (a *EthereumAddressAgentID) read(rr *util.Reader) {
	rr.ReadN(a.eth[:])
}

func (a *EthereumAddressAgentID) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteUint8(uint8(a.Kind()))
	ww.WriteN(a.eth[:])
	return ww.Err
}
