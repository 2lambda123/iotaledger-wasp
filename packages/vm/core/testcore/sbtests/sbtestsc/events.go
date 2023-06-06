package sbtestsc

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventCounter(ctx isc.Sandbox, value uint64) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteUint64(value)
	ctx.Event("testcore.counter", w.Bytes())
}

func eventTest(ctx isc.Sandbox) {
	w := new(bytes.Buffer)
	ctx.Event("testcore.test", w.Bytes())
}
