// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
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
func (m *msgNextLogIndex) AsResent() *msgNextLogIndex {
	return &msgNextLogIndex{
		BasicMessage: gpa.NewBasicMessage(m.Recipient()),
		nextLogIndex: m.nextLogIndex,
		nextBaseAO:   m.nextBaseAO,
		pleaseRepeat: false,
	}
}

func (m *msgNextLogIndex) String() string {
	return fmt.Sprintf(
		"{msgNextLogIndex, sender=%v, nextLogIndex=%v, nextBaseAO=%v, pleaseRepeat=%v",
		m.Sender().ShortString(), m.nextLogIndex, m.nextBaseAO, m.pleaseRepeat,
	)
}

func (m *msgNextLogIndex) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeNextLogIndex)
	ww.WriteUint32(m.nextLogIndex.AsUint32())
	ww.WriteBytes(m.nextBaseAO.Bytes())
	ww.WriteBool(m.pleaseRepeat)
	return w.Bytes(), nil
}

func (m *msgNextLogIndex) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if msgType := rr.ReadByte(); msgType != msgTypeNextLogIndex {
		return fmt.Errorf("unexpected msgType=%v in cmtLog.msgNextLogIndex", msgType)
	}

	nextLogIndex := rr.ReadUint32()
	m.nextLogIndex = LogIndex(nextLogIndex)
	nextAOBin := rr.ReadBytes()
	nextBaseAO, err := isc.NewAliasOutputWithIDFromBytes(nextAOBin)
	if err != nil {
		return fmt.Errorf("cannot decode msgNextLogIndex.nextBaseAO: %w", err)
	}
	m.nextBaseAO = nextBaseAO
	m.pleaseRepeat = rr.ReadBool()
	return nil
}
