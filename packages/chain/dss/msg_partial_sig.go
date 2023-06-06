// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dss

import (
	"bytes"
	"fmt"

	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/dss"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

type msgPartialSig struct {
	suite      suites.Suite // Transient, for un-marshaling only.
	sender     gpa.NodeID   // Transient.
	recipient  gpa.NodeID   // Transient.
	partialSig *dss.PartialSig
}

var _ gpa.Message = &msgPartialSig{}

func (m *msgPartialSig) Recipient() gpa.NodeID {
	return m.recipient
}

func (m *msgPartialSig) SetSender(sender gpa.NodeID) {
	m.sender = sender
}

func (m *msgPartialSig) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypePartialSig)
	ww.WriteUint16(uint16(m.partialSig.Partial.I)) // TODO: Resolve it from the context, instead of marshaling.
	ww.WriteMarshaled(m.partialSig.Partial.V)
	ww.WriteBytes(m.partialSig.SessionID)
	ww.WriteBytes(m.partialSig.Signature)
	return w.Bytes(), nil
}

func (m *msgPartialSig) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if msgType := rr.ReadByte(); msgType != msgTypePartialSig {
		return fmt.Errorf("unexpected msgType=%v in dss.msgPartialSig", msgType)
	}
	partialI := rr.ReadUint16()
	partialV := m.suite.Scalar()
	if err2 := util.ReadMarshaled(r, partialV); err2 != nil {
		return fmt.Errorf("cannot unmarshal partialSig.V: %w", err2)
	}
	m.partialSig = &dss.PartialSig{
		Partial: &share.PriShare{I: int(partialI), V: partialV},
	}
	m.partialSig.SessionID = rr.ReadBytes()
	m.partialSig.Signature = rr.ReadBytes()
	return nil
}
