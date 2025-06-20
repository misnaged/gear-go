package gear_utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"strings"
)

func GetMetaData(cli gear_client.IClient) (*types.MetadataStruct, error) {
	decoder := &scalecodec.MetadataDecoder{}
	var postReq any
	resp, err := cli.PostRequest(postReq, "state_getMetadata")
	if err != nil {
		return nil, fmt.Errorf("post request err: %v", err)
	}

	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &decoder.Metadata, nil
}

func ModulesMap(meta *types.MetadataStruct) map[string]*types.MetadataModules {
	m := make(map[string]*types.MetadataModules)

	for i := range meta.Metadata.Modules {
		m[meta.Metadata.Modules[i].Name] = &meta.Metadata.Modules[i]
	}
	return m
}

func CallsMapByModuleName(meta *types.MetadataStruct, moduleName string) map[string]*types.MetadataCalls {
	m := make(map[string]*types.MetadataCalls)
	for i := range meta.Metadata.Modules {
		if meta.Metadata.Modules[i].Name == moduleName {
			for ii := range meta.Metadata.Modules[i].Calls {
				m[meta.Metadata.Modules[i].Calls[ii].Name] = &meta.Metadata.Modules[i].Calls[ii]
			}
		}
	}
	return m
}

/*

	Do not use GetCallLookupIndexByModuleCallName and GetCallArgsByModuleCallName until you 100% sure it doesn't have any dupes

	Use GetCallLookupIndexByModuleAndCallNames and GetCallArgsByModuleAndCallNames instead

	Example: both `ChildBounties` and `Bounties` have `accept_curator` extrinsic. So you would be given `Bounties` one
	most likely
*/

// GetCallLookupIndexByModuleCallName finds extrinction's CallIndex by its CallName
//
//	Warning! It's safer to use GetCallLookupIndexByModuleAndCallNames instead
func GetCallLookupIndexByModuleCallName(meta *types.MetadataStruct, moduleCallName string) string {
	for i := range meta.Metadata.Modules {
		for j := range meta.Metadata.Modules[i].Calls {
			if meta.Metadata.Modules[i].Calls[j].Name == moduleCallName {
				return meta.Metadata.Modules[i].Calls[j].Lookup
			}
		}
	}
	return ""
}

// GetCallArgsByModuleCallName finds extrinction's Args by its CallName
//
// Warning! It's safer to use GetCallArgsByModuleAndCallNames instead
func GetCallArgsByModuleCallName(meta *types.MetadataStruct, moduleCallName string) (args []types.MetadataModuleCallArgument) {
	for i := range meta.Metadata.Modules {
		for j := range meta.Metadata.Modules[i].Calls {
			if meta.Metadata.Modules[i].Calls[j].Name == moduleCallName {
				args = append(args, meta.Metadata.Modules[i].Calls[j].Args...)
				return
			}
		}
	}
	return nil
}

// GetCallLookupIndexByModuleAndCallNames finds extrinction's LookupIndex by parent ModuleName and its CallName
//
// e.g: ModuleName: `Gear`  CallName: `upload_code`
func GetCallLookupIndexByModuleAndCallNames(meta *types.MetadataStruct, moduleName, callName string) string {
	for i := range meta.Metadata.Modules {
		if meta.Metadata.Modules[i].Name == moduleName {
			for j := range meta.Metadata.Modules[i].Calls {
				if meta.Metadata.Modules[i].Calls[j].Name == callName {
					return meta.Metadata.Modules[i].Calls[j].Lookup
				}
			}
		}
	}
	return ""
}

// GetCallArgsByModuleAndCallNames finds extrinction's Args by parent ModuleName and its CallName
//
// e.g: ModuleName: `Gear`  CallName: `upload_code`
func GetCallArgsByModuleAndCallNames(meta *types.MetadataStruct, moduleName, callName string) (args []types.MetadataModuleCallArgument) {
	for i := range meta.Metadata.Modules {
		if meta.Metadata.Modules[i].Name == moduleName {
			for j := range meta.Metadata.Modules[i].Calls {
				args = append(args, meta.Metadata.Modules[i].Calls[j].Args...)
				return
			}
		}
	}
	return nil
}

func GetStoragesByModuleName(meta *types.MetadataStruct, moduleName string) (storages []types.MetadataStorage) {
	for i := range meta.Metadata.Modules {
		if meta.Metadata.Modules[i].Name == moduleName {
			storages = append(storages, meta.Metadata.Modules[i].Storage...)
			return
		}
	}
	return nil
}

func GetEventsByModuleName(meta *types.MetadataStruct, moduleName string) (events []types.MetadataEvents) {
	for i := range meta.Metadata.Modules {
		if meta.Metadata.Modules[i].Name == moduleName {
			events = append(events, meta.Metadata.Modules[i].Events...)
			return
		}
	}
	return nil
}
func Getttt(meta *types.MetadataStruct) {
	fmt.Println(meta.Extrinsic.SignedExtensions)
}

func AddToHex(addr []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(addr))
}

func GetStorageTypeByModuleAndMethodNames(moduleName, methodName string, modules []types.MetadataModules) (*types.StorageType, error) {
	for _, v := range modules {
		if strings.EqualFold(v.Name, moduleName) {
			for _, stor := range v.Storage {
				fmt.Println(stor.Name)
				if strings.EqualFold(stor.Name, methodName) {
					return &stor.Type, nil
				}
			}
		}
	}
	return nil, errors.New("module not found")
}
