package models

import "github.com/goccy/go-json"

type RpcGenericRequest struct {
	Jsonrpc any    `json:"jsonrpc"`
	Id      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

func (rpc *RpcGenericRequest) MarshalBody() ([]byte, error) {
	return json.Marshal(rpc)
}
