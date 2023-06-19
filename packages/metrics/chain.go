package metrics

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/nodeclient"
	"github.com/iotaledger/wasp/packages/isc"
)

const (
	labelNameChain                                  = "chain"
	labelNamePipeName                               = "pipe_name"
	labelNameMessageType                            = "message_type"
	labelNameInMilestone                            = "in_milestone"
	labelNameInStateOutputMetrics                   = "in_state_output"
	labelNameInAliasOutputMetrics                   = "in_alias_output"
	labelNameInOutputMetrics                        = "in_output"
	labelNameInOnLedgerRequestMetrics               = "in_on_ledger_request"
	labelNameInTxInclusionStateMetrics              = "in_tx_inclusion_state"
	labelNameOutPublishStateTransactionMetrics      = "out_publish_state_transaction"
	labelNameOutPublishGovernanceTransactionMetrics = "out_publish_gov_transaction"
	labelNameOutPullLatestOutputMetrics             = "out_pull_latest_output"
	labelNameOutPullTxInclusionStateMetrics         = "out_pull_tx_inclusion_state"
	labelNameOutPullOutputByIDMetrics               = "out_pull_output_by_id"
	labelTxPublishResult                            = "result"
	labelNameWebapiRequestOperation                 = "api_req_type"
	labelNameWebapiRequestStatusCode                = "api_req_status_code"
	labelNameWebapiEvmRPCSuccess                    = "success"
)

type IChainMetrics interface {
	IChainPipeMetrics
	IChainBlockWALMetrics
	IChainConsensusMetrics
	IChainMempoolMetrics
	IChainMessageMetrics
	IChainStateMetrics
	IChainStateManagerMetrics
	IChainNodeConnMetrics
	IWebAPIMetrics
	IStateMetrics
}

var (
	_ IChainMetrics = &emptyChainMetrics{}
	_ IChainMetrics = &chainMetrics{}
)

type messageMetric[T any] struct {
	provider        *ChainMetricsProvider
	metricsLabels   prometheus.Labels
	messagesCount   atomic.Uint32
	lastMessage     T
	lastMessageTime time.Time
}

func newMessageMetric[T any](provider *ChainMetricsProvider, msgType string) *messageMetric[T] {
	metricsLabels := prometheus.Labels{
		labelNameMessageType: msgType,
	}

	// init values so they appear in prometheus
	provider.messagesL1.With(metricsLabels)
	provider.lastL1MessageTime.With(metricsLabels)

	return &messageMetric[T]{
		provider:      provider,
		metricsLabels: metricsLabels,
	}
}

func (m *messageMetric[T]) IncMessages(msg T, ts ...time.Time) {
	timestamp := time.Now()
	if len(ts) > 0 {
		timestamp = ts[0]
	}

	m.messagesCount.Add(1)
	m.lastMessage = msg
	m.lastMessageTime = timestamp

	m.provider.messagesL1.With(m.metricsLabels).Inc()
	m.provider.lastL1MessageTime.With(m.metricsLabels).Set(float64(timestamp.Unix()))
}

func (m *messageMetric[T]) MessagesTotal() uint32 {
	return m.messagesCount.Load()
}

func (m *messageMetric[T]) LastMessageTime() time.Time {
	return m.lastMessageTime
}

func (m *messageMetric[T]) LastMessage() T {
	return m.lastMessage
}

type emptyChainMetrics struct {
	IChainPipeMetrics
	IChainBlockWALMetrics
	IChainConsensusMetrics
	IChainMempoolMetrics
	IChainMessageMetrics
	IChainStateMetrics
	IChainStateManagerMetrics
	IChainNodeConnMetrics
	IWebAPIMetrics
	IStateMetrics
}

func NewEmptyChainMetrics() IChainMetrics {
	return &emptyChainMetrics{
		IChainPipeMetrics:         NewEmptyChainPipeMetrics(),
		IChainBlockWALMetrics:     NewEmptyChainBlockWALMetrics(),
		IChainConsensusMetrics:    NewEmptyChainConsensusMetric(),
		IChainMempoolMetrics:      NewEmptyChainMempoolMetric(),
		IChainMessageMetrics:      NewEmptyChainMessageMetrics(),
		IChainStateMetrics:        NewEmptyChainStateMetric(),
		IChainStateManagerMetrics: NewEmptyChainStateManagerMetric(),
		IChainNodeConnMetrics:     NewEmptyChainNodeConnMetric(),
		IWebAPIMetrics:            NewEmptyWebAPIMetrics(),
		IStateMetrics:             NewEmptyStateMetrics(),
	}
}

