package gear_storage_methods

import (
	"errors"
	"fmt"
	"github.com/itering/scale.go/types"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	gear_storage "github.com/misnaged/gear-go/internal/scale/storage"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/substrate-api-rpc/storage"
	"github.com/misnaged/substrate-api-rpc/storageKey"
	"github.com/misnaged/substrate-api-rpc/util/twox"
)

type Storage struct {
	moduleName,
	methodName,
	scaleType string
	meta   *types.MetadataStruct
	params []any
}

func NewStorage(moduleName, methodName string, meta *types.MetadataStruct) gear_storage.IGearStorage {
	return &Storage{moduleName: moduleName, methodName: methodName, meta: meta}
}
func (stor *Storage) encodeModuleAndMethodNames() []byte {
	module := twox.NewXXHash128([]byte(stor.moduleName))
	method := twox.NewXXHash128([]byte(stor.methodName))

	return append(append([]byte{}, module[:]...), method[:]...)
}

func (stor *Storage) getScaleType() (*storageKey.StorageOption, error) {
	storageType, err := gear_utils.GetStorageTypeByModuleAndMethodNames(stor.moduleName, stor.methodName, stor.meta.Metadata.Modules)
	if err != nil {
		return nil, fmt.Errorf("error getting storage type: %v", err)
	}
	return storageKey.CheckoutHasherAndType(storageType), nil
}
func (stor *Storage) getTypeName() error {
	name, err := stor.getScaleType()
	if err != nil {
		return fmt.Errorf("error getting scale type: %v", err)
	}
	stor.scaleType = name.Value
	return nil
}

func (stor *Storage) getEncodedStorageKey() (string, error) {
	b := stor.encodeModuleAndMethodNames()
	if stor.params != nil {
		if len(stor.params) > 0 {
			for _, param := range stor.params {
				if param.([]byte) == nil {
					return gear_utils.AddToHex(b), nil
				} else {
					b = append(b, param.([]byte)...)
				}
			}
			return gear_utils.AddToHex(b), nil
		}
	}
	return gear_utils.AddToHex(b), nil
}

func (stor *Storage) GetStorageKey() (string, error) {
	return stor.getEncodedStorageKey()
}

func (stor *Storage) DecodeStorageDataArray(gearRPC gear_rpc.IGearRPC) ([]map[string]any, error) {
	keys, err := stor.getPagedKeys(gearRPC)
	if err != nil {
		return nil, fmt.Errorf("error getting paged keys: %v", err)
	}

	if len(keys) <= 0 {
		return nil, fmt.Errorf("%w", errors.New("storageKeys length is 0"))
	}

	storageDataArr := make([]map[string]any, len(keys))

	for i := range keys {
		m, err := stor.DecodeStorageDataMap(gearRPC, keys[i])
		if err != nil {
			return nil, fmt.Errorf("error decoding storage data map: %v", err)
		}
		storageDataArr = append(storageDataArr, m)
	}
	return storageDataArr, nil
}

// TODO: refactoring is needed (shall add enum in the next updates)

func (stor *Storage) getPagedKeys(gearRPC gear_rpc.IGearRPC) ([]string, error) {
	storKey, err := stor.getEncodedStorageKey()
	if err != nil {
		return nil, fmt.Errorf(" gear.scale.StorageRequest failed: %v", err)
	}
	var pagedKeys []string

	keyPaged, err := gearRPC.StateGetKeyPaged(storKey)
	if err != nil {
		return nil, fmt.Errorf(" gear.scale.StateGetStorageLatest failed: %v", err)
	}
	toAnyArr := keyPaged.Result.([]any)
	for i := range toAnyArr {
		pagedKeys = append(pagedKeys, toAnyArr[i].(string))
	}
	return pagedKeys, nil

}

func (stor *Storage) GetStorageKeys(gearRPC gear_rpc.IGearRPC) ([]string, error) {
	return stor.getPagedKeys(gearRPC)
}
func (stor *Storage) getStorageRpc(gearRPC gear_rpc.IGearRPC, storkey string) (string, error) {
	resp, err := gearRPC.StateGetStorageLatest(storkey)
	if err != nil {
		return "", fmt.Errorf(" gear.scale.StateGetStorageLatest failed: %v", err)
	}
	if resp.Result == nil {
		return "", errors.New("response result is nil")
	}
	return resp.Result.(string), nil
}

func (stor *Storage) DecodeStorageDataMap(gearRPC gear_rpc.IGearRPC, storkey string) (map[string]any, error) {
	storageEncoded, err := stor.getStorageRpc(gearRPC, storkey)
	if err != nil {
		return nil, fmt.Errorf(" gear.scale.GetStorageRpc failed: %v", err)
	}
	err = stor.getTypeName()
	if err != nil {
		return nil, fmt.Errorf("getTypeName failed: %v", err)
	}
	a, _, err := storage.Decode(storageEncoded, stor.scaleType, &types.ScaleDecoderOption{Metadata: stor.meta})
	if err != nil {
		return nil, fmt.Errorf("storage.Decode failed: %v", err)
	}

	return a.ToMapInterface(), nil
}
func (stor *Storage) DecodeStorage(gearRPC gear_rpc.IGearRPC, decodeData any, storkey string) error {
	storageEncoded, err := stor.getStorageRpc(gearRPC, storkey)
	if err != nil {
		return fmt.Errorf(" gear.scale.GetStorageRpc failed: %v", err)
	}
	err = stor.getTypeName()
	if err != nil {
		return fmt.Errorf("getTypeName failed: %v", err)
	}
	a, _, err := storage.Decode(storageEncoded, stor.scaleType, &types.ScaleDecoderOption{Metadata: stor.meta})
	if err != nil {
		return fmt.Errorf("storage.Decode failed: %v", err)
	}
	a.ToAny(&decodeData)
	return nil
}
