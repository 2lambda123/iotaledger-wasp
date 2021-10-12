// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package jsonrpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotaledger/wasp/contracts/native/evm/evmlight"
	"github.com/iotaledger/wasp/packages/evm/evmtypes"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

type EVMChain struct {
	backend      ChainBackend
	chainID      int
	contractName string
}

func NewEVMChain(backend ChainBackend, chainID int, contractName string) *EVMChain {
	return &EVMChain{backend, chainID, contractName}
}

func (e *EVMChain) Signer() types.Signer {
	return evmtypes.Signer(big.NewInt(int64(e.chainID)))
}

func (e *EVMChain) GasPerIota() (uint64, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetGasPerIota.Name, nil)
	if err != nil {
		return 0, err
	}
	return codec.DecodeUint64(ret.MustGet(evmlight.FieldResult))
}

func (e *EVMChain) BlockNumber() (*big.Int, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetBlockNumber.Name, nil)
	if err != nil {
		return nil, err
	}

	bal := big.NewInt(0)
	bal.SetBytes(ret.MustGet(evmlight.FieldResult))
	return bal, nil
}

func (e *EVMChain) FeeColor() (colored.Color, error) {
	feeInfo, err := e.backend.CallView(governance.Contract.Name, governance.FuncGetFeeInfo.Name, dict.Dict{
		root.ParamHname: evmlight.Contract.Hname().Bytes(),
	})
	if err != nil {
		return colored.Color{}, err
	}
	return codec.DecodeColor(feeInfo.MustGet(governance.ParamFeeColor))
}

func (e *EVMChain) GasLimitFee(tx *types.Transaction) (colored.Color, uint64, error) {
	gpi, err := e.GasPerIota()
	if err != nil {
		return colored.Color{}, 0, err
	}
	feeColor, err := e.FeeColor()
	if err != nil {
		return colored.Color{}, 0, err
	}
	return feeColor, tx.Gas() / gpi, nil
}

func (e *EVMChain) SendTransaction(tx *types.Transaction) error {
	feeColor, feeAmount, err := e.GasLimitFee(tx)
	if err != nil {
		return err
	}
	fee := colored.NewBalancesForColor(feeColor, feeAmount)
	// deposit fee into sender's on-chain account
	err = e.backend.PostOnLedgerRequest(accounts.Contract.Name, accounts.FuncDeposit.Name, fee, nil)
	if err != nil {
		return err
	}
	txdata, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	// send the Ethereum transaction to the evmchain contract
	return e.backend.PostOffLedgerRequest(e.contractName, evmlight.FuncSendTransaction.Name, fee, dict.Dict{
		evmlight.FieldTransactionData: txdata,
	})
}

func paramsWithOptionalBlockNumber(blockNumber *big.Int, params dict.Dict) dict.Dict {
	ret := params
	if params == nil {
		ret = dict.Dict{}
	}
	if blockNumber != nil {
		ret.Set(evmlight.FieldBlockNumber, blockNumber.Bytes())
	}
	return ret
}

func (e *EVMChain) Balance(address common.Address, blockNumber *big.Int) (*big.Int, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetBalance.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{
		evmlight.FieldAddress: address.Bytes(),
	}))
	if err != nil {
		return nil, err
	}

	bal := big.NewInt(0)
	bal.SetBytes(ret.MustGet(evmlight.FieldResult))
	return bal, nil
}

func (e *EVMChain) Code(address common.Address, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetCode.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{
		evmlight.FieldAddress: address.Bytes(),
	}))
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmlight.FieldResult), nil
}

func (e *EVMChain) BlockByNumber(blockNumber *big.Int) (*types.Block, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetBlockByNumber.Name, paramsWithOptionalBlockNumber(blockNumber, nil))
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmlight.FieldResult) {
		return nil, nil
	}

	block, err := evmtypes.DecodeBlock(ret.MustGet(evmlight.FieldResult))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (e *EVMChain) getTransactionBy(funcName string, args dict.Dict) (tx *types.Transaction, blockHash common.Hash, blockNumber, index uint64, err error) { // nolint:unparam
	var ret dict.Dict
	ret, err = e.backend.CallView(e.contractName, funcName, args)
	if err != nil {
		return
	}

	if !ret.MustHas(evmlight.FieldTransaction) {
		return
	}

	tx, err = evmtypes.DecodeTransaction(ret.MustGet(evmlight.FieldTransaction))
	if err != nil {
		return
	}
	blockHash = common.BytesToHash(ret.MustGet(evmlight.FieldBlockHash))
	blockNumber, err = codec.DecodeUint64(ret.MustGet(evmlight.FieldBlockNumber), 0)
	if err != nil {
		return
	}
	// index is always 0
	return
}

