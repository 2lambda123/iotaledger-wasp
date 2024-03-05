package chain

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/iota.go/v4/hexutil"
	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/controllers/controllerutils"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
	"github.com/iotaledger/wasp/packages/webapi/models"
	"github.com/iotaledger/wasp/packages/webapi/params"
	"github.com/iotaledger/wasp/packages/webapi/services"
)

func (c *Controller) getCommitteeInfo(e echo.Context) error {
	controllerutils.SetOperation(e, "get_committee_info")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService, c.l1API)
	if err != nil {
		return err
	}

	chain, err := c.chainService.GetChainInfoByChainID(chainID, "")
	if err != nil {
		return apierrors.ChainNotFoundError(chainID.Bech32(c.l1API.ProtocolParameters().Bech32HRP()))
	}

	chainNodeInfo, err := c.committeeService.GetCommitteeInfo(chainID)
	if err != nil {
		if errors.Is(err, services.ErrNotInCommittee) {
			return e.JSON(http.StatusOK, models.CommitteeInfoResponse{})
		}
		return err
	}

	chainInfo := models.CommitteeInfoResponse{
		ChainID:        chainID.Bech32(c.l1API.ProtocolParameters().Bech32HRP()),
		Active:         chain.IsActive,
		StateAddress:   chainNodeInfo.Address.String(),
		CommitteeNodes: models.MapCommitteeNodes(chainNodeInfo.CommitteeNodes),
		AccessNodes:    models.MapCommitteeNodes(chainNodeInfo.AccessNodes),
		CandidateNodes: models.MapCommitteeNodes(chainNodeInfo.CandidateNodes),
	}

	return e.JSON(http.StatusOK, chainInfo)
}

func (c *Controller) getChainInfo(e echo.Context) error {
	controllerutils.SetOperation(e, "get_chain_info")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService, c.l1API)
	if err != nil {
		return err
	}

	chainInfo, err := c.chainService.GetChainInfoByChainID(chainID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if errors.Is(err, interfaces.ErrChainNotFound) {
		return e.NoContent(http.StatusNotFound)
	} else if err != nil {
		return err
	}

	evmChainID := uint16(0)
	if chainInfo.IsActive {
		evmChainID, err = c.chainService.GetEVMChainID(chainID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return err
		}
	}

	chainInfoResponse := models.MapChainInfoResponse(chainInfo, evmChainID, c.l1API)

	return e.JSON(http.StatusOK, chainInfoResponse)
}

func (c *Controller) getChainList(e echo.Context) error {
	chainIDs, err := c.chainService.GetAllChainIDs()
	c.logger.LogInfof("After allChainIDS %v", err)
	if err != nil {
		return err
	}

	chainList := make([]models.ChainInfoResponse, 0)

	for _, chainID := range chainIDs {
		chainInfo, err := c.chainService.GetChainInfoByChainID(chainID, "")
		c.logger.LogInfof("getchaininfo %v", err)

		if errors.Is(err, interfaces.ErrChainNotFound) {
			// TODO: Validate this logic here. Is it possible to still get more chain info?
			chainList = append(chainList, models.ChainInfoResponse{
				IsActive: false,
				ChainID:  chainID.Bech32(c.l1API.ProtocolParameters().Bech32HRP()),
			})
			continue
		} else if err != nil {
			return err
		}

		evmChainID := uint16(0)
		if chainInfo.IsActive {
			evmChainID, err = c.chainService.GetEVMChainID(chainID, "")
			c.logger.LogInfof("getevmchainid %v", err)

			if err != nil {
				return err
			}
		}

		chainInfoResponse := models.MapChainInfoResponse(chainInfo, evmChainID, c.l1API)
		c.logger.LogInfof("mapchaininfo %v", err)

		chainList = append(chainList, chainInfoResponse)
	}

	return e.JSON(http.StatusOK, chainList)
}

func (c *Controller) getState(e echo.Context) error {
	controllerutils.SetOperation(e, "get_state")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService, c.l1API)
	if err != nil {
		return err
	}

	stateKey, err := hexutil.DecodeHex(e.Param(params.ParamStateKey))
	if err != nil {
		return apierrors.InvalidPropertyError(params.ParamStateKey, err)
	}

	state, err := c.chainService.GetState(chainID, stateKey)
	if err != nil {
		panic(err)
	}

	response := models.StateResponse{
		State: hexutil.EncodeHex(state),
	}

	return e.JSON(http.StatusOK, response)
}
