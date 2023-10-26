package chainutil

import (
	"time"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
)

func SimulateRequest(
	ch chain.ChainCore,
	req isc.Request,
	estimateGas bool,
) (*blocklog.RequestReceipt, error) {
	aliasOutput, err := ch.LatestAliasOutput(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, ierrors.Errorf("could not get latest AliasOutput: %w", err)
	}
	res, err := runISCRequest(ch, aliasOutput, time.Now(), req, estimateGas)
	if err != nil {
		return nil, err
	}
	return res.Receipt, nil
}
