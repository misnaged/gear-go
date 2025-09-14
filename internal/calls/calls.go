package calls

import (
	"errors"
	"fmt"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/misnaged/substrate-api-rpc/rpc"
)

type Calls struct {
	GearRpc    gear_rpc.IGearRPC
	Meta       *metadata.Metadata
	KeyRing    keyring.IKeyRing
	customTx   rpc.ICustomTranscation
	ModuleName string
}

func NewCalls(meta *metadata.Metadata, gearRpc gear_rpc.IGearRPC, kr keyring.IKeyRing) *Calls {
	return &Calls{
		Meta:    meta,
		GearRpc: gearRpc,
		KeyRing: kr,
	}
}
func (calls *Calls) CallBuilder(callName string, args []any) (string, error) {
	if calls.ModuleName == "" {
		return "", fmt.Errorf("%w", errors.New("module name is empty"))
	}
	params, err := extrinsic_params.InitBuilder(calls.ModuleName, callName, calls.Meta.GetMetadata().Metadata.Modules, args)
	if err != nil {
		return "", fmt.Errorf(" extrinsic_params.InitBuilder failed: %w", err)
	}

	aa, err := calls.SignTransaction(calls.ModuleName, callName, params)
	if err != nil {
		return "", fmt.Errorf("SignTransaction failed: %w", err)
	}
	return aa, nil
}

func (calls *Calls) DoCall(callHash string) error {
	resp, err := calls.GearRpc.AuthorSubmitExtrinsic(callHash)
	if err != nil {
		return fmt.Errorf("error submitting extrinsic: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("chain response returned error: %v", resp.Error)
	}
	return nil
}
