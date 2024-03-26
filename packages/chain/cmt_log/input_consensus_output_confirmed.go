// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type inputConsensusOutputConfirmed struct {
	chainOutputs *isc.ChainOutputs
	logIndex     LogIndex
}

func NewInputConsensusOutputConfirmed(chainOutputs *isc.ChainOutputs, logIndex LogIndex) gpa.Input {
	return &inputConsensusOutputConfirmed{
		chainOutputs: chainOutputs,
		logIndex:     logIndex,
	}
}

func (inp *inputConsensusOutputConfirmed) String() string {
	return fmt.Sprintf("{cmtLog.inputConsensusOutputConfirmed, %v, li=%v}", inp.chainOutputs, inp.logIndex)
}
