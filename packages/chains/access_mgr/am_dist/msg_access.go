// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"errors"
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// Send by a node which has a chain enabled to a node it considers an access node.
type msgAccess struct {
	gpa.BasicMessage
	senderLClock    int
	receiverLClock  int
	accessForChains []isc.ChainID
	serverForChains []isc.ChainID
}

var _ gpa.Message = &msgAccess{}

func newMsgAccess(
	recipient gpa.NodeID,
	senderLClock, receiverLClock int,
	accessForChains []isc.ChainID,
	serverForChains []isc.ChainID,
) gpa.Message {
	return &msgAccess{
		BasicMessage:    gpa.NewBasicMessage(recipient),
		senderLClock:    senderLClock,
		receiverLClock:  receiverLClock,
		accessForChains: accessForChains,
		serverForChains: serverForChains,
	}
}

func (msg *msgAccess) MarshalBinary() ([]byte, error) {
	return rwutil.WriterToBytes(msg), nil
}

func (msg *msgAccess) UnmarshalBinary(data []byte) error {
	_, err := rwutil.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgAccess) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeAccess {
		return errors.New("msgType != msgTypeAccess")
	}
	msg.senderLClock = int(rr.ReadUint32())
	msg.receiverLClock = int(rr.ReadUint32())

	size := rr.ReadSize()
	msg.accessForChains = make([]isc.ChainID, size)
	for i := range msg.accessForChains {
		rr.ReadN(msg.accessForChains[i][:])
	}

	size = rr.ReadSize()
	msg.serverForChains = make([]isc.ChainID, size)
	for i := range msg.serverForChains {
		rr.ReadN(msg.serverForChains[i][:])
	}
	return rr.Err
}

func (msg *msgAccess) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteByte(msgTypeAccess)
	ww.WriteUint32(uint32(msg.senderLClock))
	ww.WriteUint32(uint32(msg.receiverLClock))

	ww.WriteSize(len(msg.accessForChains))
	for i := range msg.accessForChains {
		ww.WriteN(msg.accessForChains[i][:])
	}

	ww.WriteSize(len(msg.serverForChains))
	for i := range msg.serverForChains {
		ww.WriteN(msg.serverForChains[i][:])
	}
	return ww.Err
}