func (e *EVMChain) TransactionByHash(hash common.Hash) (tx *types.Transaction, blockHash common.Hash, blockNumber, index uint64, err error) {
	return e.getTransactionBy(evmlight.FuncGetTransactionByHash.Name, dict.Dict{
		evmlight.FieldTransactionHash: hash.Bytes(),
	})
}

func (e *EVMChain) TransactionByBlockHashAndIndex(hash common.Hash, index uint64) (tx *types.Transaction, blockHash common.Hash, blockNumber, indexRet uint64, err error) {
	if index != 0 {
		// all blocks have 1 tx
		return
	}
	return e.getTransactionBy(evmlight.FuncGetTransactionByBlockHash.Name, dict.Dict{
		evmlight.FieldBlockHash: hash.Bytes(),
	})
}

func (e *EVMChain) TransactionByBlockNumberAndIndex(blockNumber *big.Int, index uint64) (tx *types.Transaction, blockHash common.Hash, blockNumberRet, indexRet uint64, err error) {
	if index != 0 {
		// all blocks have 1 tx
		return
	}
	return e.getTransactionBy(evmlight.FuncGetTransactionByBlockNumber.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{}))
}

func (e *EVMChain) BlockByHash(hash common.Hash) (*types.Block, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetBlockByHash.Name, dict.Dict{
		evmlight.FieldBlockHash: hash.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmlight.FieldResult) {
		return nil, nil
	}

	block, err := evmtypes.DecodeBlock(ret.MustGet(evmlight.FieldResult))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (e *EVMChain) TransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetReceipt.Name, dict.Dict{
		evmlight.FieldTransactionHash: txHash.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	if !ret.MustHas(evmlight.FieldResult) {
		return nil, nil
	}

	receipt, err := evmtypes.DecodeReceiptFull(ret.MustGet(evmlight.FieldResult))
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (e *EVMChain) TransactionCount(address common.Address, blockNumber *big.Int) (uint64, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetNonce.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{
		evmlight.FieldAddress: address.Bytes(),
	}))
	if err != nil {
		return 0, err
	}
	return codec.DecodeUint64(ret.MustGet(evmlight.FieldResult), 0)
}

func (e *EVMChain) CallContract(args ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncCallContract.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{
		evmlight.FieldCallMsg: evmtypes.EncodeCallMsg(args),
	}))
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmlight.FieldResult), nil
}

func (e *EVMChain) EstimateGas(args ethereum.CallMsg) (uint64, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncEstimateGas.Name, dict.Dict{
		evmlight.FieldCallMsg: evmtypes.EncodeCallMsg(args),
	})
	if err != nil {
		return 0, err
	}
	return codec.DecodeUint64(ret.MustGet(evmlight.FieldResult), 0)
}

func (e *EVMChain) StorageAt(address common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetStorage.Name, paramsWithOptionalBlockNumber(blockNumber, dict.Dict{
		evmlight.FieldAddress: address.Bytes(),
		evmlight.FieldKey:     key.Bytes(),
	}))
	if err != nil {
		return nil, err
	}
	return ret.MustGet(evmlight.FieldResult), nil
}

func (e *EVMChain) BlockTransactionCountByHash(blockHash common.Hash) (uint64, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetTransactionCountByBlockHash.Name, dict.Dict{
		evmlight.FieldBlockHash: blockHash.Bytes(),
	})
	if err != nil {
		return 0, err
	}
	return codec.DecodeUint64(ret.MustGet(evmlight.FieldResult), 0)
}

func (e *EVMChain) BlockTransactionCountByNumber(blockNumber *big.Int) (uint64, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetTransactionCountByBlockNumber.Name, paramsWithOptionalBlockNumber(blockNumber, nil))
	if err != nil {
		return 0, err
	}
	return codec.DecodeUint64(ret.MustGet(evmlight.FieldResult), 0)
}

func (e *EVMChain) Logs(q *ethereum.FilterQuery) ([]*types.Log, error) {
	ret, err := e.backend.CallView(e.contractName, evmlight.FuncGetLogs.Name, dict.Dict{
		evmlight.FieldFilterQuery: evmtypes.EncodeFilterQuery(q),
	})
	if err != nil {
		return nil, err
	}
	return evmtypes.DecodeLogs(ret.MustGet(evmlight.FieldResult))
}
