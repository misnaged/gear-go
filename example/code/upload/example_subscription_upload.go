package main

import (
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	gear_calls "github.com/misnaged/gear-go/internal/calls/gear"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"math/rand"
	"os"
	"time"
)

const wasmPath = "assets/wasm/test/demo_ping.opt.wasm"

func main() {
	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("error creating gear: %v", err)
		os.Exit(1)
	}
	var types = []string{"submitAndWatchExtrinsic1", "submitAndWatchExtrinsic2", "subscribeStorage"}

	err = gear.GetWsClient().AddResponseTypesAndMakeWsConnectionsPool(types...)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	calls := gear_calls.New(gear.GetCalls())
	codeId, err := gear_utils.GetCodeIdFromWasmFile(wasmPath)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	upload, err := calls.UploadCode(wasmPath)
	var uploadArgs, createArgs []string
	uploadArgs = append(uploadArgs, upload)
	owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
	payload := gear_utils.TextToHex("PING")
	resp, err := gear.GetRPC().GearCalculateInitCreateGas(owner, codeId, payload, 1, false)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100)
	gas, err := gear_utils.GetMinimalGasForProgram(resp)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	salt := fmt.Sprintf("0x%d", randomNumber)
	prog := &extrinsic_params.GearProgram{
		CodeId:      codeId,
		Salt:        salt,
		InitPayload: payload,
		GasLimit:    *gas,
		Value:       "1",
		KeepAlive:   true,
	}
	fmt.Println(salt)
	nonce, err := gear.GetRPC().SystemAccountNextIndex(gear.GetCalls().KeyRing.PublicKey())
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}
	logger.Log().Infof("nonce is: %d  \n", int(nonce.Result.(float64)))
	create, err := calls.CreateProgram(wasmPath, prog)
	createArgs = append(createArgs, create)
	event := gear.EventsSubscription()
	gear.MergeSubscriptionFunctions(event, gear.EnqueuedHandler(strToAny(uploadArgs), strToAny(createArgs)))
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
