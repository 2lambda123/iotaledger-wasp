package accounts

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &FoundryCreateNewEvent{}

type FoundryCreateNewEvent struct {
	SerialNumber uint32
}

func (e *FoundryCreateNewEvent) Encode() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncFoundryCreateNew.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncFoundryCreateNew.Hname(): %w", err))
	}
	if err := util.WriteUint32(&w, e.SerialNumber); err != nil {
		panic(fmt.Errorf("failed to write event.SerialNumber: %w", err))
	}
	return w.Bytes()
}
