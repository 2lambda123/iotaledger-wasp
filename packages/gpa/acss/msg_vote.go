// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"errors"
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
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

func (msg *msgVote) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgVote) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgVote) MarshalBinary() ([]byte, error) {
	return util.WriterToBytes(msg), nil
}

func (msg *msgVote) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgVote) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeVote {
		return errors.New("unexpected message type")
	}
	msg.kind = msgVoteKind(rr.ReadByte())
	return rr.Err
}

func (msg *msgVote) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeVote)
	ww.WriteByte(byte(msg.kind))
	return ww.Err
}
