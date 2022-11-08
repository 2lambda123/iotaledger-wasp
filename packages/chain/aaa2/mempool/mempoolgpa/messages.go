package mempoolgpa

import (
	"time"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type MempoolMessages struct {
	gpa.OutMessages
	SendNperIteration int
	SendInterval      time.Duration
	ShouldStopSending func() bool
	Staggered         bool
}

func NewStaggeredMessages(
	msgs gpa.OutMessages,
	sendN int,
	sendEvery time.Duration,
	shouldStopSending func() bool,
) *MempoolMessages {
	return &MempoolMessages{
		OutMessages:       msgs,
		SendNperIteration: sendN,
		SendInterval:      sendEvery,
		ShouldStopSending: shouldStopSending,
		Staggered:         true,
	}
}

func NoMessages() *MempoolMessages {
	return &MempoolMessages{
		OutMessages: gpa.NoMessages(),
		Staggered:   false,
	}
}

func SingleMessage(msg gpa.Message) *MempoolMessages {
	return &MempoolMessages{
		OutMessages: gpa.NoMessages().Add(msg),
		Staggered:   false,
	}
}

// ----------------------------------------------------------------

const (
	msgTypeShareRequest byte = iota
	msgTypeMissingRequest
)

// share offledger req
type msgShareRequest struct {
	gpa.BasicMessage
	req             isc.Request
	shouldPropagate bool // indecates whether the message should be shared (false when responding to a "missing request message")
}

var _ gpa.Message = &msgShareRequest{}

func newMsgShareRequest(req isc.Request, shouldPropagate bool, receipient gpa.NodeID) *msgShareRequest {
	return &msgShareRequest{
		BasicMessage:    gpa.NewBasicMessage(receipient),
		req:             req,
		shouldPropagate: shouldPropagate,
	}
}

func (msg *msgShareRequest) MarshalBinary() (data []byte, err error) {
	ret := []byte{msgTypeMissingRequest}
	ret = append(ret, msg.req.Bytes()...)
	return ret, nil
}

func (msg *msgShareRequest) UnmarshalBinary(data []byte) (err error) {
	msg.req, err = isc.NewRequestFromBytes(data)
	return err
}

// ----------------------------------------------------------------

// ask for missing req
type msgMissingRequest struct {
	gpa.BasicMessage
	ref *isc.RequestRef
}

var _ gpa.Message = &msgMissingRequest{}

func newMsgMissingRequest(ref *isc.RequestRef, receipient gpa.NodeID) *msgMissingRequest {
	return &msgMissingRequest{
		BasicMessage: gpa.NewBasicMessage(receipient),
		ref:          ref,
	}
}

func (msg *msgMissingRequest) MarshalBinary() (data []byte, err error) {
	ret := []byte{msgTypeMissingRequest}
	ret = append(ret, msg.ref.Bytes()...)
	return ret, nil
}

func (msg *msgMissingRequest) UnmarshalBinary(data []byte) (err error) {
	msg.ref, err = isc.RequestRefFromBytes(data[1:])
	return err
}
