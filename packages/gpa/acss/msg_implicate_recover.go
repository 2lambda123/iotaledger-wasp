// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"errors"
	"io"

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

func (msg *msgImplicateRecover) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgImplicateRecover) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgImplicateRecover) MarshalBinary() ([]byte, error) {
	return util.WriterToBytes(msg), nil
}

func (msg *msgImplicateRecover) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgImplicateRecover) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeImplicateRecover {
		return errors.New("unexpected message type")
	}
	msg.kind = msgImplicateKind(rr.ReadByte())
	msg.i = int(rr.ReadUint16())
	msg.data = rr.ReadBytes()
	return rr.Err
}

func (msg *msgImplicateRecover) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeImplicateRecover)
	ww.WriteByte(byte(msg.kind))
	ww.WriteUint16(uint16(msg.i))
	ww.WriteBytes(msg.data)
	return ww.Err
}
