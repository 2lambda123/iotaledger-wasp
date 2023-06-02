package sbtestsc

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/isc"
)

func eventCounter(ctx isc.Sandbox, value uint64) {
	mu := marshalutil.New()
	mu.WriteUint64(value)
	ctx.Event("testcore.counter", mu.Bytes())
}

func eventTest(ctx isc.Sandbox) {
	mu := marshalutil.New()
	ctx.Event("testcore.test", mu.Bytes())
}
