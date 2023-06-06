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
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeVote)
	ww.WriteUint16(uint16(m.round))
	ww.WriteByte(byte(m.voteType))
	ww.WriteBool(m.value)
	return w.Bytes(), nil
}

func (m *msgVote) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if msgType := rr.ReadByte(); msgType != msgTypeVote {
		return fmt.Errorf("expected msgTypeVote, got %v", msgType)
	}
	m.round = int(rr.ReadUint16())
	voteType := rr.ReadByte()
	m.voteType = msgVoteType(voteType)
	m.value = rr.ReadBool()
	return nil
}