type chainMetrics struct {
	*chainPipeMetrics
	*chainBlockWALMetrics
	*chainConsensusMetric
	*chainMempoolMetric
	*chainMessageMetrics
	*chainStateMetric
	*chainStateManagerMetric
	*chainNodeConnMetric
	*webAPIChainMetrics
	*stateMetrics
}

func newChainMetrics(provider *ChainMetricsProvider, chainID isc.ChainID) *chainMetrics {
	return &chainMetrics{
		chainPipeMetrics:        newChainPipeMetric(provider, chainID),
		chainBlockWALMetrics:    newChainBlockWALMetrics(provider, chainID),
		chainConsensusMetric:    newChainConsensusMetric(provider, chainID),
		chainMempoolMetric:      newChainMempoolMetric(provider, chainID),
		chainMessageMetrics:     newChainMessageMetrics(provider, chainID),
		chainStateMetric:        newChainStateMetric(provider, chainID),
		chainStateManagerMetric: newChainStateManagerMetric(provider, chainID),
		chainNodeConnMetric:     newChainNodeConnMetric(provider, chainID),
		webAPIChainMetrics:      newWebAPIChainMetrics(provider, chainID),
		stateMetrics:            newStateMetrics(provider, chainID),
	}
}

// ChainMetricsProvider holds all metrics for all chains per chain
type ChainMetricsProvider struct {
	chainsLock       *sync.RWMutex
	chainsRegistered map[isc.ChainID]*chainMetrics

	// We use Func variant of a metric here, thus we register them
	// explicitly when they are created. Therefore we need a registry here.
	pipeLenRegistry *prometheus.Registry

	// blockWAL
	blockWALFailedWrites *prometheus.CounterVec
	blockWALFailedReads  *prometheus.CounterVec
	blockWALBlocksAdded  *countAndMaxMetrics

	// consensus
	consensusVMRunTime       *prometheus.HistogramVec
	consensusVMRunTimePerReq *prometheus.HistogramVec
	consensusVMRunReqCount   *prometheus.HistogramVec

	// mempool
	blocksTotalPerChain       *prometheus.CounterVec // TODO: Outdated and should be removed?
	requestsReceivedOffLedger *prometheus.CounterVec // TODO: Outdated and should be removed?
	requestsReceivedOnLedger  *prometheus.CounterVec // TODO: Outdated and should be removed?
	requestsProcessed         *prometheus.CounterVec // TODO: Outdated and should be removed?
	requestsAckMessages       *prometheus.CounterVec // TODO: Outdated and should be removed?
	requestsProcessingTime    *prometheus.GaugeVec   // TODO: Outdated and should be removed?

	mempoolTimePoolSize      *prometheus.GaugeVec
	mempoolOnLedgerPoolSize  *prometheus.GaugeVec
	mempoolOnLedgerReqTime   *prometheus.HistogramVec
	mempoolOffLedgerPoolSize *prometheus.GaugeVec
	mempoolOffLedgerReqTime  *prometheus.HistogramVec
	mempoolTotalSize         *prometheus.GaugeVec
	mempoolMissingReqs       *prometheus.GaugeVec

	// messages
	messagesL1             *prometheus.CounterVec // TODO: Outdated and should be removed?
	lastL1MessageTime      *prometheus.GaugeVec   // TODO: Outdated and should be removed?
	messagesL1Chain        *prometheus.CounterVec // TODO: Outdated and should be removed?
	lastL1MessageTimeChain *prometheus.GaugeVec   // TODO: Outdated and should be removed?

	inMilestoneMetrics                     *messageMetric[*nodeclient.MilestoneInfo] // TODO: Outdated and should be removed?
	inStateOutputMetrics                   *messageMetric[*InStateOutput]            // TODO: Outdated and should be removed?
	inAliasOutputMetrics                   *messageMetric[*iotago.AliasOutput]       // TODO: Outdated and should be removed?
	inOutputMetrics                        *messageMetric[*InOutput]                 // TODO: Outdated and should be removed?
	inOnLedgerRequestMetrics               *messageMetric[isc.OnLedgerRequest]       // TODO: Outdated and should be removed?
	inTxInclusionStateMetrics              *messageMetric[*TxInclusionStateMsg]      // TODO: Outdated and should be removed?
	outPublishStateTransactionMetrics      *messageMetric[*StateTransaction]         // TODO: Outdated and should be removed?
	outPublishGovernanceTransactionMetrics *messageMetric[*iotago.Transaction]       // TODO: Outdated and should be removed?
	outPullLatestOutputMetrics             *messageMetric[interface{}]               // TODO: Outdated and should be removed?
	outPullTxInclusionStateMetrics         *messageMetric[iotago.TransactionID]      // TODO: Outdated and should be removed?
	outPullOutputByIDMetrics               *messageMetric[iotago.OutputID]           // TODO: Outdated and should be removed?

	// chain state / tips
	chainActiveStateWant    *prometheus.GaugeVec
	chainActiveStateHave    *prometheus.GaugeVec
	chainConfirmedStateWant *prometheus.GaugeVec
	chainConfirmedStateHave *prometheus.GaugeVec
	chainConfirmedStateLag  ChainStateLag

	// state manager
	smCacheSize                  *prometheus.GaugeVec
	smBlocksFetching             *prometheus.GaugeVec
	smBlocksPending              *prometheus.GaugeVec
	smBlocksCommitted            *countAndMaxMetrics
	smRequestsWaiting            *prometheus.GaugeVec
	smCSPHandlingDuration        *prometheus.HistogramVec
	smCDSHandlingDuration        *prometheus.HistogramVec
	smCBPHandlingDuration        *prometheus.HistogramVec
	smFSDHandlingDuration        *prometheus.HistogramVec
	smTTHandlingDuration         *prometheus.HistogramVec
	smBlockFetchDuration         *prometheus.HistogramVec
	smPruningRunDuration         *prometheus.HistogramVec
	smPruningSingleStateDuration *prometheus.HistogramVec
	smPruningStatesInRun         *prometheus.HistogramVec
	smStatesPruned               *countAndMaxMetrics

	// node conn
	ncL1RequestReceived     *prometheus.CounterVec
	ncL1AliasOutputReceived *prometheus.CounterVec
	ncTXPublishStarted      *prometheus.CounterVec
	ncTXPublishResult       *prometheus.HistogramVec

	// webapi
	webAPIRequests    *prometheus.HistogramVec
	webAPIEvmRPCCalls *prometheus.HistogramVec

	// state
	stateBlockCommitTimes            *prometheus.HistogramVec
	stateBlockCommitNewTrieNodes     *prometheus.CounterVec
	stateBlockCommitNewTrieValues    *prometheus.CounterVec
	stateBlockPruneTimes             *prometheus.HistogramVec
	stateBlockPruneDeletedTrieNodes  *prometheus.CounterVec
	stateBlockPruneDeletedTrieValues *prometheus.CounterVec
}

