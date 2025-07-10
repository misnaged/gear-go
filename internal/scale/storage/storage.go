package gear_storage

import gear_rpc "github.com/misnaged/gear-go/internal/rpc"

type IGearStorage interface {
	DecodeStorage(gearRPC gear_rpc.IGearRPC, decodeData any, storkey string) error
	GetStorageKeys(gearRPC gear_rpc.IGearRPC) ([]string, error)
	BuildParams(accountId string) error
	GetStorageKey() (string, error)
	DecodeStorageDataArray(gearRPC gear_rpc.IGearRPC) ([]map[string]any, error)
	DecodeStorageDataMap(gearRPC gear_rpc.IGearRPC, storkey string) (map[string]any, error)

	GetProgramsId(gearRPC gear_rpc.IGearRPC) ([]string, error)
}
