package main

import (
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

const BobAccountId = "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("NewGear err: %v", err)
		os.Exit(1)
	}

	storage := gear_storage_methods.NewStorage("System", "Account", gear.GetMeta(), gear.GetRPC())

	err = storage.AddAccountIdToStorageParams(BobAccountId)
	if err != nil {
		logger.Log().Errorf("BuildParams err: %v", err)
		os.Exit(1)
	}

	storkey, err := storage.GetStorageKey()
	if err != nil {
		logger.Log().Errorf(" gear.GetStorageKey failed: %v", err)
		os.Exit(1)
	}

	arr, err := storage.DecodeStorageDataMap(storkey)
	if err != nil {
		logger.Log().Errorf(" gear.DecodeStorage failed: %v", err)
		os.Exit(1)
	}
	fmt.Println(arr)
}
