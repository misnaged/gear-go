package gear_calls

import "github.com/misnaged/gear-go/internal/models/extrinsic_params"

type IGearCalls interface {
	UploadCode(pathToWasm string) (string, error)
	CreateProgram(pathToWasm string, p *extrinsic_params.GearProgram) (string, error)
	SendMessage(dest, value, payload string, gaslimit int, keepAlive bool) (string, error)
}
