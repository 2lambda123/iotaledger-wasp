package chain

import (
	"net/http"

	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/controllers/controllerutils"
	"github.com/iotaledger/wasp/packages/webapi/models"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/wasp/packages/isc"
)

func (c *Controller) executeCallView(e echo.Context) error {
	controllerutils.SetOperation(e, "call_view")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService)
	if err != nil {
		return err
	}

	var callViewRequest models.ContractCallViewRequest
	if err = e.Bind(&callViewRequest); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	if !c.chainService.HasChain(chainID) {
		return apierrors.ChainNotFoundError(chainID.String())
	}

	// Get contract and function. The request model supports HName and common string names. HNames are preferred.
	var contractHName isc.Hname
	var functionHName isc.Hname

	if callViewRequest.ContractHName == "" {
		contractHName = isc.Hn(callViewRequest.ContractName)
	} else {
		contractHName, err = isc.HnameFromHexString(callViewRequest.ContractHName)
		if err != nil {
			return apierrors.InvalidPropertyError("contractHName", err)
		}
	}

	if callViewRequest.FunctionHName == "" {
		functionHName = isc.Hn(callViewRequest.FunctionName)
	} else {
		functionHName, err = isc.HnameFromHexString(callViewRequest.FunctionHName)
		if err != nil {
			return apierrors.InvalidPropertyError("contractHName", err)
		}
	}

	args, err := dict.FromJSONDict(callViewRequest.Arguments)
	if err != nil {
		return apierrors.InvalidPropertyError("arguments", err)
	}

	result, err := c.vmService.CallViewByChainID(chainID, contractHName, functionHName, args)
	if err != nil {
		return apierrors.ContractExecutionError(err)
	}

	return e.JSON(http.StatusOK, result)
}
