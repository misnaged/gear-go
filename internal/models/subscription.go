package models

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

func IsFinalized(resp *SubscriptionResponse) bool {
	if resultMap, ok := resp.Params.(map[string]any); ok {
		if n, ok := resultMap["result"]; ok {
			if rss, ok := n.(map[string]any); ok {
				if _, ok = rss["finalized"]; ok {
					return true
				}
			}
		}
	}
	return false
}
