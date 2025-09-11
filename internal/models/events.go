package models

import (
	"fmt"
	scalecodec "github.com/itering/scale.go"
)

type Event struct {
	EventID        string                  `json:"event_id"`
	ModuleID       string                  `json:"module_id"`
	EventIndex     int                     `json:"event_idx"`
	ExtrinsicIndex int                     `json:"extrinsic_idx"`
	Params         []scalecodec.EventParam `json:"params"`
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

type TransactionPaymentEvent struct {
	// TransactionFeePaid
	Who       string
	ActualFee string
	Tip       string
}

type BalancesEvent struct {
	// Withdraw
	// Deposit
	Who    string
	Amount string
}

type TreasuryEvent struct {
	// UpdatedInactive
	Reactivated string
	Deactivated string

	// Deposit
	Value string
}

type Message struct {
	Destination string   `json:"destination"`
	Details     *Details `json:"details"`
	To          string   `json:"to"`
	Id          string   `json:"id"`
	Payload     string   `json:"payload"`
	Source      string   `json:"source"`
	Value       any      `json:"value"`
	Exp         any      `json:"exp"`
}

type Details struct {
	Code *Code `json:"code"`
}
type Code struct {
	Success string        `json:"success,omitempty"`
	Error   *MessageError `json:"error,omitempty"`
}
type MessageError struct {
	Execution string `json:"execution"`
}
