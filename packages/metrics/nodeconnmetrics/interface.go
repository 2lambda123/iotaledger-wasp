package nodeconnmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/nodeclient"
	"github.com/iotaledger/wasp/packages/isc"
)

type StateTransaction struct {
	StateIndex  uint32
	Transaction *iotago.Transaction
}

type InStateOutput struct {
	OutputID iotago.OutputID
	Output   iotago.Output
}

type InOutput struct {
	OutputID iotago.OutputID
	Output   iotago.Output
}

type TxInclusionStateMsg struct {
	TxID  iotago.TransactionID
	State string
}

type NodeConnectionMessageMetrics[T any] interface {
	IncL1Messages(T)
	GetL1MessagesTotal() uint32
	GetLastL1MessageTime() time.Time
	GetLastL1Message() T
}

type NodeConnectionMessagesMetrics interface {
	GetOutPublishStateTransaction() NodeConnectionMessageMetrics[*StateTransaction]
	GetOutPublishGovernanceTransaction() NodeConnectionMessageMetrics[*iotago.Transaction]
	GetOutPullLatestOutput() NodeConnectionMessageMetrics[interface{}]
	GetOutPullTxInclusionState() NodeConnectionMessageMetrics[iotago.TransactionID]
	GetOutPullOutputByID() NodeConnectionMessageMetrics[iotago.OutputID]
	GetInStateOutput() NodeConnectionMessageMetrics[*InStateOutput]
	GetInAliasOutput() NodeConnectionMessageMetrics[*iotago.AliasOutput]
	GetInOutput() NodeConnectionMessageMetrics[*InOutput]
	GetInOnLedgerRequest() NodeConnectionMessageMetrics[isc.OnLedgerRequest]
	GetInTxInclusionState() NodeConnectionMessageMetrics[*TxInclusionStateMsg]
}

type NodeConnectionMetrics interface {
	NodeConnectionMessagesMetrics
	GetInMilestone() NodeConnectionMessageMetrics[*nodeclient.MilestoneInfo]
	SetRegistered(isc.ChainID)
	SetUnregistered(isc.ChainID)
	GetRegistered() []isc.ChainID
	PrometheusCollectors() []prometheus.Collector
	NewMessagesMetrics(isc.ChainID) NodeConnectionMessagesMetrics
}
