package main

import (
	gear_go "github.com/misnaged/gear-go"
	gear_calls "github.com/misnaged/gear-go/internal/calls/gear"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

const wasmPath = "assets/wasm/test/message.opt.wasm"

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}

	var types = []string{"send_message", "subscribeStorage"}

	err = gear.GetWsClient().AddResponseTypesAndMakeWsConnectionsPool(types...)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	codeID, err := gear_utils.GetCodeIdFromWasmFile(wasmPath)
	if err != nil {
		logger.Log().Errorf("error getting code id: %v", err)
		os.Exit(1)
	}
	storage := gear_storage_methods.NewStorage("GearProgram", "ProgramStorage", gear.GetMeta(), gear.GetRPC())

	program, err := storage.GetActiveProgramByCodeId(codeID)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	/*

		fmt.Println("program id is")
		fmt.Println(program.ProgramId)
		owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
		payload := gear_utils.TextToHex("04")
		resp, err := gear.GetRPC().GearCalculateHandleGas(owner, program.ProgramId, "0x04", 10854583136, false)
		if err != nil {
			logger.Log().Errorf("%v", err)
			os.Exit(1)
		}

		gas, err := gear_utils.GetMinimalGasForProgram(resp)
		if err != nil {
			logger.Log().Errorf("%v", err)
			os.Exit(1)
		}
		fmt.Println(*gas)
	*/
	calls := gear_calls.New(gear.GetCalls())
	hash, err := calls.SendMessage(program.ProgramId, "100000000000000", "0x04", 30845509441, true)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	gear.MergeSubscriptionFunctions(gear.EventsSubscription(), gear.SubmitAndWatchExtrinsic([]any{hash}, "send_message"))
	err = gear.InitSubscriptions()
	if err != nil {
		logger.Log().Errorf(" gear.ProcessEventsSubscription failed: %v", err)
		os.Exit(1)
	}
}
