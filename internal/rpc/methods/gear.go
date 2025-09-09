package gear_rpc_method

import "github.com/misnaged/gear-go/internal/models"

// gear_calculateInitCreateGas
func (gearRPC *GearRpc) GearCalculateInitCreateGas(owner, codeId, payload string, value any, allowOtherPanic bool) (*models.RpcGenericResponse, error) {
	var params []any
	params = append(params, codeId, owner, payload, value, allowOtherPanic)
	return gearRPC.client.PostRequest(params, "gear_calculateInitCreateGas")
}
