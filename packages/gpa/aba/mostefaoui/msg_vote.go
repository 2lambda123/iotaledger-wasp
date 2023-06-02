// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mostefaoui

import (
	"encoding"
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/gpa"
)

type msgVoteType byte

const (
	BVAL msgVoteType = iota
	AUX
)

type msgVote struct {
	sender    gpa.NodeID
	recipient gpa.NodeID
	round     int
	voteType  msgVoteType
	value     bool
}

var (
	_ gpa.Message                = &msgVote{}
	_ encoding.BinaryUnmarshaler = &msgVote{}
)

func multicastMsgVote(recipients []gpa.NodeID, round int, voteType msgVoteType, value bool) gpa.OutMessages {
	msgs := gpa.NoMessages()
	for _, nid := range recipients {
		msgs.Add(&msgVote{recipient: nid, round: round, voteType: voteType, value: value})
	}
	return msgs
}

func (m *msgVote) Recipient() gpa.NodeID {
	return m.recipient
}

func (m *msgVote) SetSender(sender gpa.NodeID) {
	m.sender = sender
}

func (m *msgVote) MarshalBinary() ([]byte, error) {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeVote)
	mu.WriteUint16(uint16(m.round))
	mu.WriteByte(byte(m.voteType))
	mu.WriteBool(m.value)
	return mu.Bytes(), nil
}

func (m *msgVote) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	msgType, err := mu.ReadByte()
	if err != nil {
		return err
	}
	if msgType != msgTypeVote {
		return fmt.Errorf("expected msgTypeVote, got %v", msgType)
	}
	round, err := mu.ReadUint16()
	if err != nil {
		return err
	}
	m.round = int(round)
	voteType, err := mu.ReadByte()
	if err != nil {
		return err
	}
	m.voteType = msgVoteType(voteType)
	m.value, err = mu.ReadBool()
	return err
}
