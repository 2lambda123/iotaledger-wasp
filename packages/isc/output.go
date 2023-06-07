// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

var emptyTransactionID = iotago.TransactionID{}

type OutputInfo struct {
	OutputID           iotago.OutputID
	Output             iotago.Output
	TransactionIDSpent iotago.TransactionID
}

func (o *OutputInfo) Consumed() bool {
	return o.TransactionIDSpent != emptyTransactionID
}

func NewOutputInfo(outputID iotago.OutputID, output iotago.Output, transactionIDSpent iotago.TransactionID) *OutputInfo {
	return &OutputInfo{
		OutputID:           outputID,
		Output:             output,
		TransactionIDSpent: transactionIDSpent,
	}
}

func (o *OutputInfo) AliasOutputWithID() *AliasOutputWithID {
	return NewAliasOutputWithID(o.Output.(*iotago.AliasOutput), o.OutputID)
}

type AliasOutputWithID struct {
	outputID    iotago.OutputID
	aliasOutput *iotago.AliasOutput
}

func NewAliasOutputWithID(aliasOutput *iotago.AliasOutput, outputID iotago.OutputID) *AliasOutputWithID {
	return &AliasOutputWithID{
		outputID:    outputID,
		aliasOutput: aliasOutput,
	}
}

func (a *AliasOutputWithID) Bytes() []byte {
	return rwutil.WriterToBytes(a)
}

func (a *AliasOutputWithID) GetAliasOutput() *iotago.AliasOutput {
	return a.aliasOutput
}

func (a *AliasOutputWithID) OutputID() iotago.OutputID {
	return a.outputID
}

func (a *AliasOutputWithID) TransactionID() iotago.TransactionID {
	return a.outputID.TransactionID()
}

func (a *AliasOutputWithID) GetStateIndex() uint32 {
	return a.aliasOutput.StateIndex
}

func (a *AliasOutputWithID) GetStateMetadata() []byte {
	return a.aliasOutput.StateMetadata
}

func (a *AliasOutputWithID) GetStateAddress() iotago.Address {
	return a.aliasOutput.StateController()
}

func (a *AliasOutputWithID) GetAliasID() iotago.AliasID {
	return util.AliasIDFromAliasOutput(a.aliasOutput, a.outputID)
}

func (a *AliasOutputWithID) Equals(other *AliasOutputWithID) bool {
	if a != nil && other == nil {
		return false
	}
	if a.outputID != other.outputID {
		return false
	}
	out1, _ := a.aliasOutput.Serialize(serializer.DeSeriModeNoValidation, nil)
	out2, _ := other.aliasOutput.Serialize(serializer.DeSeriModeNoValidation, nil)
	return bytes.Equal(out1, out2)
}

func (a *AliasOutputWithID) Hash() hashing.HashValue {
	return hashing.HashDataBlake2b(a.Bytes())
}

func (a *AliasOutputWithID) String() string {
	if a == nil {
		return "nil"
	}
	return fmt.Sprintf("AO[si#%v]%v", a.GetStateIndex(), a.outputID.ToHex())
}

func (a *AliasOutputWithID) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(a.outputID[:])
	a.aliasOutput = new(iotago.AliasOutput)
	rr.ReadSerialized(a.aliasOutput)
	return rr.Err
}

func (a *AliasOutputWithID) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(a.outputID[:])
	ww.WriteSerialized(a.aliasOutput)
	return ww.Err
}

func OutputSetToOutputIDs(outputSet iotago.OutputSet) iotago.OutputIDs {
	outputIDs := make(iotago.OutputIDs, len(outputSet))
	i := 0
	for id := range outputSet {
		outputIDs[i] = id
		i++
	}
	return outputIDs
}

func AliasOutputWithIDFromTx(tx *iotago.Transaction, aliasAddr iotago.Address) (*AliasOutputWithID, error) {
	txID, err := tx.ID()
	if err != nil {
		return nil, err
	}

	for index, output := range tx.Essence.Outputs {
		if aliasOutput, ok := output.(*iotago.AliasOutput); ok {
			outputID := iotago.OutputIDFromTransactionIDAndIndex(txID, uint16(index))

			aliasID := aliasOutput.AliasID
			if aliasID.Empty() {
				aliasID = iotago.AliasIDFromOutputID(outputID)
			}

			if aliasID.ToAddress().Equal(aliasAddr) {
				// output found
				return NewAliasOutputWithID(aliasOutput, outputID), nil
			}
		}
	}

	return nil, fmt.Errorf("cannot find alias output for address %v in transaction", aliasAddr.String())
}
