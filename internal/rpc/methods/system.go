package gear_rpc_method

import (
	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/substrate-api-rpc/util/ss58"
)

func (gearRPC *GearRpc) SystemAccountNextIndex(accountId string) (*models.RpcGenericResponse, error) {
	var params []string
	//TODO: add switching to Gear's addressType
	params = append(params, ss58.Encode(accountId, 42))

	return gearRPC.client.PostRequest(params, "system_accountNextIndex")
}
