package chainclient

import (
	"context"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/wasp/packages/isc"
)

// CheckRequestResult fetches the receipt for the given request ID, and returns
// an error indicating whether the request was processed successfully.
func (c *Client) CheckRequestResult(ctx context.Context, reqID isc.RequestID) error {
	receipt, _, err := c.WaspClient.CorecontractsApi.BlocklogGetRequestReceipt(ctx, c.ChainID.String(), reqID.String()).Execute()
	if err != nil {
		return ierrors.New("could not fetch receipt for request: not found in blocklog")
	}

	if receipt.ErrorMessage != nil {
		return ierrors.Errorf("the request was rejected: %v", receipt.ErrorMessage)
	}

	return nil
}
