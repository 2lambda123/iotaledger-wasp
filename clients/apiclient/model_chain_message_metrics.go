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

// checks if the ChainMessageMetrics type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ChainMessageMetrics{}

// ChainMessageMetrics struct for ChainMessageMetrics
type ChainMessageMetrics struct {
	InAccountOutput AccountOutputMetricItem `json:"inAccountOutput"`
	InOnLedgerRequest OnLedgerRequestMetricItem `json:"inOnLedgerRequest"`
	InOutput InOutputMetricItem `json:"inOutput"`
	InStateOutput InStateOutputMetricItem `json:"inStateOutput"`
	InTxInclusionState TxInclusionStateMsgMetricItem `json:"inTxInclusionState"`
	OutPublishGovernanceTransaction TransactionMetricItem `json:"outPublishGovernanceTransaction"`
	OutPublisherStateTransaction PublisherStateTransactionItem `json:"outPublisherStateTransaction"`
	OutPullLatestOutput InterfaceMetricItem `json:"outPullLatestOutput"`
	OutPullOutputByID UTXOInputMetricItem `json:"outPullOutputByID"`
	OutPullTxInclusionState TransactionIDMetricItem `json:"outPullTxInclusionState"`
}

// NewChainMessageMetrics instantiates a new ChainMessageMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewChainMessageMetrics(inAccountOutput AccountOutputMetricItem, inOnLedgerRequest OnLedgerRequestMetricItem, inOutput InOutputMetricItem, inStateOutput InStateOutputMetricItem, inTxInclusionState TxInclusionStateMsgMetricItem, outPublishGovernanceTransaction TransactionMetricItem, outPublisherStateTransaction PublisherStateTransactionItem, outPullLatestOutput InterfaceMetricItem, outPullOutputByID UTXOInputMetricItem, outPullTxInclusionState TransactionIDMetricItem) *ChainMessageMetrics {
	this := ChainMessageMetrics{}
	this.InAccountOutput = inAccountOutput
	this.InOnLedgerRequest = inOnLedgerRequest
	this.InOutput = inOutput
	this.InStateOutput = inStateOutput
	this.InTxInclusionState = inTxInclusionState
	this.OutPublishGovernanceTransaction = outPublishGovernanceTransaction
	this.OutPublisherStateTransaction = outPublisherStateTransaction
	this.OutPullLatestOutput = outPullLatestOutput
	this.OutPullOutputByID = outPullOutputByID
	this.OutPullTxInclusionState = outPullTxInclusionState
	return &this
}

// NewChainMessageMetricsWithDefaults instantiates a new ChainMessageMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewChainMessageMetricsWithDefaults() *ChainMessageMetrics {
	this := ChainMessageMetrics{}
	return &this
}

// GetInAccountOutput returns the InAccountOutput field value
func (o *ChainMessageMetrics) GetInAccountOutput() AccountOutputMetricItem {
	if o == nil {
		var ret AccountOutputMetricItem
		return ret
	}

	return o.InAccountOutput
}

// GetInAccountOutputOk returns a tuple with the InAccountOutput field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetInAccountOutputOk() (*AccountOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InAccountOutput, true
}

// SetInAccountOutput sets field value
func (o *ChainMessageMetrics) SetInAccountOutput(v AccountOutputMetricItem) {
	o.InAccountOutput = v
}

// GetInOnLedgerRequest returns the InOnLedgerRequest field value
func (o *ChainMessageMetrics) GetInOnLedgerRequest() OnLedgerRequestMetricItem {
	if o == nil {
		var ret OnLedgerRequestMetricItem
		return ret
	}

	return o.InOnLedgerRequest
}

// GetInOnLedgerRequestOk returns a tuple with the InOnLedgerRequest field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetInOnLedgerRequestOk() (*OnLedgerRequestMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InOnLedgerRequest, true
}

// SetInOnLedgerRequest sets field value
func (o *ChainMessageMetrics) SetInOnLedgerRequest(v OnLedgerRequestMetricItem) {
	o.InOnLedgerRequest = v
}

// GetInOutput returns the InOutput field value
func (o *ChainMessageMetrics) GetInOutput() InOutputMetricItem {
	if o == nil {
		var ret InOutputMetricItem
		return ret
	}

	return o.InOutput
}

// GetInOutputOk returns a tuple with the InOutput field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetInOutputOk() (*InOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InOutput, true
}

// SetInOutput sets field value
func (o *ChainMessageMetrics) SetInOutput(v InOutputMetricItem) {
	o.InOutput = v
}

// GetInStateOutput returns the InStateOutput field value
func (o *ChainMessageMetrics) GetInStateOutput() InStateOutputMetricItem {
	if o == nil {
		var ret InStateOutputMetricItem
		return ret
	}

	return o.InStateOutput
}

// GetInStateOutputOk returns a tuple with the InStateOutput field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetInStateOutputOk() (*InStateOutputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InStateOutput, true
}

// SetInStateOutput sets field value
func (o *ChainMessageMetrics) SetInStateOutput(v InStateOutputMetricItem) {
	o.InStateOutput = v
}

// GetInTxInclusionState returns the InTxInclusionState field value
func (o *ChainMessageMetrics) GetInTxInclusionState() TxInclusionStateMsgMetricItem {
	if o == nil {
		var ret TxInclusionStateMsgMetricItem
		return ret
	}

	return o.InTxInclusionState
}

// GetInTxInclusionStateOk returns a tuple with the InTxInclusionState field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetInTxInclusionStateOk() (*TxInclusionStateMsgMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InTxInclusionState, true
}

