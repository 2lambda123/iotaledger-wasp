package cons

import (
	"fmt"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
)

type inputNodeConnBlockTipSet struct {
	strongParents iotago.BlockIDs
}

func NewInputNodeConnBlockTipSet(strongParents iotago.BlockIDs) gpa.Input {
	return &inputNodeConnBlockTipSet{strongParents: strongParents}
}

func (inp *inputNodeConnBlockTipSet) String() string {
	return fmt.Sprintf("{cons.inputNodeConnBlockTipSet, |strongParents|=%v}", len(inp.strongParents))
}
