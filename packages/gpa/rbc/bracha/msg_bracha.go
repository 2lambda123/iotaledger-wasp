// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bracha

import (
	"io"

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

func (msg *msgBracha) Recipient() gpa.NodeID {
	return msg.r
}

func (msg *msgBracha) SetSender(sender gpa.NodeID) {
	msg.s = sender
}

func (msg *msgBracha) MarshalBinary() ([]byte, error) {
	return util.WriterToBytes(msg), nil
}

func (msg *msgBracha) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgBracha) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msg.t = msgBrachaType(rr.ReadByte())
	msg.v = rr.ReadBytes()
	return rr.Err
}

func (msg *msgBracha) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(byte(msg.t))
	ww.WriteBytes(msg.v)
	return ww.Err
}
