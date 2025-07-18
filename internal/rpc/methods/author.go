package gear_rpc_method

import (
	"github.com/misnaged/gear-go/internal/models"
)

// author_submitExtrinsic
func (gearRPC *GearRpc) AuthorSubmitExtrinsic(signed string) (*models.RpcGenericResponse, error) {
	var params []string
	params = append(params, signed)
	return gearRPC.client.PostRequest(params, "author_submitExtrinsic")
}
