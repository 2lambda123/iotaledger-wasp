/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the NodeMessageMetrics type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NodeMessageMetrics{}

// NodeMessageMetrics struct for NodeMessageMetrics
type NodeMessageMetrics struct {
	InAliasOutput                   AliasOutputMetricItem         `json:"inAliasOutput"`
	InMilestone                     MilestoneMetricItem           `json:"inMilestone"`
	InOnLedgerRequest               OnLedgerRequestMetricItem     `json:"inOnLedgerRequest"`
	InOutput                        InOutputMetricItem            `json:"inOutput"`
	InStateOutput                   InStateOutputMetricItem       `json:"inStateOutput"`
	InTxInclusionState              TxInclusionStateMsgMetricItem `json:"inTxInclusionState"`
	OutPublishGovernanceTransaction TransactionMetricItem         `json:"outPublishGovernanceTransaction"`
	OutPublisherStateTransaction    PublisherStateTransactionItem `json:"outPublisherStateTransaction"`
	OutPullLatestOutput             InterfaceMetricItem           `json:"outPullLatestOutput"`
	OutPullOutputByID               UTXOInputMetricItem           `json:"outPullOutputByID"`
	OutPullTxInclusionState         TransactionIDMetricItem       `json:"outPullTxInclusionState"`
	RegisteredChainIDs              []string                      `json:"registeredChainIDs"`
}

// NewNodeMessageMetrics instantiates a new NodeMessageMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodeMessageMetrics(inAliasOutput AliasOutputMetricItem, inMilestone MilestoneMetricItem, inOnLedgerRequest OnLedgerRequestMetricItem, inOutput InOutputMetricItem, inStateOutput InStateOutputMetricItem, inTxInclusionState TxInclusionStateMsgMetricItem, outPublishGovernanceTransaction TransactionMetricItem, outPublisherStateTransaction PublisherStateTransactionItem, outPullLatestOutput InterfaceMetricItem, outPullOutputByID UTXOInputMetricItem, outPullTxInclusionState TransactionIDMetricItem, registeredChainIDs []string) *NodeMessageMetrics {
	this := NodeMessageMetrics{}
	this.InAliasOutput = inAliasOutput
	this.InMilestone = inMilestone
	this.InOnLedgerRequest = inOnLedgerRequest
	this.InOutput = inOutput
	this.InStateOutput = inStateOutput
	this.InTxInclusionState = inTxInclusionState
	this.OutPublishGovernanceTransaction = outPublishGovernanceTransaction
	this.OutPublisherStateTransaction = outPublisherStateTransaction
	this.OutPullLatestOutput = outPullLatestOutput
	this.OutPullOutputByID = outPullOutputByID
	this.OutPullTxInclusionState = outPullTxInclusionState
	this.RegisteredChainIDs = registeredChainIDs
	return &this
}

// NewNodeMessageMetricsWithDefaults instantiates a new NodeMessageMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodeMessageMetricsWithDefaults() *NodeMessageMetrics {
	this := NodeMessageMetrics{}
	return &this
}

// GetInAliasOutput returns the InAliasOutput field value
func (o *NodeMessageMetrics) GetInAliasOutput() AliasOutputMetricItem {
	if o == nil {
		var ret AliasOutputMetricItem
		return ret
	}

	return o.InAliasOutput
}

// GetInAliasOutputOk returns a tuple with the InAliasOutput field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInAliasOutputOk() (*AliasOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InAliasOutput, true
}

// SetInAliasOutput sets field value
func (o *NodeMessageMetrics) SetInAliasOutput(v AliasOutputMetricItem) {
	o.InAliasOutput = v
}

// GetInMilestone returns the InMilestone field value
func (o *NodeMessageMetrics) GetInMilestone() MilestoneMetricItem {
	if o == nil {
		var ret MilestoneMetricItem
		return ret
	}

	return o.InMilestone
}

// GetInMilestoneOk returns a tuple with the InMilestone field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInMilestoneOk() (*MilestoneMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InMilestone, true
}

// SetInMilestone sets field value
func (o *NodeMessageMetrics) SetInMilestone(v MilestoneMetricItem) {
	o.InMilestone = v
}

