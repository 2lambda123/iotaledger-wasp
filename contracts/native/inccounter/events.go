package inccounter

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/isc"
)

func eventCounter(ctx isc.Sandbox, val int64) {
	mu := marshalutil.New()
	mu.WriteInt64(val)
	ctx.Event("inccounter.counter", mu.Bytes())
}
