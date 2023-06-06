// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"bytes"
	"fmt"

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
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeImplicateRecover)
	ww.WriteByte(byte(m.kind))
	ww.WriteUint16(uint16(m.i))
	ww.WriteBytes(m.data)
	return w.Bytes(), nil
}

func (m *msgImplicateRecover) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if t := rr.ReadByte(); t != msgTypeImplicateRecover {
		return fmt.Errorf("unexpected msgType: %v in acss.msgImplicateRecover", t)
	}
	k := rr.ReadByte()
	m.kind = msgImplicateKind(k)
	m.i = int(rr.ReadUint16())
	m.data = rr.ReadBytes()
	return nil
}
