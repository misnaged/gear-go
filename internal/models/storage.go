package models

type FrameSystemAccountInfo struct {
	Nonce       any   `json:"nonce"`
	Consumers   any   `json:"consumers"`
	Provider    any   `json:"providers"`
	Sufficients any   `json:"sufficients"`
	Data        *Data `json:"data"`
}
type Data struct {
	Free     any `json:"free"`
	Reserved any `json:"reserved"`
	Frozen   any `json:"frozen"`
	Flags    any `json:"flags"`
}
