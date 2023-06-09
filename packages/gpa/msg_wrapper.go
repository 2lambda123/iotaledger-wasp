// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gpa

import (
	"io"

	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// MsgWrapper can be used to compose an algorithm out of other abstractions.
// These messages are meant to wrap and route the messages of the sub-algorithms.
type MsgWrapper struct {
	msgType       byte
	subsystemFunc func(subsystem byte, index int) (GPA, error) // Resolve a subsystem GPA based on its code and index.
}

func NewMsgWrapper(msgType byte, subsystemFunc func(subsystem byte, index int) (GPA, error)) *MsgWrapper {
	return &MsgWrapper{msgType, subsystemFunc}
}

func (w *MsgWrapper) WrapMessage(subsystem byte, index int, msg Message) Message {
	return &WrappingMsg{w.msgType, subsystem, index, msg}
}

func (w *MsgWrapper) WrapMessages(subsystem byte, index int, msgs OutMessages) OutMessages {
	if msgs == nil {
		return nil
	}
	wrapped := NoMessages()
	msgs.MustIterate(func(msg Message) {
		wrapped.Add(w.WrapMessage(subsystem, index, msg))
	})
	return wrapped
}

func (w *MsgWrapper) DelegateInput(subsystem byte, index int, input Input) (GPA, OutMessages, error) {
	sub, err := w.subsystemFunc(subsystem, index)
	if err != nil {
		return nil, nil, err
	}
	return sub, w.WrapMessages(subsystem, index, sub.Input(input)), nil
}

func (w *MsgWrapper) DelegateMessage(msg *WrappingMsg) (GPA, OutMessages, error) {
	sub, err := w.subsystemFunc(msg.Subsystem(), msg.Index())
	if err != nil {
		return nil, nil, err
	}
	return sub, w.WrapMessages(msg.Subsystem(), msg.Index(), sub.Message(msg.Wrapped())), nil
}

func (w *MsgWrapper) UnmarshalMessage(data []byte) (Message, error) {
	rr := rwutil.NewBytesReader(data)
	rr.ReadMessageTypeAndVerify(w.msgType)
	ret := &WrappingMsg{
		msgType:   w.msgType,
		subsystem: rr.ReadByte(),
		index:     int(rr.ReadUint16()),
	}
	wrappedData := rr.ReadBytes()
	if rr.Err != nil {
		return nil, rr.Err
	}

	subGPA, err := w.subsystemFunc(ret.subsystem, ret.index)
	if err != nil {
		return nil, err
	}
	ret.wrapped, err = subGPA.UnmarshalMessage(wrappedData)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// The message that contains another, and its routing info.
type WrappingMsg struct {
	msgType   byte
	subsystem byte
	index     int
	wrapped   Message
}

var _ Message = new(WrappingMsg)

func NewWrappingMsg(msgType, subsystem byte, index int, wrapped Message) *WrappingMsg {
	return &WrappingMsg{msgType: msgType, subsystem: subsystem, index: index, wrapped: wrapped}
}

func (msg *WrappingMsg) Subsystem() byte {
	return msg.subsystem
}

func (msg *WrappingMsg) Index() int {
	return msg.index
}

func (msg *WrappingMsg) Wrapped() Message {
	return msg.wrapped
}

func (msg *WrappingMsg) Recipient() NodeID {
	return msg.wrapped.Recipient()
}

func (msg *WrappingMsg) SetSender(sender NodeID) {
	msg.wrapped.SetSender(sender)
}

func (msg *WrappingMsg) MarshalBinary() ([]byte, error) {
	return rwutil.MarshalBinary(msg)
}

func (msg *WrappingMsg) UnmarshalBinary(data []byte) error {
	// return rwutil.UnmarshalBinary(data, msg)
	panic("this message is un-marshaled by the gpa.MsgWrapper")
}

// note: never called, unfinished concept version
func (msg *WrappingMsg) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadMessageTypeAndVerify(msg.msgType)
	msg.subsystem = rr.ReadByte()
	msg.index = int(rr.ReadUint16())
	// TODO: allocate proper message
	rr.ReadMarshaled(msg.wrapped)
	return rr.Err
}

func (msg *WrappingMsg) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteMessageType(msg.msgType)
	ww.WriteByte(msg.subsystem)
	ww.WriteUint16(uint16(msg.index))
	ww.WriteMarshaled(msg.wrapped)
	return ww.Err
}
