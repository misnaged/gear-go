package gear_calls

import (
	"fmt"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
)

func (gc *GearCalls) UploadCode(pathToWasm string) error {
	codeParams, err := extrinsic_params.NewGearCode(pathToWasm)
	if err != nil {
		return fmt.Errorf("error creating code: %w", err)
	}
	args := []any{codeParams.Code}
	call, err := gc.c.CallBuilder("upload_code", args)
	if err != nil {
		return fmt.Errorf("error calling extrinsic params builder: %w", err)
	}
	resp, err := gc.c.GearRpc.AuthorSubmitExtrinsic(call)
	if err != nil {
		return fmt.Errorf("error submitting extrinsic: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("chain response returned error: %v", resp.Error)
	}
	return nil
}
