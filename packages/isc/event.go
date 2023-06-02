package isc

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/util"
)

type Event struct {
	ContractID Hname  `json:"contractID"`
	Payload    []byte `json:"payload"`
	Topic      string `json:"topic"`
	Timestamp  uint64 `json:"timestamp"`
}

func NewEvent(event []byte) (ret *Event, err error) {
	mu := marshalutil.New(event)
	ret = &Event{}
	ret.ContractID, err = HnameFromMarshalUtil(mu)
	if err != nil {
		return nil, err
	}
	ret.Topic, err = util.ReadStringMu(mu)
	if err != nil {
		return nil, err
	}
	ret.Timestamp, err = mu.ReadUint64()
	if err != nil {
		return nil, err
	}
	ret.Payload, err = util.ReadBytesMu(mu)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *Event) Bytes() []byte {
	mu := marshalutil.New()
	mu.WriteUint32(uint32(e.ContractID))
	util.WriteStringMu(mu, e.Topic)
	mu.WriteUint64(e.Timestamp)
	util.WriteBytesMu(mu, e.Payload)
	return mu.Bytes()
}
