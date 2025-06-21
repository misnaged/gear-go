package gear_rpc_method

import "github.com/misnaged/gear-go/internal/models"

// StateGetRuntimeVersionLatest returns response for `state_getRuntimeVersion` json-rpc request with latest blockhash
func (gearRPC *GearRpc) StateGetRuntimeVersionLatest() (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(nil, "state_getRuntimeVersion")
}

// StateGetMetadataLatest returns response for json-rpc `state_getMetadata` request with latest blockhash
func (gearRPC *GearRpc) StateGetMetadataLatest() (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(nil, "state_getMetadata")
}

func (gearRPC *GearRpc) StateGetRuntimeVersion(blockHash string) (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(blockHash, "state_getRuntimeVersion")
}

func (gearRPC *GearRpc) StateGetMetadata(blockHash string) (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(blockHash, "state_getMetadata")
}

func (gearRPC *GearRpc) StateGetStorageLatest(accountId string) (*models.RpcGenericResponse, error) {
	var params []string
	params = append(params, accountId)
	return gearRPC.client.PostRequest(params, "state_getStorage")
}
