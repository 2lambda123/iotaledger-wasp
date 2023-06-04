package isc

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/util"
)

type Event struct {
	ContractID Hname  `json:"contractID"`
	Payload    []byte `json:"payload"`
	Topic      string `json:"topic"`
	Timestamp  uint64 `json:"timestamp"`
}

func NewEvent(event []byte) (*Event, error) {
	r := bytes.NewBuffer(event)
	ret := &Event{}
	err := ret.ContractID.Read(r)
	if err != nil {
		return nil, err
	}
	ret.Topic, err = util.ReadString(r)
	if err != nil {
		return nil, err
	}
	ret.Timestamp, err = util.ReadUint64(r)
	if err != nil {
		return nil, err
	}
	ret.Payload, err = util.ReadBytes(r)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *Event) Bytes() []byte {
	w := new(bytes.Buffer)
	_ = e.ContractID.Write(w)
	_ = util.WriteString(w, e.Topic)
	_ = util.WriteUint64(w, e.Timestamp)
	_ = util.WriteBytes(w, e.Payload)
	return w.Bytes()
}
