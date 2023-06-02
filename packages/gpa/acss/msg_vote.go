// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/gpa"
)

type msgVoteKind byte

const (
	msgVoteOK msgVoteKind = iota
	msgVoteREADY
)

// This message is used a vote for the "Bracha-style totality" agreement.
type msgVote struct {
	sender    gpa.NodeID
	recipient gpa.NodeID
	kind      msgVoteKind
}

var _ gpa.Message = &msgVote{}

func (m *msgVote) Recipient() gpa.NodeID {
	return m.recipient
}

func (m *msgVote) SetSender(sender gpa.NodeID) {
	m.sender = sender
}

func (m *msgVote) MarshalBinary() ([]byte, error) {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeVote)
	mu.WriteByte(byte(m.kind))
	return mu.Bytes(), nil
}

func (m *msgVote) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	t, err := mu.ReadByte()
	if err != nil {
		return err
	}
	if t != msgTypeVote {
		return fmt.Errorf("unexpected msgType: %v in acss.msgVote", t)
	}
	k, err := mu.ReadByte()
	if err != nil {
		return err
	}
	m.kind = msgVoteKind(k)
	return nil
}
