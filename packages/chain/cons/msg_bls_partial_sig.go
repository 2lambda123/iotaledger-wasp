// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"go.dedis.ch/kyber/v3/suites"

	"github.com/iotaledger/wasp/packages/gpa"
)

type msgBLSPartialSig struct {
	blsSuite   suites.Suite
	sender     gpa.NodeID
	recipient  gpa.NodeID
	partialSig []byte
}

var _ gpa.Message = &msgBLSPartialSig{}

func newMsgBLSPartialSig(blsSuite suites.Suite, recipient gpa.NodeID, partialSig []byte) *msgBLSPartialSig {
	return &msgBLSPartialSig{blsSuite: blsSuite, recipient: recipient, partialSig: partialSig}
}

func (msg *msgBLSPartialSig) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgBLSPartialSig) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgBLSPartialSig) MarshalBinary() ([]byte, error) {
	return msg.partialSig, nil
}

func (msg *msgBLSPartialSig) UnmarshalBinary(data []byte) error {
	msg.partialSig = data
	return nil
}
