package gear_calls

import (
	"fmt"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
)

func (gc *GearCalls) CreateProgram(pathToWasm string, p *extrinsic_params.GearProgram) error {
	codeId, err := gear_utils.GetCodeIdFromWasmFile(pathToWasm)
	if err != nil {
		return fmt.Errorf("gear_utils.GetCodeIdFromWasmFile: %v", err)
	}

	p.CodeId = codeId
	args := []any{p.CodeId, p.Salt, p.InitPayload, p.GasLimit, p.Value, p.KeepAlive}
	call, err := gc.c.CallBuilder("create_program", args)
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
	fmt.Println(resp.Result)
	return nil

}
