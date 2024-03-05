// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
)

// That's the main/initial input for the consensus.
type inputProposal struct {
	params Input
}

func NewInputProposal(params Input) gpa.Input {
	return &inputProposal{params: params}
}

func (ip *inputProposal) String() string {
	return fmt.Sprintf("{cons.inputProposal: %v}", ip.params)
}
