package gear_utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	blake2b2 "github.com/ethereum/go-ethereum/crypto/blake2b"
	"github.com/goccy/go-json"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models"
	"os"
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
func GetExtrinsicDecoderByRawHex(raw string, meta *types.MetadataStruct) *scalecodec.ExtrinsicDecoder {
	option := types.ScaleDecoderOption{Metadata: meta}
	dec := scalecodec.ExtrinsicDecoder{}
	dec.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
	dec.Process()
	return &dec
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
				if strings.EqualFold(stor.Name, methodName) {
					return &stor.Type, nil
				}
			}
		}
	}
	return nil, errors.New("module not found")
}

func GetArgumentsForExtrinsicParams(moduleName, callName string, modules []types.MetadataModules) []types.MetadataModuleCallArgument {
	for _, v := range modules {
		if v.Name == moduleName {
			for _, vv := range v.Calls {
				if vv.Name == callName {
					return vv.Args
				}
			}
		}
	}
	return nil
}

// GetCodeIdFromWasmFile reads wasm file and returns it Blake2 hash
//
// https://docs.rs/gprimitives/latest/gprimitives/struct.CodeId.html
func GetCodeIdFromWasmFile(filePath string) (string, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	h := blake2b2.Sum256(f)
	return AddToHex(h[:]), nil
}

func TextToHex(text string) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString([]byte(text)))
}

func DecodeToString(str string) (string, error) {
	res, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
func GetMinimalGasForProgram(calculateGasResponse *models.RpcGenericResponse) (*int, error) {
	b, err := json.Marshal(calculateGasResponse.Result)
	if err != nil {
		return nil, fmt.Errorf(" json.Marshal failed: %w", err)
	}
	var calculateGasResult models.GasCalculateResult
	err = json.Unmarshal(b, &calculateGasResult)
	if err != nil {
		return nil, fmt.Errorf(" json.Unmarshal failed: %w", err)
	}
	return &calculateGasResult.MinLimit, nil
}
