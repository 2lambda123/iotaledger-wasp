package corecontracts

import (
	"errors"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
)

type BlockLog struct {
	vmService interfaces.VMService
}

func NewBlockLog(vmService interfaces.VMService) *BlockLog {
	return &BlockLog{
		vmService: vmService,
	}
}

func (b *BlockLog) GetControlAddresses(chainID isc.ChainID) (*blocklog.ControlAddresses, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewControlAddresses.Hname(), nil)
	if err != nil {
		return nil, err
	}

	par := kvdecoder.New(ret)

	stateAddress, err := par.GetAddress(blocklog.ParamStateControllerAddress)
	if err != nil {
		return nil, err
	}

	governingAddress, err := par.GetAddress(blocklog.ParamGoverningAddress)
	if err != nil {
		return nil, err
	}

	sinceBlockIndex, err := par.GetUint32(blocklog.ParamBlockIndex)
	if err != nil {
		return nil, err
	}

	controlAddresses := &blocklog.ControlAddresses{
		StateAddress:     stateAddress,
		GoverningAddress: governingAddress,
		SinceBlockIndex:  sinceBlockIndex,
	}

	return controlAddresses, nil
}

func handleBlockInfo(info dict.Dict) (*blocklog.BlockInfo, error) {
	resultDecoder := kvdecoder.New(info)

	blockInfoBin, err := resultDecoder.GetBytes(blocklog.ParamBlockInfo)
	if err != nil {
		return nil, err
	}

	blockInfo, err := blocklog.BlockInfoFromBytes(blockInfoBin)
	if err != nil {
		return nil, err
	}

	return blockInfo, nil
}

func (b *BlockLog) GetLatestBlockInfo(chainID isc.ChainID) (*blocklog.BlockInfo, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetBlockInfo.Hname(), nil)
	if err != nil {
		return nil, err
	}

	return handleBlockInfo(ret)
}

func (b *BlockLog) GetBlockInfo(chainID isc.ChainID, blockIndex uint32) (*blocklog.BlockInfo, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetBlockInfo.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamBlockIndex: blockIndex,
	}))
	if err != nil {
		return nil, err
	}

	return handleBlockInfo(ret)
}

func handleRequestIDs(requestIDsDict dict.Dict) ([]isc.RequestID, error) {
	requestIDCollection := collections.NewArray16ReadOnly(requestIDsDict, blocklog.ParamRequestID)
	requestIDsCount := requestIDCollection.Len()

	requestIDs := make([]isc.RequestID, requestIDsCount)

	for i := range requestIDs {
		reqIDBin := requestIDCollection.GetAt(uint16(i))

		var err error
		requestIDs[i], err = isc.RequestIDFromBytes(reqIDBin)
		if err != nil {
			return nil, err
		}
	}

	return requestIDs, nil
}

func (b *BlockLog) GetRequestIDsForLatestBlock(chainID isc.ChainID) ([]isc.RequestID, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetRequestIDsForBlock.Hname(), nil)
	if err != nil {
		return nil, err
	}

	return handleRequestIDs(ret)
}

func (b *BlockLog) GetRequestIDsForBlock(chainID isc.ChainID, blockIndex uint32) ([]isc.RequestID, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetRequestIDsForBlock.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamBlockIndex: blockIndex,
	}))
	if err != nil {
		return nil, err
	}

	return handleRequestIDs(ret)
}

var ErrNoRecord = errors.New("no request record")

func (b *BlockLog) GetRequestReceipt(chainID isc.ChainID, requestID isc.RequestID) (*blocklog.RequestReceipt, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetRequestReceipt.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamRequestID: requestID,
	}))
	if err != nil {
		return nil, err
	}

	resultDecoder := kvdecoder.New(ret)
	binRec, err := resultDecoder.GetBytes(blocklog.ParamRequestRecord)
	if err != nil || binRec == nil {
		return nil, ErrNoRecord
	}

	requestReceipt, err := blocklog.RequestReceiptFromBytes(binRec)
	if err != nil {
		return nil, err
	}

	requestReceipt.BlockIndex, err = resultDecoder.GetUint32(blocklog.ParamBlockIndex)
	if err != nil {
		return nil, err
	}

	requestReceipt.RequestIndex, err = resultDecoder.GetUint16(blocklog.ParamRequestIndex)
	if err != nil {
		return nil, err
	}

	return requestReceipt, err
}

func (b *BlockLog) GetRequestReceiptsForBlock(chainID isc.ChainID, blockIndex uint32) ([]*blocklog.RequestReceipt, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetRequestReceiptsForBlock.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamBlockIndex: blockIndex,
	}))
	if err != nil {
		return nil, err
	}

	returnedBlockIndex, err := codec.DecodeUint32(ret.Get(blocklog.ParamBlockIndex))
	if err != nil {
		return nil, err
	}

	requestRecordCollection := collections.NewArray16ReadOnly(ret, blocklog.ParamRequestRecord)
	requestRecordCount := requestRecordCollection.Len()

	requestReceipts := make([]*blocklog.RequestReceipt, requestRecordCount)

	for i := range requestReceipts {
		data := requestRecordCollection.GetAt(uint16(i))
		requestReceipts[i], err = blocklog.RequestReceiptFromBytes(data)
		if err != nil {
			return nil, err
		}
		requestReceipts[i].WithBlockData(returnedBlockIndex, uint16(i))
	}

	return requestReceipts, nil
}

func (b *BlockLog) IsRequestProcessed(chainID isc.ChainID, requestID isc.RequestID) (bool, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewIsRequestProcessed.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamRequestID: requestID,
	}))
	if err != nil {
		return false, err
	}

	resultDecoder := kvdecoder.New(ret)
	isProcessed, err := resultDecoder.GetBool(blocklog.ParamRequestProcessed)
	if err != nil {
		return false, err
	}

	return isProcessed, nil
}

func eventsFromViewResult(viewResult dict.Dict) ([]string, error) {
	eventCollection := collections.NewArray16ReadOnly(viewResult, blocklog.ParamEvent)
	eventCount := eventCollection.Len()

	events := make([]string, eventCount)
	for i := range events {
		data := eventCollection.GetAt(uint16(i))
		events[i] = string(data)
	}

	return events, nil
}

func (b *BlockLog) GetEventsForRequest(chainID isc.ChainID, requestID isc.RequestID) ([]string, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetEventsForRequest.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamRequestID: requestID,
	}))
	if err != nil {
		return nil, err
	}

	return eventsFromViewResult(ret)
}

func (b *BlockLog) GetEventsForBlock(chainID isc.ChainID, blockIndex uint32) ([]string, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetEventsForBlock.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamBlockIndex: blockIndex,
	}))
	if err != nil {
		return nil, err
	}

	return eventsFromViewResult(ret)
}

func (b *BlockLog) GetEventsForContract(chainID isc.ChainID, contractHname isc.Hname) ([]string, error) {
	ret, err := b.vmService.CallViewByChainID(chainID, blocklog.Contract.Hname(), blocklog.ViewGetEventsForContract.Hname(), codec.MakeDict(map[string]interface{}{
		blocklog.ParamContractHname: contractHname,
	}))
	if err != nil {
		return nil, err
	}

	return eventsFromViewResult(ret)
}