var (
	postTimeBuckets = prometheus.ExponentialBucketsRange(0.1, 60*60, 17) // Time to confirm/reject a TX in L1 [0.1s - 1h].
	execTimeBuckets = prometheus.ExponentialBucketsRange(0.01, 100, 17)  // Execution of misc functions.
	recCountBuckets = prometheus.ExponentialBucketsRange(1, 1000, 16)
)

//nolint:funlen
func NewChainMetricsProvider() *ChainMetricsProvider {
	m := &ChainMetricsProvider{
		chainsLock:       &sync.RWMutex{},
		chainsRegistered: map[isc.ChainID]*chainMetrics{},

		//
		// blockWAL
		//
		blockWALFailedWrites: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "wal",
			Name:      "failed_writes_total",
			Help:      "Total number of writes to WAL that failed",
		}, []string{labelNameChain}),
		blockWALFailedReads: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "wal",
			Name:      "failed_reads_total",
			Help:      "Total number of reads failed while replaying WAL",
		}, []string{labelNameChain}),
		blockWALBlocksAdded: newCountAndMaxMetrics(
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "wal",
				Name:      "blocks_added",
				Help:      "Total number of blocks added into WAL",
			}, []string{labelNameChain}),
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "wal",
				Name:      "max_block_index",
				Help:      "Largest index of block added into WAL",
			}, []string{labelNameChain}),
		),

		//
		// consensus
		//
		consensusVMRunTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "consensus",
			Name:      "vm_run_time",
			Help:      "Time (s) it takes to run the VM per chain block.",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		consensusVMRunTimePerReq: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "consensus",
			Name:      "vm_run_time_per_req",
			Help:      "Time (s) it takes to run the VM per request.",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		consensusVMRunReqCount: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "consensus",
			Name:      "vm_run_req_count",
			Help:      "Number of requests processed per VM run.",
			Buckets:   recCountBuckets,
		}, []string{labelNameChain}),

		//
		// mempool
		//
		blocksTotalPerChain: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "blocks",
			Name:      "total",
			Help:      "Number of blocks per chain",
		}, []string{labelNameChain}),
		requestsReceivedOffLedger: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "requests",
			Name:      "off_ledger_total",
			Help:      "Number of off-ledger requests made to chain",
		}, []string{labelNameChain}),
		requestsReceivedOnLedger: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "requests",
			Name:      "on_ledger_total",
			Help:      "Number of on-ledger requests made to the chain",
		}, []string{labelNameChain}),
		requestsProcessed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "requests",
			Name:      "processed_total",
			Help:      "Number of requests processed per chain",
		}, []string{labelNameChain}),
		requestsAckMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "requests",
			Name:      "received_acks_total",
			Help:      "Number of received request acknowledgements per chain",
		}, []string{labelNameChain}),
		requestsProcessingTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "requests",
			Name:      "processing_time",
			Help:      "Time to process requests per chain",
		}, []string{labelNameChain}),
		mempoolTimePoolSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "time_pool_size",
			Help:      "Number of postponed requests in mempool.",
		}, []string{labelNameChain}),
		mempoolOnLedgerPoolSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "on_ledger_pool_size",
			Help:      "Number of On Ledger requests in mempool.",
		}, []string{labelNameChain}),
		mempoolOnLedgerReqTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "on_ledger_req_time",
			Help:      "Time (s) an on-ledger request stayed in the mempool before removing it.",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		mempoolOffLedgerPoolSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "off_ledger_pool_size",
			Help:      "Number of Off Ledger requests in mempool.",
		}, []string{labelNameChain}),
		mempoolOffLedgerReqTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "off_ledger_req_time",
			Help:      "Time (s) an off-ledger request stayed in the mempool before removing it.",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		mempoolTotalSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "total_pool_size",
			Help:      "Total requests in mempool.",
		}, []string{labelNameChain}),
		mempoolMissingReqs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "mempool",
			Name:      "missing_reqs",
			Help:      "Number of requests missing at this node (asking others to send them).",
		}, []string{labelNameChain}),

		//
		// messages // TODO: Review, if they are used/needed.
		//

		messagesL1: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "messages",
			Name:      "messages_total",
			Help:      "Number of messages sent/received by L1 connection",
		}, []string{labelNameMessageType}),
		lastL1MessageTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "messages",
			Name:      "last_message_time",
			Help:      "Last time when a message was sent/received by L1 connection",
		}, []string{labelNameMessageType}),
		messagesL1Chain: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "messages",
			Name:      "chain_messages_total",
			Help:      "Number of messages sent/received by L1 connection of the chain",
		}, []string{labelNameChain, labelNameMessageType}),
		lastL1MessageTimeChain: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "messages",
			Name:      "chain_last_message_time",
			Help:      "Last time when a message was sent/received by L1 connection of the chain",
		}, []string{labelNameChain, labelNameMessageType}),

		//
		// chain state / tips.
		//
		chainActiveStateWant: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "chain",
			Name:      "active_state_want",
			Help:      "We try to get blocks till this StateIndex for the active state.",
		}, []string{labelNameChain}),
		chainActiveStateHave: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "chain",
			Name:      "active_state_have",
			Help:      "We received blocks till this StateIndex for the active state.",
		}, []string{labelNameChain}),
		chainConfirmedStateWant: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "chain",
			Name:      "confirmed_state_want",
			Help:      "We try to get blocks till this StateIndex for the confirmed state.",
		}, []string{labelNameChain}),
		chainConfirmedStateHave: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "chain",
			Name:      "confirmed_state_have",
			Help:      "We received blocks till this StateIndex for the confirmed state.",
		}, []string{labelNameChain}),

		//
		// state manager
		//
		smCacheSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "cache_size",
			Help:      "Number of blocks stored in cache",
		}, []string{labelNameChain}),
		smBlocksFetching: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "blocks_fetching",
			Help:      "Number of blocks the node is waiting from other nodes",
		}, []string{labelNameChain}),
		smBlocksPending: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "blocks_pending",
			Help:      "Number of blocks the node has fetched but hasn't committed, because the node doesn't have their ancestors",
		}, []string{labelNameChain}),
		smBlocksCommitted: newCountAndMaxMetrics(
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "state_manager",
				Name:      "blocks_committed",
				Help:      "Number of blocks the node has committed to the store",
			}, []string{labelNameChain}),
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "state_manager",
				Name:      "max_blocks_index_committed",
				Help:      "Largest index of block committed to the store",
			}, []string{labelNameChain}),
		),
		smRequestsWaiting: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "requests_waiting",
			Help:      "Number of requests from other components of the node waiting for response from the state manager. Note that StateDiff request is counted as two requests as it has to obtain two possibly different blocks.",
		}, []string{labelNameChain}),
		smCSPHandlingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "consensus_state_proposal_duration",
			Help:      "The duration (s) from starting handling ConsensusStateProposal request till responding to the consensus",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smCDSHandlingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "consensus_decided_state_duration",
			Help:      "The duration (s) from starting handling ConsensusDecidedState request till responding to the consensus",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smCBPHandlingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "consensus_block_produced_duration",
			Help:      "The duration (s) from starting till finishing handling ConsensusBlockProduced, which includes responding to the consensus",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smFSDHandlingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "chain_fetch_state_diff_duration",
			Help:      "The duration (s) from starting handling ChainFetchStateDiff request till responding to the chain",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smTTHandlingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "timer_tick_duration",
			Help:      "The duration (s) from starting till finishing handling StateManagerTimerTick request",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smBlockFetchDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "block_fetch_duration",
			Help:      "The duration (s) from starting fetching block from other till it is received in this node",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smPruningRunDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "pruning_run_duration",
			Help:      "The duration (s) from starting till finishing pruning run, which may include pruning several states from store",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smPruningSingleStateDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "pruning_single_state_duration",
			Help:      "The duration (s) from starting till finishing pruning single state from store",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain}),
		smPruningStatesInRun: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "state_manager",
			Name:      "pruning_states_in_run",
			Help:      "Number of states pruned in single pruning run (should be 1 in normally running nodes)",
			Buckets:   recCountBuckets,
		}, []string{labelNameChain}),
		smStatesPruned: newCountAndMaxMetrics(
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "state_manager",
				Name:      "sates_pruned",
				Help:      "Number of states pruned in total since starting the node",
			}, []string{labelNameChain}),
			prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "iota_wasp",
				Subsystem: "state_manager",
				Name:      "max_state_index_pruned",
				Help:      "Largest index of state pruned from the store",
			}, []string{labelNameChain}),
		),

		//
		// node conn
		//
		ncL1RequestReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "node_conn",
			Name:      "l1_request_received",
			Help:      "A number of confirmed requests received from L1.",
		}, []string{labelNameChain}),
		ncL1AliasOutputReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "node_conn",
			Name:      "l1_alias_output_received",
			Help:      "A number of confirmed alias outputs received from L1.",
		}, []string{labelNameChain}),
		ncTXPublishStarted: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "iota_wasp",
			Subsystem: "node_conn",
			Name:      "tx_publish_started",
			Help:      "A number of transactions submitted for publication in L1.",
		}, []string{labelNameChain}),
		ncTXPublishResult: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "node_conn",
			Name:      "tx_publish_result",
			Help:      "The duration (s) to publish a transaction.",
			Buckets:   postTimeBuckets,
		}, []string{labelNameChain, labelTxPublishResult}),

		//
		// webapi
		//
		webAPIRequests: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "webapi",
			Name:      "webapi_requests",
			Help:      "Time elapsed (s) processing requests",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain, labelNameWebapiRequestOperation, labelNameWebapiRequestStatusCode}),
		webAPIEvmRPCCalls: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "iota_wasp",
			Subsystem: "webapi",
			Name:      "webapi_evm_rpc_calls",
			Help:      "Time elapsed (s) processing evm rpc requests",
			Buckets:   execTimeBuckets,
		}, []string{labelNameChain, labelNameWebapiRequestOperation, labelNameWebapiEvmRPCSuccess}),
	}

	m.inMilestoneMetrics = newMessageMetric[*nodeclient.MilestoneInfo](m, labelNameInMilestone)
	m.inStateOutputMetrics = newMessageMetric[*InStateOutput](m, labelNameInStateOutputMetrics)
	m.inAliasOutputMetrics = newMessageMetric[*iotago.AliasOutput](m, labelNameInAliasOutputMetrics)
	m.inOutputMetrics = newMessageMetric[*InOutput](m, labelNameInOutputMetrics)
	m.inOnLedgerRequestMetrics = newMessageMetric[isc.OnLedgerRequest](m, labelNameInOnLedgerRequestMetrics)
	m.inTxInclusionStateMetrics = newMessageMetric[*TxInclusionStateMsg](m, labelNameInTxInclusionStateMetrics)
	m.outPublishStateTransactionMetrics = newMessageMetric[*StateTransaction](m, labelNameOutPublishStateTransactionMetrics)
	m.outPublishGovernanceTransactionMetrics = newMessageMetric[*iotago.Transaction](m, labelNameOutPublishGovernanceTransactionMetrics)
	m.outPullLatestOutputMetrics = newMessageMetric[interface{}](m, labelNameOutPullLatestOutputMetrics)
	m.outPullTxInclusionStateMetrics = newMessageMetric[iotago.TransactionID](m, labelNameOutPullTxInclusionStateMetrics)
	m.outPullOutputByIDMetrics = newMessageMetric[iotago.OutputID](m, labelNameOutPullOutputByIDMetrics)

	m.chainConfirmedStateLag = make(ChainStateLag)

	initStateMetrics(m)

	return m
}

