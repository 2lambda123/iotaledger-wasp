package inccounter

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventCounter(ctx isc.Sandbox, val int64) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteInt64(val)
	ctx.Event("inccounter.counter", w.Bytes())
}
