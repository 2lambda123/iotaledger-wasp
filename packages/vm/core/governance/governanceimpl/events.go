package governanceimpl

import (
	"bytes"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventRotate(ctx isc.Sandbox, newAddr iotago.Address, oldAddr iotago.Address) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteN(isc.BytesFromAddress(newAddr))
	ww.WriteN(isc.BytesFromAddress(oldAddr))
	ctx.Event("coregovernance.rotate", w.Bytes())
}
