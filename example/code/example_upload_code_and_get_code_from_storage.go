/*
In this example we're uploading new code using Extrinsic upload_code
and then get its data from the storage
*/

package main

import (
	"fmt"
	blake2b2 "github.com/ethereum/go-ethereum/crypto/blake2b"
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_scale "github.com/misnaged/gear-go/internal/scale"
	gear_storage_methods "github.com/misnaged/gear-go/internal/scale/storage/methods"
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
	f, err := os.ReadFile("./example/code/demo_messenger.opt.wasm")
	if err != nil {
		logger.Log().Errorf("failed to read *.wasm file: : %v", err)
		os.Exit(1)
	}
	code, err := uploadCodeTemp(gear.GetScale(), f)
	if err != nil {
		logger.Log().Errorf("error uploading code: %v", err)
		os.Exit(1)
	}

	var args []string
	args = append(args, code)
	gear.GetClient().Subscribe(args, "author_submitAndWatchExtrinsic")
	storage := gear_storage_methods.NewStorage("GearProgram", "CodeStorage", gear.GetScale().GetMetadata())
	key, err := storage.GetStorageKey()
	if err != nil {
		logger.Log().Errorf("error while getting storage key: %v", err)
		os.Exit(1)
	}
	vv, err := storage.DecodeStorageDataMap(gear.GetRPC(), key)
	if err != nil {
		logger.Log().Errorf("error decoding storage: %v", err)
		os.Exit(1)
	}
	//Print CodeId
	h := blake2b2.Sum256(f)
	gear_utils.AddToHex(h[:])

	// CodeId also could be gathered via gear_utils.GetCodeIdFromWasmFile() function

	fmt.Println("Code Id is:", gear_utils.AddToHex(h[:]))

	exports := vv["exports"]
	codeLen := vv["original_code_len"]
	stackEnd := vv["stack_end"]
	version := vv["version"]
	staticPages := vv["static_pages"]
	//
	sections := vv["instantiated_section_sizes"].(map[string]any)
	codeSection := sections["code_section"]
	dataSection := sections["data_section"]
	elementSection := sections["element_section"]
	globalSection := sections["global_section"]
	tableSection := sections["table_section"]
	typeSection := sections["type_section"]
	fmt.Printf(`
		exports %v 
		original_code_len %v 
		stack_end %v 
		version %v 
		static_pages %v 
		code_section %v 
		data_section %v 
		element_section %v 
		global_section %v 
		table_section %v 
		type_section %v 
		`,
		exports,
		codeLen,
		stackEnd,
		version,
		staticPages,
		codeSection,
		dataSection,
		elementSection,
		globalSection,
		tableSection,
		typeSection)

}
func uploadCodeTemp(scale *gear_scale.Scale, file []byte) (string, error) {

	var args []any

	toHex := gear_utils.AddToHex(file)
	args = append(args, toHex)
	params, err := extrinsic_params.InitBuilder("Gear", "upload_code", scale.GetMetadata().Metadata.Modules, args)
	if err != nil {
		return "", fmt.Errorf(" extrinsic_params.InitBuilder failed: %w", err)
	}

	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	aa, err := scale.SignTransaction("Gear", "upload_code", kr, params)
	if err != nil {
		return "", fmt.Errorf(" gear.scale.SignTransaction failed: %w", err)
	}
	//TODO: debug only. Removal is needed
	return aa, nil
}