// GetInOnLedgerRequest returns the InOnLedgerRequest field value
func (o *NodeMessageMetrics) GetInOnLedgerRequest() OnLedgerRequestMetricItem {
	if o == nil {
		var ret OnLedgerRequestMetricItem
		return ret
	}

	return o.InOnLedgerRequest
}

// GetInOnLedgerRequestOk returns a tuple with the InOnLedgerRequest field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInOnLedgerRequestOk() (*OnLedgerRequestMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InOnLedgerRequest, true
}

// SetInOnLedgerRequest sets field value
func (o *NodeMessageMetrics) SetInOnLedgerRequest(v OnLedgerRequestMetricItem) {
	o.InOnLedgerRequest = v
}

// GetInOutput returns the InOutput field value
func (o *NodeMessageMetrics) GetInOutput() InOutputMetricItem {
	if o == nil {
		var ret InOutputMetricItem
		return ret
	}

	return o.InOutput
}

// GetInOutputOk returns a tuple with the InOutput field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInOutputOk() (*InOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InOutput, true
}

// SetInOutput sets field value
func (o *NodeMessageMetrics) SetInOutput(v InOutputMetricItem) {
	o.InOutput = v
}

// GetInStateOutput returns the InStateOutput field value
func (o *NodeMessageMetrics) GetInStateOutput() InStateOutputMetricItem {
	if o == nil {
		var ret InStateOutputMetricItem
		return ret
	}

	return o.InStateOutput
}

// GetInStateOutputOk returns a tuple with the InStateOutput field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInStateOutputOk() (*InStateOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InStateOutput, true
}

// SetInStateOutput sets field value
func (o *NodeMessageMetrics) SetInStateOutput(v InStateOutputMetricItem) {
	o.InStateOutput = v
}

// GetInTxInclusionState returns the InTxInclusionState field value
func (o *NodeMessageMetrics) GetInTxInclusionState() TxInclusionStateMsgMetricItem {
	if o == nil {
		var ret TxInclusionStateMsgMetricItem
		return ret
	}

	return o.InTxInclusionState
}

// GetInTxInclusionStateOk returns a tuple with the InTxInclusionState field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetInTxInclusionStateOk() (*TxInclusionStateMsgMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InTxInclusionState, true
}

// SetInTxInclusionState sets field value
func (o *NodeMessageMetrics) SetInTxInclusionState(v TxInclusionStateMsgMetricItem) {
	o.InTxInclusionState = v
}

// GetOutPublishGovernanceTransaction returns the OutPublishGovernanceTransaction field value
func (o *NodeMessageMetrics) GetOutPublishGovernanceTransaction() TransactionMetricItem {
	if o == nil {
		var ret TransactionMetricItem
		return ret
	}

	return o.OutPublishGovernanceTransaction
}

// GetOutPublishGovernanceTransactionOk returns a tuple with the OutPublishGovernanceTransaction field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetOutPublishGovernanceTransactionOk() (*TransactionMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPublishGovernanceTransaction, true
}

// SetOutPublishGovernanceTransaction sets field value
func (o *NodeMessageMetrics) SetOutPublishGovernanceTransaction(v TransactionMetricItem) {
	o.OutPublishGovernanceTransaction = v
}

// GetOutPublisherStateTransaction returns the OutPublisherStateTransaction field value
func (o *NodeMessageMetrics) GetOutPublisherStateTransaction() PublisherStateTransactionItem {
	if o == nil {
		var ret PublisherStateTransactionItem
		return ret
	}

	return o.OutPublisherStateTransaction
}

// GetOutPublisherStateTransactionOk returns a tuple with the OutPublisherStateTransaction field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetOutPublisherStateTransactionOk() (*PublisherStateTransactionItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPublisherStateTransaction, true
}

// SetOutPublisherStateTransaction sets field value
func (o *NodeMessageMetrics) SetOutPublisherStateTransaction(v PublisherStateTransactionItem) {
	o.OutPublisherStateTransaction = v
}

// GetOutPullLatestOutput returns the OutPullLatestOutput field value
func (o *NodeMessageMetrics) GetOutPullLatestOutput() InterfaceMetricItem {
	if o == nil {
		var ret InterfaceMetricItem
		return ret
	}

	return o.OutPullLatestOutput
}

