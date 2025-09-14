package models

type RpcGenericResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      any            `json:"id"`
	Result  any            `json:"result,omitempty"`
	Method  string         `json:"method"`
	Error   map[string]any `json:"error,omitempty"`
}

type GasCalculateResult struct {
	MinLimit      int  `json:"min_limit"`
	Reserved      int  `json:"reserved"`
	Burned        int  `json:"burned"`
	MayBeReturned int  `json:"may_be_returned"`
	Waited        bool `json:"waited"`
}
