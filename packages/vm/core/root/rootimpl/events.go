package rootimpl

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

func eventDeploy(ctx isc.Sandbox, progHash hashing.HashValue, name string, description string) {
	mu := marshalutil.New()
	util.MarshalBytes(mu, progHash.Bytes())
	util.MarshalString(mu, name)
	util.MarshalString(mu, description)
	ctx.Event("coreroot.deploy", mu.Bytes())
}

func eventGrant(ctx isc.Sandbox, deployer isc.AgentID) {
	mu := marshalutil.New()
	util.MarshalBytes(mu, deployer.Bytes())
	ctx.Event("coreroot.grant", mu.Bytes())
}

func eventRevoke(ctx isc.Sandbox, deployer isc.AgentID) {
	mu := marshalutil.New()
	util.MarshalBytes(mu, deployer.Bytes())
	ctx.Event("coreroot.revoke", mu.Bytes())
}
