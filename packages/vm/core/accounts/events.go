package accounts

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/isc"
)

func eventFoundryCreated(ctx isc.Sandbox, foundrySN uint32) {
	mu := marshalutil.New()
	mu.WriteUint32(foundrySN)
	ctx.Event("coreaccounts.foundryCreated", mu.Bytes())
}

func eventFoundryDestroyed(ctx isc.Sandbox, foundrySN uint32) {
	mu := marshalutil.New()
	mu.WriteUint32(foundrySN)
	ctx.Event("coreaccounts.foundryDestroyed", mu.Bytes())
}

func eventFoundryModified(ctx isc.Sandbox, foundrySN uint32) {
	mu := marshalutil.New()
	mu.WriteUint32(foundrySN)
	ctx.Event("coreaccounts.foundryModified", mu.Bytes())
}
