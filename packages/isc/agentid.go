// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
)

type AgentIDKind uint8

const (
	AgentIDKindNil AgentIDKind = iota
	AgentIDKindAddress
	AgentIDKindContract
	AgentIDKindEthereumAddress
)

// AgentID represents any entity that can hold assets on L2 and/or call contracts.
type AgentID interface {
	Bytes() []byte
	Equals(other AgentID) bool
	Kind() AgentIDKind
	read(rr *util.Reader)
	String() string
	Write(w io.Writer) error
}

// AgentIDWithL1Address is an AgentID backed by an L1 address (either AddressAgentID or ContractAgentID).
type AgentIDWithL1Address interface {
	AgentID
	Address() iotago.Address
}

// AddressFromAgentID returns the L1 address of the AgentID, if applicable.
func AddressFromAgentID(a AgentID) (iotago.Address, bool) {
	wa, ok := a.(AgentIDWithL1Address)
	if !ok {
		return nil, false
	}
	return wa.Address(), true
}

// HnameFromAgentID returns the hname of the AgentID, if applicable.
func HnameFromAgentID(a AgentID) (Hname, bool) {
	ca, ok := a.(*ContractAgentID)
	if !ok {
		return 0, false
	}
	return ca.Hname(), true
}

// NewAgentID creates an AddressAgentID if the address is not an AliasAddress;
// otherwise a ContractAgentID with hname = HnameNil.
func NewAgentID(addr iotago.Address) AgentID {
	if addr.Type() == iotago.AddressAlias {
		chainID := ChainIDFromAddress(addr.(*iotago.AliasAddress))
		return NewContractAgentID(chainID, 0)
	}
	return &AddressAgentID{a: addr}
}

func AgentIDFromMarshalUtil(mu *marshalutil.MarshalUtil) (AgentID, error) {
	var err error
	kind, err := mu.ReadByte()
	if err != nil {
		return nil, err
	}
	switch AgentIDKind(kind) {
	case AgentIDKindNil:
		return &NilAgentID{}, nil
	case AgentIDKindAddress:
		return addressAgentIDFromMarshalUtil(mu)
	case AgentIDKindContract:
		return contractAgentIDFromMarshalUtil(mu)
	case AgentIDKindEthereumAddress:
		return ethAgentIDFromMarshalUtil(mu)
	}
	return nil, fmt.Errorf("no handler for AgentID kind %d", kind)
}

func AgentIDFromBytes(data []byte) (ret AgentID, err error) {
	rr := util.NewBytesReader(data)
	return AgentIDFromReader(rr), rr.Err
}

func AgentIDFromReader(rr *util.Reader) (ret AgentID) {
	switch AgentIDKind(rr.ReadInt8()) {
	case AgentIDKindNil:
		ret = new(NilAgentID)
	case AgentIDKindAddress:
		ret = new(AddressAgentID)
	case AgentIDKindContract:
		ret = new(ContractAgentID)
	case AgentIDKindEthereumAddress:
		ret = new(EthereumAddressAgentID)
	default:
		if rr.Err == nil {
			rr.Err = errors.New("invalid agent ID kind")
			return nil
		}
	}
	// note: this local read() will not attempt to read the kind byte
	ret.read(rr)
	return ret
}

// NewAgentIDFromString parses the human-readable string representation
func NewAgentIDFromString(s string) (AgentID, error) {
	if s == nilAgentIDString {
		return &NilAgentID{}, nil
	}
	var hnamePart, addrPart string
	{
		parts := strings.Split(s, "@")
		switch len(parts) {
		case 1:
			addrPart = parts[0]
		case 2:
			addrPart = parts[1]
			hnamePart = parts[0]
		default:
			return nil, errors.New("NewAgentIDFromString: wrong format")
		}
	}

	if hnamePart != "" {
		return contractAgentIDFromString(hnamePart, addrPart)
	}
	if strings.HasPrefix(addrPart, string(parameters.L1().Protocol.Bech32HRP)) {
		return addressAgentIDFromString(s)
	}
	if strings.HasPrefix(addrPart, "0x") {
		return ethAgentIDFromString(s)
	}
	return nil, errors.New("NewAgentIDFromString: wrong format")
}

// NewRandomAgentID creates random AgentID
func NewRandomAgentID() AgentID {
	return NewContractAgentID(RandomChainID(), Hn("testName"))
}
