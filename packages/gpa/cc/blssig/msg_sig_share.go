// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package blssig

import (
	"encoding"

	"github.com/iotaledger/wasp/packages/gpa"
)

type msgSigShare struct {
	recipient gpa.NodeID
	sender    gpa.NodeID
	sigShare  []byte
}

var (
	_ gpa.Message              = &msgSigShare{}
	_ encoding.BinaryMarshaler = &msgSigShare{}
)

func (msg *msgSigShare) Recipient() gpa.NodeID {
	return msg.recipient
}

func (msg *msgSigShare) SetSender(sender gpa.NodeID) {
	msg.sender = sender
}

func (msg *msgSigShare) MarshalBinary() ([]byte, error) {
	return msg.sigShare, nil
}

func (msg *msgSigShare) UnmarshalBinary(data []byte) error {
	msg.sigShare = data
	return nil
}
