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
	rr := util.NewReader(r)

	ret := &Event{}
	err := ret.ContractID.Read(r)
	if err != nil {
		return nil, err
	}
	ret.Topic = rr.ReadString()
	ret.Timestamp = rr.ReadUint64()
	ret.Payload = rr.ReadBytes()
	return ret, nil
}

func (e *Event) Bytes() []byte {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	e.ContractID.Write(w)
	ww.WriteString(e.Topic)
	ww.WriteUint64(e.Timestamp)
	ww.WriteBytes(e.Payload)
	return w.Bytes()
}
