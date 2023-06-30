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

// checks if the ContractCallViewRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ContractCallViewRequest{}

// ContractCallViewRequest struct for ContractCallViewRequest
type ContractCallViewRequest struct {
	Arguments JSONDict `json:"arguments"`
	// The contract name as HName (Hex)
	ContractHName string `json:"contractHName"`
	// The contract name
	ContractName string `json:"contractName"`
	// The function name as HName (Hex)
	FunctionHName string `json:"functionHName"`
	// The function name
	FunctionName string `json:"functionName"`
}

// NewContractCallViewRequest instantiates a new ContractCallViewRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewContractCallViewRequest(arguments JSONDict, contractHName string, contractName string, functionHName string, functionName string) *ContractCallViewRequest {
	this := ContractCallViewRequest{}
	this.Arguments = arguments
	this.ContractHName = contractHName
	this.ContractName = contractName
	this.FunctionHName = functionHName
	this.FunctionName = functionName
	return &this
}

// NewContractCallViewRequestWithDefaults instantiates a new ContractCallViewRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewContractCallViewRequestWithDefaults() *ContractCallViewRequest {
	this := ContractCallViewRequest{}
	return &this
}

// GetArguments returns the Arguments field value
func (o *ContractCallViewRequest) GetArguments() JSONDict {
	if o == nil {
		var ret JSONDict
		return ret
	}

	return o.Arguments
}

// GetArgumentsOk returns a tuple with the Arguments field value
// and a boolean to check if the value has been set.
func (o *ContractCallViewRequest) GetArgumentsOk() (*JSONDict, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Arguments, true
}

// SetArguments sets field value
func (o *ContractCallViewRequest) SetArguments(v JSONDict) {
	o.Arguments = v
}

// GetContractHName returns the ContractHName field value
func (o *ContractCallViewRequest) GetContractHName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ContractHName
}

// GetContractHNameOk returns a tuple with the ContractHName field value
// and a boolean to check if the value has been set.
func (o *ContractCallViewRequest) GetContractHNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ContractHName, true
}

// SetContractHName sets field value
func (o *ContractCallViewRequest) SetContractHName(v string) {
	o.ContractHName = v
}

// GetContractName returns the ContractName field value
func (o *ContractCallViewRequest) GetContractName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ContractName
}

// GetContractNameOk returns a tuple with the ContractName field value
// and a boolean to check if the value has been set.
func (o *ContractCallViewRequest) GetContractNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ContractName, true
}

// SetContractName sets field value
func (o *ContractCallViewRequest) SetContractName(v string) {
	o.ContractName = v
}

// GetFunctionHName returns the FunctionHName field value
func (o *ContractCallViewRequest) GetFunctionHName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FunctionHName
}

// GetFunctionHNameOk returns a tuple with the FunctionHName field value
// and a boolean to check if the value has been set.
func (o *ContractCallViewRequest) GetFunctionHNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FunctionHName, true
}

// SetFunctionHName sets field value
func (o *ContractCallViewRequest) SetFunctionHName(v string) {
	o.FunctionHName = v
}

// GetFunctionName returns the FunctionName field value
func (o *ContractCallViewRequest) GetFunctionName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FunctionName
}

// GetFunctionNameOk returns a tuple with the FunctionName field value
// and a boolean to check if the value has been set.
func (o *ContractCallViewRequest) GetFunctionNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FunctionName, true
}

// SetFunctionName sets field value
func (o *ContractCallViewRequest) SetFunctionName(v string) {
	o.FunctionName = v
}

func (o ContractCallViewRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ContractCallViewRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["arguments"] = o.Arguments
	toSerialize["contractHName"] = o.ContractHName
	toSerialize["contractName"] = o.ContractName
	toSerialize["functionHName"] = o.FunctionHName
	toSerialize["functionName"] = o.FunctionName
	return toSerialize, nil
}

type NullableContractCallViewRequest struct {
	value *ContractCallViewRequest
	isSet bool
}

func (v NullableContractCallViewRequest) Get() *ContractCallViewRequest {
	return v.value
}

func (v *NullableContractCallViewRequest) Set(val *ContractCallViewRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableContractCallViewRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableContractCallViewRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableContractCallViewRequest(val *ContractCallViewRequest) *NullableContractCallViewRequest {
	return &NullableContractCallViewRequest{value: val, isSet: true}
}

func (v NullableContractCallViewRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableContractCallViewRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
