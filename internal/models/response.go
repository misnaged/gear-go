package models

type RpcGenericResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      any            `json:"id"`
	Result  any            `json:"result"`
	Method  string         `json:"method"`
	Error   map[string]any `json:"error,omitempty"`
}
