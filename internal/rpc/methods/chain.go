package gear_rpc_method

import "github.com/misnaged/gear-go/internal/models"

// ChainGetBlockHashLatest returns response for json-rpc request `chain_getBlockHash` with the latest blockNumber
func (gearRPC *GearRpc) ChainGetBlockHashLatest() (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(nil, "chain_getBlockHash")
}
