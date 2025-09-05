package main

import (
	"fmt"
	blake2b2 "github.com/ethereum/go-ethereum/crypto/blake2b"
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/internal/calls"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"os"
)

func main() {

	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}
	gear.GetConfig().Client.IsWebSocket = true // override if it's false
	f, err := os.ReadFile("./example/code/demo_messenger.opt.wasm")
	if err != nil {
		logger.Log().Errorf("failed to read *.wasm file: : %v", err)
		os.Exit(1)
	}
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	call := calls.NewCalls(gear.GetMeta(), gear.GetRPC(), kr)
	code, err := uploadCodeTemp(call, f)
	if err != nil {
		logger.Log().Errorf("error uploading code: %v", err)
		os.Exit(1)
	}
	var args []string
	args = append(args, code)
	_, err = gear.GetClient().PostRequest(args, "author_submitExtrinsic")
	if err != nil {
		logger.Log().Errorf("error posting request: %v", err)
		os.Exit(1)
	}
	storage := gear_storage_methods.NewStorage("GearProgram", "CodeMetadataStorage", gear.GetMeta(), gear.GetRPC())
	storageDataArr, err := storage.DecodeStorageDataArray()
	if err != nil {
		logger.Log().Errorf("error decoding storage: %v", err)
		os.Exit(1)
	}
	storageData := storageDataArr[0]

	h := blake2b2.Sum256(f)
	gear_utils.AddToHex(h[:])

	exports := storageData["exports"]
	origCodeLen := storageData["original_code_len"]
	stackEnd := storageData["stack_end"]
	staticPages := storageData["static_pages"]
	//
	sections := storageData["instrumentation_status"].(map[string]any)
	instrumented := sections["Instrumented"].(map[string]any)
	version := instrumented["version"]
	codeLen := instrumented["code_len"]
	fmt.Printf(`
				exports %v
				original_code_len %v
				stack_end %v
				static_pages %v
				instrumentationStatus.Instrumented.version %v
				instrumentationStatus.Instrumented.codelen %v
				`,

		exports,
		origCodeLen,
		stackEnd,
		staticPages,
		version,
		codeLen)

}
func uploadCodeTemp(calls *calls.Calls, file []byte) (string, error) {

	var args []any

	toHex := gear_utils.AddToHex(file)
	args = append(args, toHex)
	params, err := extrinsic_params.InitBuilder("Gear", "upload_code", calls.Meta.GetMetadata().Metadata.Modules, args)
	if err != nil {
		return "", fmt.Errorf(" extrinsic_params.InitBuilder failed: %w", err)
	}

	aa, err := calls.SignTransaction("Gear", "upload_code", params)
	if err != nil {
		return "", fmt.Errorf(" gear.scale.SignTransaction failed: %w", err)
	}
	return aa, nil
}
