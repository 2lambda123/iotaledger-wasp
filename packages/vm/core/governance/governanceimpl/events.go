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
	// TODO should use byte instead of string
	if err := util.WriteString16(&w, e.NewStateControllerAddr.String()); err != nil {
		panic(fmt.Errorf("failed to write event.NewStateControllerAddr: %w", err))
	}
	if err := util.WriteString16(&w, e.StoredStateController.String()); err != nil {
		panic(fmt.Errorf("failed to write event.StoredStateController: %w", err))
	}
	return w.Bytes()
}

func (e *RotateStateControllerEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
}
