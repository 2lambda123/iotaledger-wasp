// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/chain/cons"
	"github.com/iotaledger/wasp/packages/gpa"
)

type inputConsensusOutputDone struct {
	logIndex LogIndex
	result   *cons.Result
}

// This message is internal one, but should be sent by other components (e.g. consensus or the chain).
func NewInputConsensusOutputDone(
	logIndex LogIndex,
	result *cons.Result,
) gpa.Input {
	return &inputConsensusOutputDone{
		logIndex: logIndex,
		result:   result,
	}
}

func (inp *inputConsensusOutputDone) String() string {
	return fmt.Sprintf(
		"{cmtLog.inputConsensusOutputDone, logIndex=%v, %v}",
		inp.logIndex, inp.result,
	)
}
