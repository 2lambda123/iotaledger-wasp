// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"errors"
	"fmt"
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

type msgNextLogIndex struct {
	gpa.BasicMessage
	nextLogIndex LogIndex               // Proposal is to go to this LI without waiting for a consensus.
	nextBaseAO   *isc.AliasOutputWithID // Using this AO as a base.
	pleaseRepeat bool                   // If true, the receiver should resend its latest message back to the sender.
}

var _ gpa.Message = &msgNextLogIndex{}

func newMsgNextLogIndex(recipient gpa.NodeID, nextLogIndex LogIndex, nextBaseAO *isc.AliasOutputWithID, pleaseRepeat bool) *msgNextLogIndex {
	return &msgNextLogIndex{
		BasicMessage: gpa.NewBasicMessage(recipient),
		nextLogIndex: nextLogIndex,
		nextBaseAO:   nextBaseAO,
		pleaseRepeat: pleaseRepeat,
	}
}

// Make a copy for re-sending the message.
// We set pleaseResend to false to avoid accidental loops.
func (msg *msgNextLogIndex) AsResent() *msgNextLogIndex {
	return &msgNextLogIndex{
		BasicMessage: gpa.NewBasicMessage(msg.Recipient()),
		nextLogIndex: msg.nextLogIndex,
		nextBaseAO:   msg.nextBaseAO,
		pleaseRepeat: false,
	}
}

func (msg *msgNextLogIndex) String() string {
	return fmt.Sprintf(
		"{msgNextLogIndex, sender=%v, nextLogIndex=%v, nextBaseAO=%v, pleaseRepeat=%v",
		msg.Sender().ShortString(), msg.nextLogIndex, msg.nextBaseAO, msg.pleaseRepeat,
	)
}

func (msg *msgNextLogIndex) MarshalBinary() ([]byte, error) {
	return rwutil.WriterToBytes(msg), nil
}

func (msg *msgNextLogIndex) UnmarshalBinary(data []byte) error {
	_, err := rwutil.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgNextLogIndex) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeNextLogIndex {
		return errors.New("msgType != msgTypeNextLogIndex")
	}
	msg.nextLogIndex = LogIndex(rr.ReadUint32())
	msg.nextBaseAO = new(isc.AliasOutputWithID)
	rr.Read(msg.nextBaseAO)
	msg.pleaseRepeat = rr.ReadBool()
	return rr.Err
}

func (msg *msgNextLogIndex) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteByte(msgTypeNextLogIndex)
	ww.WriteUint32(msg.nextLogIndex.AsUint32())
	ww.Write(msg.nextBaseAO)
	ww.WriteBool(msg.pleaseRepeat)
	return ww.Err
}
