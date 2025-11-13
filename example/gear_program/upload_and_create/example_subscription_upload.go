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

	var args [][]any
	args = append(args, strToAny(uploadArgs))
	inter := gear_go.NewInterruption("create_program", CalculateGas(gear, codeId))
	gear.MergeSubscriptionFunctions(gear.EventsSubscription(), gear.EnqueuedSubscriptions(methods, strToRespType(types), callNames, moduleNames, args, inter))
	err = gear.InitSubscriptions()
	if err != nil {
		logger.Log().Errorf(" gear.ProcessEventsSubscription failed: %v", err)
		os.Exit(1)
	}
}
func CalculateGas(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {
		owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
		payload := gear_utils.TextToHex("0x")
		resp, err := gear.GetRPC().GearCalculateInitCreateGas(owner, codeId, payload, 0, true)
		if err != nil {
			logger.Log().Errorf("%v", err)
			os.Exit(1)
		}
		gas, err := gear_utils.GetMinimalGasForProgram(resp)
		if err != nil {
			logger.Log().Errorf("%v", err)
			os.Exit(1)
		}
		logger.Log().Printf("\n\n gas: %v\n\n", *gas)
		p := &extrinsic_params.GearProgram{
			CodeId:      codeId,
			Salt:        "0x11",
			InitPayload: payload,
			GasLimit:    *gas,
			Value:       "1",
			KeepAlive:   true,
		}

		createArgs := []any{p.CodeId, p.Salt, p.InitPayload, p.GasLimit, p.Value, p.KeepAlive}
		return createArgs, nil
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
