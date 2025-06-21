package gear_rpc

import (
	"github.com/misnaged/gear-go/internal/models"
)

type IGearRPC interface {
	NoArgRpcRequest(method NoArgsMethod) (*models.RpcGenericResponse, error)
	// Chain methods
	ChainGetBlockHashLatest() (*models.RpcGenericResponse, error)

	//State methods
	StateGetRuntimeVersionLatest() (*models.RpcGenericResponse, error)
	StateGetMetadataLatest() (*models.RpcGenericResponse, error)
	StateGetRuntimeVersion(blockHash string) (*models.RpcGenericResponse, error)
	StateGetMetadata(blockHash string) (*models.RpcGenericResponse, error)
	StateGetStorageLatest(accountId string) (*models.RpcGenericResponse, error)
}