func (m *ChainMetricsProvider) GetChainMetrics(chainID isc.ChainID) IChainMetrics {
	m.chainsLock.Lock()
	defer m.chainsLock.Unlock()

	if cm, ok := m.chainsRegistered[chainID]; ok {
		return cm
	}
	cm := newChainMetrics(m, chainID)
	m.chainsRegistered[chainID] = cm
	return cm
}

func (m *ChainMetricsProvider) PrometheusCollectorsBlockWAL() []prometheus.Collector {
	return append(
		[]prometheus.Collector{
			m.blockWALFailedWrites,
			m.blockWALFailedReads,
		},
		m.blockWALBlocksAdded.collectors()...,
	)
}

func (m *ChainMetricsProvider) PrometheusCollectorsConsensus() []prometheus.Collector {
	return []prometheus.Collector{
		m.consensusVMRunTime,
		m.consensusVMRunTimePerReq,
		m.consensusVMRunReqCount,
	}
}

func (m *ChainMetricsProvider) PrometheusCollectorsMempool() []prometheus.Collector {
	return []prometheus.Collector{
		m.blocksTotalPerChain,
		m.requestsReceivedOffLedger,
		m.requestsReceivedOnLedger,
		m.requestsProcessed,
		m.requestsAckMessages,
		m.requestsProcessingTime,
		m.mempoolTimePoolSize,
		m.mempoolOnLedgerPoolSize,
		m.mempoolOnLedgerReqTime,
		m.mempoolOffLedgerPoolSize,
		m.mempoolOffLedgerReqTime,
		m.mempoolTotalSize,
		m.mempoolMissingReqs,
	}
}

