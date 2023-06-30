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

// checks if the CallTarget type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CallTarget{}

// CallTarget struct for CallTarget
type CallTarget struct {
	// The contract name as HName (Hex)
	ContractHName string `json:"contractHName"`
	// The function name as HName (Hex)
	FunctionHName string `json:"functionHName"`
}

// NewCallTarget instantiates a new CallTarget object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCallTarget(contractHName string, functionHName string) *CallTarget {
	this := CallTarget{}
	this.ContractHName = contractHName
	this.FunctionHName = functionHName
	return &this
}

// NewCallTargetWithDefaults instantiates a new CallTarget object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCallTargetWithDefaults() *CallTarget {
	this := CallTarget{}
	return &this
}

// GetContractHName returns the ContractHName field value
func (o *CallTarget) GetContractHName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ContractHName
}

// GetContractHNameOk returns a tuple with the ContractHName field value
// and a boolean to check if the value has been set.
func (o *CallTarget) GetContractHNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ContractHName, true
}

// SetContractHName sets field value
func (o *CallTarget) SetContractHName(v string) {
	o.ContractHName = v
}

// GetFunctionHName returns the FunctionHName field value
func (o *CallTarget) GetFunctionHName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FunctionHName
}

// GetFunctionHNameOk returns a tuple with the FunctionHName field value
// and a boolean to check if the value has been set.
func (o *CallTarget) GetFunctionHNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FunctionHName, true
}

// SetFunctionHName sets field value
func (o *CallTarget) SetFunctionHName(v string) {
	o.FunctionHName = v
}

func (o CallTarget) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CallTarget) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["contractHName"] = o.ContractHName
	toSerialize["functionHName"] = o.FunctionHName
	return toSerialize, nil
}

type NullableCallTarget struct {
	value *CallTarget
	isSet bool
}

func (v NullableCallTarget) Get() *CallTarget {
	return v.value
}

func (v *NullableCallTarget) Set(val *CallTarget) {
	v.value = val
	v.isSet = true
}

func (v NullableCallTarget) IsSet() bool {
	return v.isSet
}

func (v *NullableCallTarget) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCallTarget(val *CallTarget) *NullableCallTarget {
	return &NullableCallTarget{value: val, isSet: true}
}

func (v NullableCallTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCallTarget) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
