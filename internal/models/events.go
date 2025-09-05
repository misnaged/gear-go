package models

import scalecodec "github.com/itering/scale.go"

type Event struct {
	EventID        string                  `json:"event_id"`
	ModuleID       string                  `json:"module_id"`
	EventIndex     int                     `json:"event_idx"`
	ExtrinsicIndex int                     `json:"extrinsic_idx"`
	Params         []scalecodec.EventParam `json:"params"`
}
