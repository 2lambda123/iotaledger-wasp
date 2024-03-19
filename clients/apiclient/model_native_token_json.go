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

// checks if the NativeTokenJSON type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NativeTokenJSON{}

// NativeTokenJSON struct for NativeTokenJSON
type NativeTokenJSON struct {
	Amount string `json:"amount"`
	Id string `json:"id"`
}

// NewNativeTokenJSON instantiates a new NativeTokenJSON object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNativeTokenJSON(amount string, id string) *NativeTokenJSON {
	this := NativeTokenJSON{}
	this.Amount = amount
	this.Id = id
	return &this
}

// NewNativeTokenJSONWithDefaults instantiates a new NativeTokenJSON object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNativeTokenJSONWithDefaults() *NativeTokenJSON {
	this := NativeTokenJSON{}
	return &this
}

// GetAmount returns the Amount field value
func (o *NativeTokenJSON) GetAmount() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Amount
}

// GetAmountOk returns a tuple with the Amount field value
// and a boolean to check if the value has been set.
func (o *NativeTokenJSON) GetAmountOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Amount, true
}

// SetAmount sets field value
func (o *NativeTokenJSON) SetAmount(v string) {
	o.Amount = v
}

// GetId returns the Id field value
func (o *NativeTokenJSON) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *NativeTokenJSON) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *NativeTokenJSON) SetId(v string) {
	o.Id = v
}

func (o NativeTokenJSON) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NativeTokenJSON) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["amount"] = o.Amount
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

type NullableNativeTokenJSON struct {
	value *NativeTokenJSON
	isSet bool
}

func (v NullableNativeTokenJSON) Get() *NativeTokenJSON {
	return v.value
}

func (v *NullableNativeTokenJSON) Set(val *NativeTokenJSON) {
	v.value = val
	v.isSet = true
}

func (v NullableNativeTokenJSON) IsSet() bool {
	return v.isSet
}

func (v *NullableNativeTokenJSON) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNativeTokenJSON(val *NativeTokenJSON) *NullableNativeTokenJSON {
	return &NullableNativeTokenJSON{value: val, isSet: true}
}

func (v NullableNativeTokenJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNativeTokenJSON) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

