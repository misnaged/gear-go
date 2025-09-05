package models

import (
	"fmt"
)

type SubscriptionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`

	Result           any            `json:"result,omitempty"`
	Error            *ResponseError `json:"error,omitempty"`
	SubscriptionHash string
}
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type BasicParams struct {
	Result       any    `json:"result"`
	Subscription string `json:"subscription"`
}

func GetChangesFromEvents(response *SubscriptionResponse) ([]any, error) {
	result, err := GetFieldFromAny("result", response.Params)
	if err != nil {
		return nil, fmt.Errorf("gear.GetFieldFromAny failed: %w", err)

	}
	changes, err := GetFieldFromAny("changes", result.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("gear.GetFieldFromAny failed: %w", err)
	}
	return changes.([]any)[0].([]any), nil //TODO: usually len of the first array is 1. But it would be reasonable to check it to be sure
}

// `author_submitAndWatchExtrinsic` related structs.
// Basically `author_submitAndWatchExtrinsic` request has two scenarios:
// The first one is when no errors have been occurred. In this scenario
// we're getting 4 subscription messages from `author_extrinsicUpdate` method:
// 1) returns subscription hash; 2) returns if subs is ready; 3) returns block hash 4) returns if finalized
//
// In the second scenario an error is received

type InBlockResponse struct {
	InBlock string `json:"inBlock"`
}
type FinalizedResponse struct {
	Finalized string `json:"finalized"`
}
