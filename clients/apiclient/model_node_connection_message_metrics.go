/*
Wasp API

REST API for the Wasp node

API version: 123
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"time"
)

// checks if the NodeConnectionMessageMetrics type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NodeConnectionMessageMetrics{}

// NodeConnectionMessageMetrics struct for NodeConnectionMessageMetrics
type NodeConnectionMessageMetrics struct {
	// Last time the message was sent/received
	LastEvent *time.Time `json:"lastEvent,omitempty"`
	// The print out of the last message
	LastMessage *string `json:"lastMessage,omitempty"`
	// Total number of messages sent/received
	Total *uint32 `json:"total,omitempty"`
}

// NewNodeConnectionMessageMetrics instantiates a new NodeConnectionMessageMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodeConnectionMessageMetrics() *NodeConnectionMessageMetrics {
	this := NodeConnectionMessageMetrics{}
	return &this
}

// NewNodeConnectionMessageMetricsWithDefaults instantiates a new NodeConnectionMessageMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodeConnectionMessageMetricsWithDefaults() *NodeConnectionMessageMetrics {
	this := NodeConnectionMessageMetrics{}
	return &this
}

// GetLastEvent returns the LastEvent field value if set, zero value otherwise.
func (o *NodeConnectionMessageMetrics) GetLastEvent() time.Time {
	if o == nil || isNil(o.LastEvent) {
		var ret time.Time
		return ret
	}
	return *o.LastEvent
}

// GetLastEventOk returns a tuple with the LastEvent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodeConnectionMessageMetrics) GetLastEventOk() (*time.Time, bool) {
	if o == nil || isNil(o.LastEvent) {
		return nil, false
	}
	return o.LastEvent, true
}

// HasLastEvent returns a boolean if a field has been set.
func (o *NodeConnectionMessageMetrics) HasLastEvent() bool {
	if o != nil && !isNil(o.LastEvent) {
		return true
	}

	return false
}

// SetLastEvent gets a reference to the given time.Time and assigns it to the LastEvent field.
func (o *NodeConnectionMessageMetrics) SetLastEvent(v time.Time) {
	o.LastEvent = &v
}

// GetLastMessage returns the LastMessage field value if set, zero value otherwise.
func (o *NodeConnectionMessageMetrics) GetLastMessage() string {
	if o == nil || isNil(o.LastMessage) {
		var ret string
		return ret
	}
	return *o.LastMessage
}

// GetLastMessageOk returns a tuple with the LastMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodeConnectionMessageMetrics) GetLastMessageOk() (*string, bool) {
	if o == nil || isNil(o.LastMessage) {
		return nil, false
	}
	return o.LastMessage, true
}

// HasLastMessage returns a boolean if a field has been set.
func (o *NodeConnectionMessageMetrics) HasLastMessage() bool {
	if o != nil && !isNil(o.LastMessage) {
		return true
	}

	return false
}

// SetLastMessage gets a reference to the given string and assigns it to the LastMessage field.
func (o *NodeConnectionMessageMetrics) SetLastMessage(v string) {
	o.LastMessage = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *NodeConnectionMessageMetrics) GetTotal() uint32 {
	if o == nil || isNil(o.Total) {
		var ret uint32
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodeConnectionMessageMetrics) GetTotalOk() (*uint32, bool) {
	if o == nil || isNil(o.Total) {
		return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *NodeConnectionMessageMetrics) HasTotal() bool {
	if o != nil && !isNil(o.Total) {
		return true
	}

	return false
}

// SetTotal gets a reference to the given uint32 and assigns it to the Total field.
func (o *NodeConnectionMessageMetrics) SetTotal(v uint32) {
	o.Total = &v
}

func (o NodeConnectionMessageMetrics) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NodeConnectionMessageMetrics) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.LastEvent) {
		toSerialize["lastEvent"] = o.LastEvent
	}
	if !isNil(o.LastMessage) {
		toSerialize["lastMessage"] = o.LastMessage
	}
	if !isNil(o.Total) {
		toSerialize["total"] = o.Total
	}
	return toSerialize, nil
}

type NullableNodeConnectionMessageMetrics struct {
	value *NodeConnectionMessageMetrics
	isSet bool
}

func (v NullableNodeConnectionMessageMetrics) Get() *NodeConnectionMessageMetrics {
	return v.value
}

func (v *NullableNodeConnectionMessageMetrics) Set(val *NodeConnectionMessageMetrics) {
	v.value = val
	v.isSet = true
}

func (v NullableNodeConnectionMessageMetrics) IsSet() bool {
	return v.isSet
}

func (v *NullableNodeConnectionMessageMetrics) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodeConnectionMessageMetrics(val *NodeConnectionMessageMetrics) *NullableNodeConnectionMessageMetrics {
	return &NullableNodeConnectionMessageMetrics{value: val, isSet: true}
}

func (v NullableNodeConnectionMessageMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodeConnectionMessageMetrics) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


