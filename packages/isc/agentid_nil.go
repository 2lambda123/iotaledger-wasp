package isc

import (
	"io"

	"github.com/iotaledger/wasp/packages/util"
)

const nilAgentIDString = "-"

type NilAgentID struct{}

var _ AgentID = &NilAgentID{}

func (a *NilAgentID) Kind() AgentIDKind {
	return AgentIDKindNil
}

func (a *NilAgentID) Bytes() []byte {
	return util.WriterToBytes(a)
}

func (a *NilAgentID) String() string {
	return nilAgentIDString
}

func (a *NilAgentID) Equals(other AgentID) bool {
	if other == nil {
		return false
	}
	return other.Kind() == a.Kind()
}

// note: local read(), no need to read type byte
func (a *NilAgentID) read(rr *util.Reader) {
}

func (a *NilAgentID) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteUint8(uint8(a.Kind()))
	return ww.Err
}
