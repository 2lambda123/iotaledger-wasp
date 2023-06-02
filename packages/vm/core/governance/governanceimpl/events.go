package governanceimpl

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventRotate(ctx isc.Sandbox, newAddr iotago.Address, oldAddr iotago.Address) {
	mu := marshalutil.New()
	util.WriteBytesMu(mu, isc.BytesFromAddress(newAddr))
	util.WriteBytesMu(mu, isc.BytesFromAddress(oldAddr))
	ctx.Event("coregovernance.rotate", mu.Bytes())
}
