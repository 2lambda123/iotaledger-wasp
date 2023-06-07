// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mostefaoui

import (
	"encoding"
	"errors"
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util/rwutil"
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

func (msg *msgVote) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgVote) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgVote) MarshalBinary() ([]byte, error) {
	return rwutil.WriterToBytes(msg), nil
}

func (msg *msgVote) UnmarshalBinary(data []byte) error {
	_, err := rwutil.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgVote) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeVote {
		return errors.New("msgType != msgTypeVote")
	}
	msg.round = int(rr.ReadUint16())
	msg.voteType = msgVoteType(rr.ReadByte())
	msg.value = rr.ReadBool()
	return rr.Err
}

func (msg *msgVote) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteByte(msgTypeVote)
	ww.WriteUint16(uint16(msg.round))
	ww.WriteByte(byte(msg.voteType))
	ww.WriteBool(msg.value)
	return ww.Err
}
