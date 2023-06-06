// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mostefaoui

import (
	"bytes"
	"encoding"
	"fmt"

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

func (m *msgDone) Recipient() gpa.NodeID {
	return m.recipient
}

func (m *msgDone) SetSender(sender gpa.NodeID) {
	m.sender = sender
}

func (m *msgDone) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeDone)
	ww.WriteUint16(uint16(m.round))
	return w.Bytes(), nil
}

func (m *msgDone) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if msgType := rr.ReadByte(); msgType != msgTypeDone {
		return fmt.Errorf("expected msgTypeDone, got %v", msgType)
	}
	m.round = int(rr.ReadUint16())
	return nil
}