func (m *ChainMetricsProvider) PrometheusCollectorsChainMessages() []prometheus.Collector {
	return []prometheus.Collector{
		m.messagesL1,
		m.lastL1MessageTime,
		m.messagesL1Chain,
		m.lastL1MessageTimeChain,
	}
}

func (m *ChainMetricsProvider) PrometheusCollectorsChainState() []prometheus.Collector {
	return []prometheus.Collector{
		m.chainActiveStateWant,
		m.chainActiveStateHave,
		m.chainConfirmedStateWant,
		m.chainConfirmedStateHave,
	}
}

func (m *ChainMetricsProvider) PrometheusCollectorsChainStateManager() []prometheus.Collector {
	return append(
		append(
			[]prometheus.Collector{
				m.smCacheSize,
				m.smBlocksFetching,
				m.smBlocksPending,
				m.smRequestsWaiting,
				m.smCSPHandlingDuration,
				m.smCDSHandlingDuration,
				m.smCBPHandlingDuration,
				m.smFSDHandlingDuration,
				m.smTTHandlingDuration,
				m.smBlockFetchDuration,
				m.smPruningRunDuration,
				m.smPruningSingleStateDuration,
				m.smPruningStatesInRun,
			},
			m.smBlocksCommitted.collectors()...,
		),
		m.smStatesPruned.collectors()...,
	)
}

