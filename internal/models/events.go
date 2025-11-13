package models

import (
	"fmt"
	"github.com/goccy/go-json"
	scalecodec "github.com/itering/scale.go"
)

type Event struct {
	EventID        string                  `json:"event_id"`
	ModuleID       string                  `json:"module_id"`
	EventIndex     int                     `json:"event_idx"`
	ExtrinsicIndex int                     `json:"extrinsic_idx"`
	Params         []scalecodec.EventParam `json:"params"`
}
type ChangesResponse struct {
	Result       any    `json:"result"`
	Subscription string `json:"subscription"`
}

type Result struct {
	Block  string  `json:"block"`
	Change [][]any `json:"changes"`
}
type Changes struct {
	ChangeHash string
}

func GetChangesFromEvents(response *SubscriptionResponse) (*Changes, error) {
	b, err := json.Marshal(response.Params)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	var resp ChangesResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	bb, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	var res Result
	err = json.Unmarshal(bb, &res)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	changes := &Changes{}
	for _, change := range res.Change {
		changes.ChangeHash = change[1].(string)
	}
	return changes, nil
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

type ExtrinsicSuccess struct {
	Class   string  `json:"class"`
	PaysFee string  `json:"pays_fee"`
	Weight  *Weight `json:"weight"`
}
type Weight struct {
	ProofSize int `json:"proof_size"`
	RefTime   int `json:"ref_time"`
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
