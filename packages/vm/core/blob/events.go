package blob

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventStore(ctx isc.Sandbox, blobHash hashing.HashValue) {
	mu := marshalutil.New()
	util.WriteBytesMu(mu, blobHash.Bytes())
	ctx.Event("coreblob.store", mu.Bytes())
}
