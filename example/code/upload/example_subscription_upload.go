package main

import (
	gear_go "github.com/misnaged/gear-go"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

const wasmPath = "assets/wasm/test/demo_ping.opt.wasm"

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}
	var types = []string{"upload_code", "create_program", "subscribeStorage"}

	// general preparation
	err = gear.GetWsClient().AddResponseTypesAndMakeWsConnectionsPool(types...)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	//  author_submitAndWatchExtrinsic section
	var methods = []string{"author_submitAndWatchExtrinsic", "author_submitAndWatchExtrinsic"}
	var callNames = []string{"upload_code", "create_program"}
	var moduleNames = []string{"Gear", "Gear"}
	// upload_code
	codeId, err := gear_utils.GetCodeIdFromWasmFile(wasmPath)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	upload, err := extrinsic_params.NewGearCode(wasmPath)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	uploadArgs := []string{upload.Code}

	// create_program
	owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
	payload := gear_utils.TextToHex("PING")
	resp, err := gear.GetRPC().GearCalculateInitCreateGas(owner, codeId, payload, 1, true)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	gas, err := gear_utils.GetMinimalGasForProgram(resp)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	p := &extrinsic_params.GearProgram{
		CodeId:      codeId,
		Salt:        "0x1",
		InitPayload: payload,
		GasLimit:    *gas,
		Value:       "1",
		KeepAlive:   true,
	}

	createArgs := []any{p.CodeId, p.Salt, p.InitPayload, p.GasLimit, p.Value, p.KeepAlive}

	var args [][]any

	args = append(args, strToAny(uploadArgs), createArgs)
	gear.MergeSubscriptionFunctions(gear.EventsSubscription(), gear.EnqueuedSubscriptions(methods, strToRespType(types), callNames, moduleNames, args))
	err = gear.InitSubscriptions()
	if err != nil {
		logger.Log().Errorf(" gear.ProcessEventsSubscription failed: %v", err)
		os.Exit(1)
	}
}

func strToAny(aArr []string) []any {
	var strArray []any
	for _, a := range aArr {
		strArray = append(strArray, a)
	}
	return strArray
}

func strToRespType(aArr []string) []gear_client.ResponseType {
	var respTypes []gear_client.ResponseType
	for _, a := range aArr {
		respTypes = append(respTypes, gear_client.ResponseType(a))
	}
	return respTypes
}
