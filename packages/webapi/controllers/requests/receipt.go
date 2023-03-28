package requests

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/corecontracts"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
	"github.com/iotaledger/wasp/packages/webapi/models"
	"github.com/iotaledger/wasp/packages/webapi/params"
)

// TODO this should reuse the code from webapi/controllers/corecontracts/blocklog getRequestReceipt
func (c *Controller) getReceipt(e echo.Context) error {
	chainID, err := params.DecodeChainID(e)
	if err != nil {
		return err
	}

	requestID, err := params.DecodeRequestID(e)
	if err != nil {
		return err
	}

	receipt, vmError, err := c.vmService.GetReceipt(chainID, requestID)
	if err != nil {
		if errors.Is(err, corecontracts.ErrNoRecord) {
			return apierrors.NoRecordFoundErrror(err)
		}
		if errors.Is(err, interfaces.ErrChainNotFound) {
			return apierrors.ChainNotFoundError(chainID.String())
		}
		return apierrors.ReceiptError(err)
	}

	mappedReceipt := models.MapReceiptResponse(receipt, vmError)

	return e.JSON(http.StatusOK, mappedReceipt)
}