// GetOutPullLatestOutputOk returns a tuple with the OutPullLatestOutput field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetOutPullLatestOutputOk() (*InterfaceMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullLatestOutput, true
}

// SetOutPullLatestOutput sets field value
func (o *NodeMessageMetrics) SetOutPullLatestOutput(v InterfaceMetricItem) {
	o.OutPullLatestOutput = v
}

// GetOutPullOutputByID returns the OutPullOutputByID field value
func (o *NodeMessageMetrics) GetOutPullOutputByID() UTXOInputMetricItem {
	if o == nil {
		var ret UTXOInputMetricItem
		return ret
	}

	return o.OutPullOutputByID
}

// GetOutPullOutputByIDOk returns a tuple with the OutPullOutputByID field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetOutPullOutputByIDOk() (*UTXOInputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullOutputByID, true
}

// SetOutPullOutputByID sets field value
func (o *NodeMessageMetrics) SetOutPullOutputByID(v UTXOInputMetricItem) {
	o.OutPullOutputByID = v
}

// GetOutPullTxInclusionState returns the OutPullTxInclusionState field value
func (o *NodeMessageMetrics) GetOutPullTxInclusionState() TransactionIDMetricItem {
	if o == nil {
		var ret TransactionIDMetricItem
		return ret
	}

	return o.OutPullTxInclusionState
}

// GetOutPullTxInclusionStateOk returns a tuple with the OutPullTxInclusionState field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetOutPullTxInclusionStateOk() (*TransactionIDMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullTxInclusionState, true
}

// SetOutPullTxInclusionState sets field value
func (o *NodeMessageMetrics) SetOutPullTxInclusionState(v TransactionIDMetricItem) {
	o.OutPullTxInclusionState = v
}

// GetRegisteredChainIDs returns the RegisteredChainIDs field value
func (o *NodeMessageMetrics) GetRegisteredChainIDs() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.RegisteredChainIDs
}

// GetRegisteredChainIDsOk returns a tuple with the RegisteredChainIDs field value
// and a boolean to check if the value has been set.
func (o *NodeMessageMetrics) GetRegisteredChainIDsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.RegisteredChainIDs, true
}

// SetRegisteredChainIDs sets field value
func (o *NodeMessageMetrics) SetRegisteredChainIDs(v []string) {
	o.RegisteredChainIDs = v
}

func (o NodeMessageMetrics) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NodeMessageMetrics) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["inAliasOutput"] = o.InAliasOutput
	toSerialize["inMilestone"] = o.InMilestone
	toSerialize["inOnLedgerRequest"] = o.InOnLedgerRequest
	toSerialize["inOutput"] = o.InOutput
	toSerialize["inStateOutput"] = o.InStateOutput
	toSerialize["inTxInclusionState"] = o.InTxInclusionState
	toSerialize["outPublishGovernanceTransaction"] = o.OutPublishGovernanceTransaction
	toSerialize["outPublisherStateTransaction"] = o.OutPublisherStateTransaction
	toSerialize["outPullLatestOutput"] = o.OutPullLatestOutput
	toSerialize["outPullOutputByID"] = o.OutPullOutputByID
	toSerialize["outPullTxInclusionState"] = o.OutPullTxInclusionState
	toSerialize["registeredChainIDs"] = o.RegisteredChainIDs
	return toSerialize, nil
}

type NullableNodeMessageMetrics struct {
	value *NodeMessageMetrics
	isSet bool
}

func (v NullableNodeMessageMetrics) Get() *NodeMessageMetrics {
	return v.value
}

func (v *NullableNodeMessageMetrics) Set(val *NodeMessageMetrics) {
	v.value = val
	v.isSet = true
}

func (v NullableNodeMessageMetrics) IsSet() bool {
	return v.isSet
}

func (v *NullableNodeMessageMetrics) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodeMessageMetrics(val *NodeMessageMetrics) *NullableNodeMessageMetrics {
	return &NullableNodeMessageMetrics{value: val, isSet: true}
}

func (v NullableNodeMessageMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodeMessageMetrics) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