func (m *ChainMetricsProvider) PrometheusCollectorsChainNodeConn() []prometheus.Collector {
	return []prometheus.Collector{
		m.ncL1RequestReceived,
		m.ncL1AliasOutputReceived,
		m.ncTXPublishStarted,
		m.ncTXPublishResult,
	}
}

func (m *ChainMetricsProvider) PrometheusRegisterChainPipeMetrics(reg *prometheus.Registry) {
	m.pipeLenRegistry = reg
}

func (m *ChainMetricsProvider) PrometheusCollectorsWebAPI() []prometheus.Collector {
	return []prometheus.Collector{
		m.webAPIRequests,
		m.webAPIEvmRPCCalls,
	}
}

func (m *ChainMetricsProvider) PrometheusCollectorsState() []prometheus.Collector {
	return []prometheus.Collector{
		m.stateBlockCommitTimes,
		m.stateBlockCommitNewTrieNodes,
		m.stateBlockCommitNewTrieValues,
		m.stateBlockPruneTimes,
		m.stateBlockPruneDeletedTrieNodes,
		m.stateBlockPruneDeletedTrieValues,
	}
}

func (m *ChainMetricsProvider) RegisterChain(chainID isc.ChainID) {
	m.GetChainMetrics(chainID)
}

func (m *ChainMetricsProvider) UnregisterChain(chainID isc.ChainID) {
	m.chainsLock.Lock()
	defer m.chainsLock.Unlock()

	if cm, ok := m.chainsRegistered[chainID]; ok {
		cm.cleanup()
		delete(m.chainsRegistered, chainID)
	}
}