// SetInTxInclusionState sets field value
func (o *ChainMessageMetrics) SetInTxInclusionState(v TxInclusionStateMsgMetricItem) {
	o.InTxInclusionState = v
}

// GetOutPublishGovernanceTransaction returns the OutPublishGovernanceTransaction field value
func (o *ChainMessageMetrics) GetOutPublishGovernanceTransaction() TransactionMetricItem {
	if o == nil {
		var ret TransactionMetricItem
		return ret
	}

	return o.OutPublishGovernanceTransaction
}

// GetOutPublishGovernanceTransactionOk returns a tuple with the OutPublishGovernanceTransaction field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetOutPublishGovernanceTransactionOk() (*TransactionMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPublishGovernanceTransaction, true
}

// SetOutPublishGovernanceTransaction sets field value
func (o *ChainMessageMetrics) SetOutPublishGovernanceTransaction(v TransactionMetricItem) {
	o.OutPublishGovernanceTransaction = v
}

// GetOutPublisherStateTransaction returns the OutPublisherStateTransaction field value
func (o *ChainMessageMetrics) GetOutPublisherStateTransaction() PublisherStateTransactionItem {
	if o == nil {
		var ret PublisherStateTransactionItem
		return ret
	}

	return o.OutPublisherStateTransaction
}

// GetOutPublisherStateTransactionOk returns a tuple with the OutPublisherStateTransaction field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetOutPublisherStateTransactionOk() (*PublisherStateTransactionItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPublisherStateTransaction, true
}

// SetOutPublisherStateTransaction sets field value
func (o *ChainMessageMetrics) SetOutPublisherStateTransaction(v PublisherStateTransactionItem) {
	o.OutPublisherStateTransaction = v
}

// GetOutPullLatestOutput returns the OutPullLatestOutput field value
func (o *ChainMessageMetrics) GetOutPullLatestOutput() InterfaceMetricItem {
	if o == nil {
		var ret InterfaceMetricItem
		return ret
	}

	return o.OutPullLatestOutput
}

// GetOutPullLatestOutputOk returns a tuple with the OutPullLatestOutput field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetOutPullLatestOutputOk() (*InterfaceMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullLatestOutput, true
}

// SetOutPullLatestOutput sets field value
func (o *ChainMessageMetrics) SetOutPullLatestOutput(v InterfaceMetricItem) {
	o.OutPullLatestOutput = v
}

// GetOutPullOutputByID returns the OutPullOutputByID field value
func (o *ChainMessageMetrics) GetOutPullOutputByID() UTXOInputMetricItem {
	if o == nil {
		var ret UTXOInputMetricItem
		return ret
	}

	return o.OutPullOutputByID
}

// GetOutPullOutputByIDOk returns a tuple with the OutPullOutputByID field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetOutPullOutputByIDOk() (*UTXOInputMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullOutputByID, true
}

// SetOutPullOutputByID sets field value
func (o *ChainMessageMetrics) SetOutPullOutputByID(v UTXOInputMetricItem) {
	o.OutPullOutputByID = v
}

// GetOutPullTxInclusionState returns the OutPullTxInclusionState field value
func (o *ChainMessageMetrics) GetOutPullTxInclusionState() TransactionIDMetricItem {
	if o == nil {
		var ret TransactionIDMetricItem
		return ret
	}

	return o.OutPullTxInclusionState
}

// GetOutPullTxInclusionStateOk returns a tuple with the OutPullTxInclusionState field value
// and a boolean to check if the value has been set.
func (o *ChainMessageMetrics) GetOutPullTxInclusionStateOk() (*TransactionIDMetricItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutPullTxInclusionState, true
}

// SetOutPullTxInclusionState sets field value
func (o *ChainMessageMetrics) SetOutPullTxInclusionState(v TransactionIDMetricItem) {
	o.OutPullTxInclusionState = v
}

func (o ChainMessageMetrics) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ChainMessageMetrics) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["inAccountOutput"] = o.InAccountOutput
	toSerialize["inOnLedgerRequest"] = o.InOnLedgerRequest
	toSerialize["inOutput"] = o.InOutput
	toSerialize["inStateOutput"] = o.InStateOutput
	toSerialize["inTxInclusionState"] = o.InTxInclusionState
	toSerialize["outPublishGovernanceTransaction"] = o.OutPublishGovernanceTransaction
	toSerialize["outPublisherStateTransaction"] = o.OutPublisherStateTransaction
	toSerialize["outPullLatestOutput"] = o.OutPullLatestOutput
	toSerialize["outPullOutputByID"] = o.OutPullOutputByID
	toSerialize["outPullTxInclusionState"] = o.OutPullTxInclusionState
	return toSerialize, nil
}

type NullableChainMessageMetrics struct {
	value *ChainMessageMetrics
	isSet bool
}

func (v NullableChainMessageMetrics) Get() *ChainMessageMetrics {
	return v.value
}

func (v *NullableChainMessageMetrics) Set(val *ChainMessageMetrics) {
	v.value = val
	v.isSet = true
}

func (v NullableChainMessageMetrics) IsSet() bool {
	return v.isSet
}

func (v *NullableChainMessageMetrics) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableChainMessageMetrics(val *ChainMessageMetrics) *NullableChainMessageMetrics {
	return &NullableChainMessageMetrics{value: val, isSet: true}
}

func (v NullableChainMessageMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableChainMessageMetrics) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


