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

func NewEvent(event []byte) (ret *Event, err error) {
	r := bytes.NewBuffer(event)
	var contractID uint32
	err = util.ReadUint32(r, &contractID)
	if err != nil {
		return nil, err
	}
	ret = &Event{ContractID: Hname(contractID)}
	ret.Topic, err = util.ReadString(r)
	if err != nil {
		return nil, err
	}
	err = util.ReadUint64(r, &ret.Timestamp)
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
	w := &bytes.Buffer{}
	_ = util.WriteUint32(w, uint32(e.ContractID))
	_ = util.WriteString(w, e.Topic)
	_ = util.WriteUint64(w, e.Timestamp)
	_ = util.WriteBytes(w, e.Payload)
	return w.Bytes()
}