func (m *ChainMetricsProvider) RegisteredChains() []isc.ChainID {
	m.chainsLock.RLock()
	defer m.chainsLock.RUnlock()

	return lo.Keys(m.chainsRegistered)
}

func (m *ChainMetricsProvider) InMilestone() IMessageMetric[*nodeclient.MilestoneInfo] {
	return m.inMilestoneMetrics
}

func (m *ChainMetricsProvider) InStateOutput() IMessageMetric[*InStateOutput] {
	return m.inStateOutputMetrics
}

func (m *ChainMetricsProvider) InAliasOutput() IMessageMetric[*iotago.AliasOutput] {
	return m.inAliasOutputMetrics
}

func (m *ChainMetricsProvider) InOutput() IMessageMetric[*InOutput] {
	return m.inOutputMetrics
}

func (m *ChainMetricsProvider) InOnLedgerRequest() IMessageMetric[isc.OnLedgerRequest] {
	return m.inOnLedgerRequestMetrics
}

func (m *ChainMetricsProvider) InTxInclusionState() IMessageMetric[*TxInclusionStateMsg] {
	return m.inTxInclusionStateMetrics
}

func (m *ChainMetricsProvider) OutPublishStateTransaction() IMessageMetric[*StateTransaction] {
	return m.outPublishStateTransactionMetrics
}

func (m *ChainMetricsProvider) OutPublishGovernanceTransaction() IMessageMetric[*iotago.Transaction] {
	return m.outPublishGovernanceTransactionMetrics
}

func (m *ChainMetricsProvider) OutPullLatestOutput() IMessageMetric[interface{}] {
	return m.outPullLatestOutputMetrics
}

func (m *ChainMetricsProvider) OutPullTxInclusionState() IMessageMetric[iotago.TransactionID] {
	return m.outPullTxInclusionStateMetrics
}

func (m *ChainMetricsProvider) OutPullOutputByID() IMessageMetric[iotago.OutputID] {
	return m.outPullOutputByIDMetrics
}

func (m *ChainMetricsProvider) MaxChainConfirmedStateLag() uint32 {
	return m.chainConfirmedStateLag.MaxLag()
}
