package cmt_log

import (
	"fmt"

	"github.com/iotaledger/hive.go/log"
	"github.com/iotaledger/wasp/packages/chain/cons"
)

// We can provide input to the next consensus when
//   - there is base output determined or block to sign.
//   - the log index is agreed.
//   - the minimal delay has passed from the previous consensus.
//
// TODO: delays should be considered only for the consensus rounds producing new blocks.
type VarOutput interface {
	StatusString() string // Summary of the internal state.
	Value() *Output
	LogIndexAgreed(li LogIndex)
	ConsInputChanged(consInput cons.Input)
	CanPropose()
	Suspended(suspended bool)
}

type varOutputImpl struct {
	candidateLI LogIndex
	consInput   cons.Input
	canPropose  bool
	suspended   bool
	outValue    *Output
	persistUsed func(li LogIndex)
	log         log.Logger
}

func NewVarOutput(persistUsed func(li LogIndex), log log.Logger) VarOutput {
	return &varOutputImpl{
		candidateLI: NilLogIndex(),
		consInput:   nil,
		canPropose:  true,
		suspended:   false,
		outValue:    nil,
		persistUsed: persistUsed,
		log:         log,
	}
}

func (vo *varOutputImpl) StatusString() string {
	return fmt.Sprintf(
		"{varOutput: output=%v, candidate{li=%v, consInput=%v}, canPropose=%v, suspended=%v}",
		vo.outValue, vo.candidateLI, vo.consInput, vo.canPropose, vo.suspended,
	)
}

func (vo *varOutputImpl) Value() *Output {
	if vo.outValue == nil || vo.suspended {
		return nil // Untyped nil.
	}
	return vo.outValue
}

func (vo *varOutputImpl) LogIndexAgreed(li LogIndex) {
	vo.candidateLI = li
	vo.tryOutput()
}

func (vo *varOutputImpl) ConsInputChanged(consInput cons.Input) {
	vo.consInput = consInput
	vo.tryOutput()
}

func (vo *varOutputImpl) CanPropose() {
	vo.canPropose = true
	vo.tryOutput()
}

func (vo *varOutputImpl) Suspended(suspended bool) {
	if vo.suspended && !suspended {
		vo.log.LogInfof("Committee resumed.")
	}
	if !vo.suspended && suspended {
		vo.log.LogInfof("Committee suspended.")
	}
	vo.suspended = suspended
}

func (vo *varOutputImpl) tryOutput() {
	if vo.candidateLI.IsNil() || vo.consInput == nil || !vo.canPropose {
		// Keep output unchanged.
		return
	}
	//
	// Output the new data.
	vo.persistUsed(vo.candidateLI)
	vo.outValue = makeOutput(vo.candidateLI, vo.consInput)
	vo.log.LogInfof("⊪ Output %v", vo.outValue)
	vo.canPropose = false
	vo.candidateLI = NilLogIndex()
}
