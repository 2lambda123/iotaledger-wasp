package governanceimpl

import (
	"bytes"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventRotate(ctx isc.Sandbox, newAddr iotago.Address, oldAddr iotago.Address) {
	w := &bytes.Buffer{}
	_ = util.WriteBytes(w, isc.BytesFromAddress(newAddr))
	_ = util.WriteBytes(w, isc.BytesFromAddress(oldAddr))
	ctx.Event("coregovernance.rotate", w.Bytes())
}
