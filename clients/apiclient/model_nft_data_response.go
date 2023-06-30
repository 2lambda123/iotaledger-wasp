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

// checks if the NFTDataResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NFTDataResponse{}

// NFTDataResponse struct for NFTDataResponse
type NFTDataResponse struct {
	Id       string `json:"id"`
	Issuer   string `json:"issuer"`
	Metadata string `json:"metadata"`
	Owner    string `json:"owner"`
}

// NewNFTDataResponse instantiates a new NFTDataResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNFTDataResponse(id string, issuer string, metadata string, owner string) *NFTDataResponse {
	this := NFTDataResponse{}
	this.Id = id
	this.Issuer = issuer
	this.Metadata = metadata
	this.Owner = owner
	return &this
}

// NewNFTDataResponseWithDefaults instantiates a new NFTDataResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNFTDataResponseWithDefaults() *NFTDataResponse {
	this := NFTDataResponse{}
	return &this
}

// GetId returns the Id field value
func (o *NFTDataResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *NFTDataResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *NFTDataResponse) SetId(v string) {
	o.Id = v
}

// GetIssuer returns the Issuer field value
func (o *NFTDataResponse) GetIssuer() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Issuer
}

// GetIssuerOk returns a tuple with the Issuer field value
// and a boolean to check if the value has been set.
func (o *NFTDataResponse) GetIssuerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Issuer, true
}

// SetIssuer sets field value
func (o *NFTDataResponse) SetIssuer(v string) {
	o.Issuer = v
}

// GetMetadata returns the Metadata field value
func (o *NFTDataResponse) GetMetadata() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *NFTDataResponse) GetMetadataOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *NFTDataResponse) SetMetadata(v string) {
	o.Metadata = v
}

// GetOwner returns the Owner field value
func (o *NFTDataResponse) GetOwner() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Owner
}

// GetOwnerOk returns a tuple with the Owner field value
// and a boolean to check if the value has been set.
func (o *NFTDataResponse) GetOwnerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Owner, true
}

// SetOwner sets field value
func (o *NFTDataResponse) SetOwner(v string) {
	o.Owner = v
}

func (o NFTDataResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NFTDataResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["issuer"] = o.Issuer
	toSerialize["metadata"] = o.Metadata
	toSerialize["owner"] = o.Owner
	return toSerialize, nil
}

type NullableNFTDataResponse struct {
	value *NFTDataResponse
	isSet bool
}

func (v NullableNFTDataResponse) Get() *NFTDataResponse {
	return v.value
}

func (v *NullableNFTDataResponse) Set(val *NFTDataResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableNFTDataResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableNFTDataResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNFTDataResponse(val *NFTDataResponse) *NullableNFTDataResponse {
	return &NullableNFTDataResponse{value: val, isSet: true}
}

func (v NullableNFTDataResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNFTDataResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
