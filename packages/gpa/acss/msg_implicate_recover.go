// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util/rwutil"
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

var _ gpa.Message = new(msgImplicateRecover)

func (msg *msgImplicateRecover) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgImplicateRecover) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgImplicateRecover) MarshalBinary() ([]byte, error) {
	return rwutil.MarshalBinary(msg)
}

func (msg *msgImplicateRecover) UnmarshalBinary(data []byte) error {
	return rwutil.UnmarshalBinary(data, msg)
}

func (msg *msgImplicateRecover) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadMessageTypeAndVerify(msgTypeImplicateRecover)
	msg.kind = msgImplicateKind(rr.ReadByte())
	msg.i = int(rr.ReadUint16())
	msg.data = rr.ReadBytes()
	return rr.Err
}

func (msg *msgImplicateRecover) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteMessageType(msgTypeImplicateRecover)
	ww.WriteByte(byte(msg.kind))
	ww.WriteUint16(uint16(msg.i))
	ww.WriteBytes(msg.data)
	return ww.Err
}
