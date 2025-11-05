package extrinsic_params

import (
	"fmt"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"os"
)

type GearCode struct {
	Code string // rust: code Vec<U8>
}

type GearSendMessage struct {
	ProgramId string
	Payload   string
	GasLimit  int
	Value     string
	KeepAlive bool
}

type GearProgram struct {
	CodeId      string // rust: code_id [U8; 32]
	Salt        string // rust: salt Vec<U8>
	InitPayload string // rust: init_payload Vec<U8>
	GasLimit    int    // rust: gas_limit U64
	Value       string // rust: value U128
	KeepAlive   bool   // keep_alive Bool
}

func NewGearCode(wasmFilePath string) (*GearCode, error) {
	f, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read wasm file: %v", err)
	}
	toHex := gear_utils.AddToHex(f)
	return &GearCode{Code: toHex}, nil

}
