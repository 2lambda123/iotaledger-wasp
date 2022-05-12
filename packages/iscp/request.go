package iscp

import (
	"time"

	"github.com/iotaledger/hive.go/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

func UTXOInputFromMarshalUtil(mu *marshalutil.MarshalUtil) (*iotago.UTXOInput, error) {
	data, err := mu.ReadBytes(iotago.OutputIDLength)
	if err != nil {
		return nil, err
	}
	id, err := DecodeOutputID(data)
	if err != nil {
		return nil, err
	}
	return id.UTXOInput(), nil
}

func UTXOInputToMarshalUtil(id *iotago.UTXOInput, mu *marshalutil.MarshalUtil) {
	mu.WriteBytes(EncodeOutputID(id.ID()))
}

// Request wraps any data which can be potentially be interpreted as a request
type Request interface {
	Calldata

	IsOffLedger() bool
	AsOffLedger() AsOffLedger
	AsOnLedger() AsOnLedger

	Bytes() []byte
	String() string
}

type TimeData struct {
	MilestoneIndex uint32
	Time           time.Time
}

type Calldata interface {
	ID() RequestID
	Params() dict.Dict
	SenderAccount() AgentID // returns nil if sender address is not available
	SenderAddress() iotago.Address
	CallTarget() CallTarget
	TargetAddress() iotago.Address   // TODO implement properly. Target depends on time assumptions and UTXO type
	FungibleTokens() *FungibleTokens // attached assets for the UTXO request, nil for off-ledger. All goes to sender
	NFT() *NFT                       // Not nil if the request is an NFT request
	Allowance() *Allowance           // transfer of assets to the smart contract. Debited from sender account
	GasBudget() uint64
}

type Features interface {
	TimeLock() *TimeData
	Expiry() (*TimeData, iotago.Address) // return expiry time data and sender address or nil, nil if does not exist
	ReturnAmount() (uint64, bool)
}

type AsOffLedger interface {
	Nonce() uint64
}

type AsOnLedger interface {
	Output() iotago.Output
	IsInternalUTXO(*ChainID) bool
	UTXOInput() iotago.UTXOInput
	Features() Features
}

type ReturnAmountOptions interface {
	ReturnTo() iotago.Address
	Amount() uint64
}

func TakeRequestIDs(reqs ...Request) []RequestID {
	ret := make([]RequestID, len(reqs))
	for i := range reqs {
		ret[i] = reqs[i].ID()
	}
	return ret
}

// RequestsInTransaction parses the transaction and extracts those outputs which are interpreted as a request to a chain
func RequestsInTransaction(tx *iotago.Transaction) (map[ChainID][]Request, error) {
	txid, err := tx.ID()
	if err != nil {
		return nil, err
	}

	ret := make(map[ChainID][]Request)
	for i, out := range tx.Essence.Outputs {
		switch out.(type) {
		case *iotago.BasicOutput, *iotago.NFTOutput:
			// process it
		default:
			// only BasicOutputs and NFTs are interpreted right now, // TODO other outputs
			continue
		}
		// wrap output into the iscp.Request
		odata, err := OnLedgerFromUTXO(out, &iotago.UTXOInput{
			TransactionID:          *txid,
			TransactionOutputIndex: uint16(i),
		})
		if err != nil {
			return nil, err // TODO: maybe log the error and keep processing?
		}

		addr := odata.TargetAddress()
		if addr.Type() != iotago.AddressAlias {
			continue
		}

		chainID := ChainIDFromAliasID(addr.(*iotago.AliasAddress).AliasID())

		if odata.IsInternalUTXO(&chainID) {
			continue
		}

		ret[chainID] = append(ret[chainID], odata)
	}
	return ret, nil
}

// don't process any request which deadline will expire within 1 minute
const RequestConsideredExpiredWindow = time.Minute * 1

func RequestIsExpired(req AsOnLedger, currentTime TimeData) bool {
	expiry, _ := req.Features().Expiry()
	if expiry == nil {
		return false
	}
	if expiry.MilestoneIndex != 0 && currentTime.MilestoneIndex >= expiry.MilestoneIndex {
		return false
	}
	return !expiry.Time.IsZero() && currentTime.Time.After(expiry.Time.Add(-RequestConsideredExpiredWindow))
}

func RequestIsUnlockable(req AsOnLedger, chainAddress iotago.Address, currentTime TimeData) bool {
	if RequestIsExpired(req, currentTime) {
		return false
	}

	output, _ := req.Output().(iotago.TransIndepIdentOutput)

	return output.UnlockableBy(chainAddress, &iotago.ExternalUnlockParameters{
		ConfMsIndex: currentTime.MilestoneIndex,
		ConfUnix:    uint32(currentTime.Time.Unix()),
	})
}

func RequestHash(req Request) hashing.HashValue {
	return hashing.HashData(req.Bytes())
}
