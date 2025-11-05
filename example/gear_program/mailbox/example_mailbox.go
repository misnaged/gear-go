package main

import (
	"errors"
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
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
	var types = []string{"upload_code", "create_program", "send_message", "subscribeStorage"}
	uplCode := gear_go.NewEnq("upload_code", "upload_code", "Gear", "author_submitAndWatchExtrinsic", true)
	createPgm := gear_go.NewEnq("create_program", "create_program", "Gear", "author_submitAndWatchExtrinsic", true)
	sendMsg := gear_go.NewEnq("send_message", "send_message", "Gear", "author_submitAndWatchExtrinsic", true)
	subscribeStorage := gear_go.NewEnq("subscribeStorage", "", "", "", false)
	subscribeStorage.CustomFunc = func() error {
		fmt.Println("Custom function called for subscribeStorage")
		return nil
	}

	sendMsg.AfterFinalizationFunc = func() error {
		storage := gear_storage_methods.NewStorage("GearMessenger", "Mailbox", gear.GetMeta(), gear.GetRPC())

		k, _ := storage.GetStorageKeys()
		storageDataArr, err := storage.DecodeStorageDataMap(k[0])
		if err != nil {
			return fmt.Errorf(" gear.DecodeStorageDataArray failed: %w", err)
		}

		//nolint:staticcheck
		mbx := &models.Mailbox{}

		for _, data := range storageDataArr {
			m, ok := data.(map[string]any)
			if !ok {
				return errors.New(" gear.DecodeStorageDataArray failed ")
			}

			if v, ok := m["id"].(string); ok {
				mbx.Id = v
			}
			if v, ok := m["payload"].(string); ok {
				mbx.Payload = v
			}
			if v, ok := m["source"].(string); ok {
				mbx.Source = v
			}
			if v, ok := m["value"].(string); ok {
				mbx.Value = v
			}
			if v, ok := m["start"].(float64); ok {
				mbx.Start = v
			}
			if v, ok := m["finish"].(float64); ok {
				mbx.Finish = v
			}
		}

		//nolint:staticcheck
		if mbx == nil {
			return errors.New(" gear.DecodeStorageDataArray failed ")
		}
		fmt.Printf(" %#v \n", mbx)
		return nil
	}

	// general preparation
	err = gear.GetWsClient().AddResponseTypesAndMakeWsConnectionsPool(types...)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	var dd []*gear_go.EnquedSubscription

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

	uplCode.Args = strToAny(uploadArgs)
	inter := gear_go.NewInterruption("create_program", CalculateGas(gear, codeId))
	inter2 := gear_go.NewInterruption("send_message", SendMessageArgs(gear, codeId))

	dd = append(dd, uplCode, createPgm, sendMsg, subscribeStorage)
	gear.MergeSubscriptionFunctions(gear.EventsSubscription(), gear.EnqueuedSubscriptionsOptional(dd, inter, inter2))
	err = gear.InitSubscriptions()
	if err != nil {
		logger.Log().Errorf(" gear.ProcessEventsSubscription failed: %v", err)
		os.Exit(1)
	}
}

func SendMessageArgs(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {

		storage := gear_storage_methods.NewStorage("GearProgram", "ProgramStorage", gear.GetMeta(), gear.GetRPC())

		program, err := storage.GetActiveProgramByCodeId(codeId)
		if err != nil {
			return nil, fmt.Errorf("gear.GetActiveProgramByCodeId failed: %w", err)
		}
		owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice

		value := "100000000000000"

		payload := "0x04" // todo;

		resp, err := gear.GetRPC().GearCalculateHandleGas(owner, program.ProgramId, payload, value, true)
		if err != nil {
			return nil, fmt.Errorf("gear.GetRPC().GearCalculateHandleGas failed: %w", err)
		}
		gas, err := gear_utils.GetMinimalGasForProgram(resp)
		if err != nil {
			return nil, fmt.Errorf("GetMinimalGasForProgram failed: %w", err)

		}

		sendMessageParams := &extrinsic_params.GearSendMessage{
			ProgramId: program.ProgramId,
			Payload:   payload,
			GasLimit:  *gas,
			Value:     value,
			KeepAlive: true,
		}
		createArgs := []any{
			sendMessageParams.ProgramId,
			sendMessageParams.Payload,
			sendMessageParams.GasLimit,
			sendMessageParams.Value,
			sendMessageParams.KeepAlive,
		}

		return createArgs, nil
	}
}

func CalculateGas(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {
		owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
		payload := "0x0"
		resp, err := gear.GetRPC().GearCalculateInitCreateGas(owner, codeId, payload, 0, true)
		if err != nil {
			return nil, fmt.Errorf("gear.GetRPC().GearCalculateInitCreateGas failed: %w", err)
		}
		gas, err := gear_utils.GetMinimalGasForProgram(resp)
		if err != nil {
			return nil, fmt.Errorf("GetMinimalGasForProgram failed: %w", err)
		}
		logger.Log().Printf("\n\n gas: %v\n\n", *gas)
		p := &extrinsic_params.GearProgram{
			CodeId:      codeId,
			Salt:        "0x1",
			InitPayload: payload,
			GasLimit:    *gas,
			Value:       "0",
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
