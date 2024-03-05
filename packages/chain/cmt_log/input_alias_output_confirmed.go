// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type inputAnchorOutputConfirmed struct {
	confirmedOutputs *isc.ChainOutputs
}

func NewInputAnchorOutputConfirmed(
	confirmedOutputs *isc.ChainOutputs,
) gpa.Input {
	return &inputAnchorOutputConfirmed{
		confirmedOutputs: confirmedOutputs,
	}
}

func (inp *inputAnchorOutputConfirmed) String() string {
	return fmt.Sprintf("{cmtLog.inputAnchorOutputConfirmed, %v}", inp.confirmedOutputs)
}
