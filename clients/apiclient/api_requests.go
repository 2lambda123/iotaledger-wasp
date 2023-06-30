/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RequestsApiService RequestsApi service
type RequestsApiService service

type ApiOffLedgerRequest struct {
	ctx              context.Context
	ApiService       *RequestsApiService
	offLedgerRequest *OffLedgerRequest
}

// Offledger request as JSON. Request encoded in Hex
func (r ApiOffLedgerRequest) OffLedgerRequest(offLedgerRequest OffLedgerRequest) ApiOffLedgerRequest {
	r.offLedgerRequest = &offLedgerRequest
	return r
}

func (r ApiOffLedgerRequest) Execute() (*http.Response, error) {
	return r.ApiService.OffLedgerExecute(r)
}

/*
OffLedger Post an off-ledger request

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiOffLedgerRequest
*/
func (a *RequestsApiService) OffLedger(ctx context.Context) ApiOffLedgerRequest {
	return ApiOffLedgerRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *RequestsApiService) OffLedgerExecute(r ApiOffLedgerRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RequestsApiService.OffLedger")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1/requests/offledger"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.offLedgerRequest == nil {
		return nil, reportError("offLedgerRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.offLedgerRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}
