package models

import "github.com/goccy/go-json"

type RpcGenericRequest struct {
	Jsonrpc any         `json:"jsonrpc"`
	Id      any         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (rpc *RpcGenericRequest) MarshalBody() ([]byte, error) {
	body := &RpcGenericRequest{
		Jsonrpc: rpc.Jsonrpc,
		Id:      rpc.Id,
		Method:  rpc.Method,
		Params:  rpc.Params,
	}
	return json.Marshal(body)
}
