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

// checks if the RequestIDsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RequestIDsResponse{}

// RequestIDsResponse struct for RequestIDsResponse
type RequestIDsResponse struct {
	RequestIds []string `json:"requestIds"`
}

// NewRequestIDsResponse instantiates a new RequestIDsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequestIDsResponse(requestIds []string) *RequestIDsResponse {
	this := RequestIDsResponse{}
	this.RequestIds = requestIds
	return &this
}

// NewRequestIDsResponseWithDefaults instantiates a new RequestIDsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequestIDsResponseWithDefaults() *RequestIDsResponse {
	this := RequestIDsResponse{}
	return &this
}

// GetRequestIds returns the RequestIds field value
func (o *RequestIDsResponse) GetRequestIds() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.RequestIds
}

// GetRequestIdsOk returns a tuple with the RequestIds field value
// and a boolean to check if the value has been set.
func (o *RequestIDsResponse) GetRequestIdsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.RequestIds, true
}

// SetRequestIds sets field value
func (o *RequestIDsResponse) SetRequestIds(v []string) {
	o.RequestIds = v
}

func (o RequestIDsResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RequestIDsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["requestIds"] = o.RequestIds
	return toSerialize, nil
}

type NullableRequestIDsResponse struct {
	value *RequestIDsResponse
	isSet bool
}

func (v NullableRequestIDsResponse) Get() *RequestIDsResponse {
	return v.value
}

func (v *NullableRequestIDsResponse) Set(val *RequestIDsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestIDsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestIDsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestIDsResponse(val *RequestIDsResponse) *NullableRequestIDsResponse {
	return &NullableRequestIDsResponse{value: val, isSet: true}
}

func (v NullableRequestIDsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestIDsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
