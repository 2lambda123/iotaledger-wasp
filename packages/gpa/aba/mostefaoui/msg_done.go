// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mostefaoui

import (
	"encoding"
	"errors"
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type msgDone struct {
	sender    gpa.NodeID
	recipient gpa.NodeID
	round     int
}

var (
	_ gpa.Message                = &msgDone{}
	_ encoding.BinaryUnmarshaler = &msgDone{}
)

func multicastMsgDone(nodeIDs []gpa.NodeID, me gpa.NodeID, round int) gpa.OutMessages {
	msgs := gpa.NoMessages()
	for _, n := range nodeIDs {
		if n != me {
			msgs.Add(&msgDone{recipient: n, round: round})
		}
	}
	return msgs
}

func (msg *msgDone) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgDone) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgDone) MarshalBinary() ([]byte, error) {
	return util.WriterToBytes(msg), nil
}

func (msg *msgDone) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgDone) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeDone {
		return errors.New("unexpected message type")
	}
	msg.round = int(rr.ReadUint16())
	return rr.Err
}

func (msg *msgDone) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeDone)
	ww.WriteUint16(uint16(msg.round))
	return ww.Err
}
