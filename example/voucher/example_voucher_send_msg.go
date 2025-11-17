package main

import (
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"os"
)

const wasmPath = "assets/wasm/test/demo_ping.opt.wasm"
const CharliePublicKey = "90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22"
const CharlieSeed = "0xbc1ede780f784bb6991a585e4f6e61522c14e1cae6ad0895fb57b9a205a8f938"

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}
	var types = []string{"upload_code", "create_program", "voucher_issue", "voucher_call", "subscribeStorage"}
	uplCode := gear_go.NewEnq("upload_code", "upload_code", "Gear", "author_submitAndWatchExtrinsic", true)
	createPgm := gear_go.NewEnq("create_program", "create_program", "Gear", "author_submitAndWatchExtrinsic", true)
	voucherIssue := gear_go.NewEnq("voucher_issue", "issue", "GearVoucher", "author_submitAndWatchExtrinsic", true)
	voucherCall := gear_go.NewEnq("voucher_call", "call", "GearVoucher", "author_submitAndWatchExtrinsic", true)
	voucherCall.CustomBuilderKeyRing = keyring.New(keyring.Sr25519Type, CharlieSeed)
	voucherCall.IsCustomBuilder = true
	subscribeStorage := gear_go.NewEnq("subscribeStorage", "", "", "", false)
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
	fmt.Println("codeId:", codeId)
	upload, err := extrinsic_params.NewGearCode(wasmPath)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	uploadArgs := []string{upload.Code}

	uplCode.Args = gear_utils.StrToAny(uploadArgs)
	inter := gear_go.NewInterruption("create_program", CalculateGas(gear, codeId))
	inter2 := gear_go.NewInterruption("voucher_issue", VoucherIssueArgs(gear, codeId))
	inter3 := gear_go.NewInterruption("voucher_call", VoucherSendMsgArgs(gear, codeId))

	dd = append(dd, uplCode, createPgm, voucherIssue, voucherCall, subscribeStorage)
	gear.MergeSubscriptionFunctions(gear.EventsSubscription(), gear.EnqueuedSubscriptionsOptional(dd, inter, inter2, inter3))
	err = gear.InitSubscriptions()
	if err != nil {
		logger.Log().Errorf(" gear.ProcessEventsSubscription failed: %v", err)
		os.Exit(1)
	}
}

func VoucherIssueArgs(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {
		progStor := gear_storage_methods.NewStorage("GearProgram", "ProgramStorage", gear.GetMeta(), gear.GetRPC())
		program, err := progStor.GetActiveProgramByCodeId(codeId)
		if err != nil {
			return nil, fmt.Errorf(" storage.GetActiveProgramByCodeId failed: %v", err)
		}
		charlie := fmt.Sprintf("%s%s", "0x", CharliePublicKey)
		return []any{charlie, "100000000000000", []any{program.ProgramId}, false, 100000}, nil
	}
}
func VoucherSendMsgArgs(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {
		storage := gear_storage_methods.NewStorage("GearVoucher", "Vouchers", gear.GetMeta(), gear.GetRPC())
		err := storage.AddAccountIdToStorageParams(CharliePublicKey)
		if err != nil {
			return nil, fmt.Errorf(" storage.AddAccountIdToStorageParams failed: %w", err)
		}

		storkey, err := storage.GetVoucherStorageKeys()
		if err != nil {
			return nil, fmt.Errorf(" storage.GetStorageKeys failed: %w", err)
		}
		progStor := gear_storage_methods.NewStorage("GearProgram", "ProgramStorage", gear.GetMeta(), gear.GetRPC())
		program, err := progStor.GetActiveProgramByCodeId(codeId)
		if err != nil {
			return nil, fmt.Errorf(" storage.GetActiveProgramByCodeId failed: %v", err)
		}
		callSendMsg := extrinsic_params.NewVoucherCallSendMessage(program.ProgramId, "0x50494e47", "0", 20_000_000_000, false)
		voucherId := storkey[0][128+2:] // +2 bytes for `0x`

		return []any{fmt.Sprintf("0x%s", voucherId), callSendMsg}, nil
	}
}

func CalculateGas(gear *gear_go.Gear, codeId string) func() ([]any, error) {
	return func() ([]any, error) {
		owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
		payload := "0x50494e47"
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
