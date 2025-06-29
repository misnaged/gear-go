/*
In this example we're uploading new code using Extrinsic upload_code
and then get its data from the storage
*/

package main

import (
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	gear_storage_methods "github.com/misnaged/gear-go/internal/scale/storage/methods"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}
	code, err := gear.UploadCodeTemp()
	if err != nil {
		logger.Log().Errorf("error uploading code: %v", err)
		os.Exit(1)
	}
	var args []string
	args = append(args, code)
	gear.GetClient().Subscribe(args, "author_submitAndWatchExtrinsic")
	storage := gear_storage_methods.NewStorage("GearProgram", "CodeStorage", gear.GetScale().GetMetadata())
	var vv map[string]any
	err = storage.DecodeStorage(gear.GetRPC(), &vv, true)
	if err != nil {
		logger.Log().Errorf("error decoding storage: %v", err)
		os.Exit(1)
	}

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
