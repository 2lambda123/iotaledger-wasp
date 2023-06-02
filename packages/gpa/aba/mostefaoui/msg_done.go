// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mostefaoui

import (
	"encoding"
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/gpa"
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
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeDone)
	mu.WriteUint16(uint16(m.round))
	return mu.Bytes(), nil
}

func (m *msgDone) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	msgType, err := mu.ReadByte()
	if err != nil {
		return err
	}
	if msgType != msgTypeDone {
		return fmt.Errorf("expected msgTypeDone, got %v", msgType)
	}
	round, err := mu.ReadUint16()
	if err != nil {
		return err
	}
	m.round = int(round)
	return nil
}
