package rootimpl

import (
	"bytes"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventDeploy(ctx isc.Sandbox, progHash hashing.HashValue, name string, description string) {
	w := &bytes.Buffer{}
	_ = util.WriteBytes(w, progHash.Bytes())
	_ = util.WriteString(w, name)
	_ = util.WriteString(w, description)
	ctx.Event("coreroot.deploy", w.Bytes())
}

func eventGrant(ctx isc.Sandbox, deployer isc.AgentID) {
	w := &bytes.Buffer{}
	_ = util.WriteBytes(w, deployer.Bytes())
	ctx.Event("coreroot.grant", w.Bytes())
}

func eventRevoke(ctx isc.Sandbox, deployer isc.AgentID) {
	w := &bytes.Buffer{}
	_ = util.WriteBytes(w, deployer.Bytes())
	ctx.Event("coreroot.revoke", w.Bytes())
}
