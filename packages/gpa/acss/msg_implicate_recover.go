// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type msgImplicateKind byte

const (
	msgImplicateRecoverKindIMPLICATE msgImplicateKind = iota
	msgImplicateRecoverKindRECOVER
)

// The <IMPLICATE, i, skᵢ> and <RECOVER, i, skᵢ> messages.
type msgImplicateRecover struct {
	sender    gpa.NodeID
	recipient gpa.NodeID
	kind      msgImplicateKind
	i         int
	data      []byte // Either implication or the recovered secret.
}

var _ gpa.Message = &msgImplicateRecover{}

func (m *msgImplicateRecover) Recipient() gpa.NodeID {
	return m.recipient
}

func (m *msgImplicateRecover) SetSender(sender gpa.NodeID) {
	m.sender = sender
}

func (m *msgImplicateRecover) MarshalBinary() ([]byte, error) {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeImplicateRecover)
	mu.WriteByte(byte(m.kind))
	mu.WriteUint16(uint16(m.i))
	util.WriteBytesMu(mu, m.data)
	return mu.Bytes(), nil
}

func (m *msgImplicateRecover) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	msgType, err := mu.ReadByte()
	if err != nil {
		return err
	}
	if msgType != msgTypeImplicateRecover {
		return fmt.Errorf("unexpected msgType: %v in acss.msgImplicateRecover", msgType)
	}
	kind, err := mu.ReadByte()
	if err != nil {
		return err
	}
	m.kind = msgImplicateKind(kind)
	i, err := mu.ReadUint16()
	if err != nil { // TODO: Resolve I from the context, trusting it might be unsafe.
		return err
	}
	m.i = int(i)
	m.data, err = util.ReadBytesMu(mu)
	return err
}
