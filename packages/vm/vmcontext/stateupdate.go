package vmcontext

import (
	"bytes"
	"io"

	"github.com/iotaledger/wasp/packages/kv/buffered"
)

type StateUpdate struct {
	Mutations *buffered.Mutations
}

// NewStateUpdate creates a state update with timestamp mutation, if provided
func NewStateUpdate() *StateUpdate {
	return &StateUpdate{
		Mutations: buffered.NewMutations(),
	}
}

func (su *StateUpdate) Clone() *StateUpdate {
	return &StateUpdate{Mutations: su.Mutations.Clone()}
}

func (su *StateUpdate) Bytes() []byte {
	w := new(bytes.Buffer)
	_ = su.Write(w)
	return w.Bytes()
}

func (su *StateUpdate) Write(w io.Writer) error {
	return su.Mutations.Write(w)
}

func (su *StateUpdate) Read(r io.Reader) error {
	return su.Mutations.Read(r)
}
