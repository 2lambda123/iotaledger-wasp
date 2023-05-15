package governanceimpl

import (
	"bytes"
	"fmt"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
)

var _ isc.Event = &RotateStateControllerEvent{}

type RotateStateControllerEvent struct {
	Timestamp              uint64
	NewStateControllerAddr iotago.Address
	StoredStateController  iotago.Address
}

func (e *RotateStateControllerEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, governance.FuncRotateStateController.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncRotateStateController.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *RotateStateControllerEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	// TODO should use byte instead of string
	if err := util.WriteString16(&w, e.NewStateControllerAddr.String()); err != nil {
		panic(fmt.Errorf("failed to write event.NewStateControllerAddr: %w", err))
	}
	if err := util.WriteString16(&w, e.StoredStateController.String()); err != nil {
		panic(fmt.Errorf("failed to write event.StoredStateController: %w", err))
	}
	return w.Bytes()
}

func (e *RotateStateControllerEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if !bytes.Equal(topic, isc.DecodeEventTopic(e)) {
		panic("decode by unmatched event type")
	}
	if err = util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	str, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.NewStateControllerAddr: %w", err))
	}
	_, e.NewStateControllerAddr, err = iotago.ParseBech32(str)
	if err != nil {
		panic(fmt.Errorf("failed to decode NewStateControllerAddr: %w", err))
	}
	str, err = util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.StoredStateController: %w", err))
	}
	_, e.StoredStateController, err = iotago.ParseBech32(str)
	if err != nil {
		panic(fmt.Errorf("failed to decode StoredStateController: %w", err))
	}
}
