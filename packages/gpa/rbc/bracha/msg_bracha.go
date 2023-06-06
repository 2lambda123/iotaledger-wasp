// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bracha

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type msgBrachaType byte

const (
	msgBrachaTypePropose msgBrachaType = iota
	msgBrachaTypeEcho
	msgBrachaTypeReady
)

type msgBracha struct {
	t msgBrachaType // Type
	s gpa.NodeID    // Transient: Sender
	r gpa.NodeID    // Transient: Recipient
	v []byte        // Value
}

var _ gpa.Message = &msgBracha{}

func (m *msgBracha) Recipient() gpa.NodeID {
	return m.r
}

func (m *msgBracha) SetSender(sender gpa.NodeID) {
	m.s = sender
}

func (m *msgBracha) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(byte(m.t))
	ww.WriteBytes(m.v)
	return w.Bytes(), nil
}

func (m *msgBracha) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	t := rr.ReadByte()
	m.t = msgBrachaType(t)
	m.v = rr.ReadBytes()
	return nil
}
