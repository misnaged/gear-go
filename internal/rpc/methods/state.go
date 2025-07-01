package gear_rpc_method

import (
	"fmt"
	"github.com/misnaged/gear-go/internal/models"
)

//TODO: Add Error handling for all

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

func (gearRPC *GearRpc) StateGetKeyPaged(encodedKey string) (*models.RpcGenericResponse, error) {
	var params []any
	params = append(params, encodedKey, 1000, encodedKey)
	fmt.Println(params)
	return gearRPC.client.PostRequest(params, "state_getKeysPaged")
}

func (gearRPC *GearRpc) StateQueryStorageAt(encodedKey string) (*models.RpcGenericResponse, error) {
	pagedKeys, err := gearRPC.StateGetKeyPaged(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var subParams []any
	var params [][]any
	subParams = append(subParams, pagedKeys.Result.(string)) //TODO panic check for pagedKeys.Result !
	params = append(params, subParams)
	return gearRPC.client.PostRequest(params, "state_queryStorageAt")
}
