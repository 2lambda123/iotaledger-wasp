package blob

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventStore(ctx isc.Sandbox, blobHash hashing.HashValue) {
	w := &bytes.Buffer{}
	_ = util.WriteBytes(w, blobHash.Bytes())
	ctx.Event("coreblob.store", w.Bytes())
}
