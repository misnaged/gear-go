package gear_rpc

import (
	"github.com/misnaged/gear-go/internal/models"
)

type IGearRPC interface {
	NoArgRpcRequest(method NoArgsMethod) (*models.RpcGenericResponse, error)
	// Author methods
	AuthorSubmitExtrinsic(signed string) (*models.RpcGenericResponse, error)
	// Chain methods
	ChainGetBlockHashLatest() (*models.RpcGenericResponse, error)
	ChainGetBlockHash(blockNum int) (*models.RpcGenericResponse, error)
	//State methods
	StateGetRuntimeVersionLatest() (*models.RpcGenericResponse, error)
	StateGetMetadataLatest() (*models.RpcGenericResponse, error)
	StateGetRuntimeVersion(blockHash string) (*models.RpcGenericResponse, error)
	StateGetMetadata(blockHash string) (*models.RpcGenericResponse, error)
	StateGetStorageLatest(accountId string) (*models.RpcGenericResponse, error)
	StateGetKeyPaged(encodedKey string) (*models.RpcGenericResponse, error)
	StateQueryStorageAt(encodedKey string) (*models.RpcGenericResponse, error)

	SystemAccountNextIndex(accountId string) (*models.RpcGenericResponse, error)

	//Gear methods
	GearCalculateInitCreateGas(codeId, owner, payload string, value any, allowOtherPanic bool) (*models.RpcGenericResponse, error)
}
