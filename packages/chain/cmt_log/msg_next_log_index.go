// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
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
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeNextLogIndex)
	mu.WriteUint32(m.nextLogIndex.AsUint32())
	util.MarshalBytes(mu, m.nextBaseAO.Bytes())
	mu.WriteBool(m.pleaseRepeat)
	return mu.Bytes(), nil
}

func (m *msgNextLogIndex) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	msgType, err := mu.ReadByte()
	if err != nil {
		return err
	}
	if msgType != msgTypeNextLogIndex {
		return fmt.Errorf("unexpected msgType=%v in cmtLog.msgNextLogIndex", msgType)
	}
	nextLogIndex, err2 := mu.ReadUint32()
	if err2 != nil {
		return fmt.Errorf("cannot unmarshal msgNextLogIndex.nextLogIndex: %w", err2)
	}
	m.nextLogIndex = LogIndex(nextLogIndex)
	nextAOBin, err := util.UnmarshalBytes(mu)
	if err != nil {
		return fmt.Errorf("cannot unmarshal msgNextLogIndex.nextBaseAO: %w", err)
	}
	m.nextBaseAO, err = isc.NewAliasOutputWithIDFromBytes(nextAOBin)
	if err != nil {
		return fmt.Errorf("cannot decode msgNextLogIndex.nextBaseAO: %w", err)
	}
	if m.pleaseRepeat, err = mu.ReadBool(); err != nil {
		return fmt.Errorf("cannot unmarshal msgNextLogIndex.pleaseRepeat: %w", err)
	}
	return nil
}
