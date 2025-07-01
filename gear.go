package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/client/ws"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	gear_scale "github.com/misnaged/gear-go/internal/scale"
	"github.com/misnaged/scriptorium/versioner"
	"time"
)

type Gear struct {
	config  *config.Scheme
	version *version.Version
	client  gear_client.IClient
	scale   *gear_scale.Scale
	gearRPC gear_rpc.IGearRPC
}

const (
	BobAccountId = "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"
	AliceSeed    = "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
	AliceSecret  = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"
)

// TODO: for debug only! Remove it in the next update!

func NewGear() (*Gear, error) {
	// Keeping subsequence of inits is must!
	gear := &Gear{
		config: initConfig(),
	}
	if err := gear.preRequests(); err != nil {
		return nil, fmt.Errorf(" gear.preRequests failed: %w", err)
	}
	if err := gear.initClient(); err != nil {
		return nil, fmt.Errorf(" gear.initClient failed: %w", err)
	}

	gear.initGearRpc()

	if err := gear.initScale(); err != nil {
		return nil, fmt.Errorf(" gear.initScale failed: %w", err)
	}
	//
	if err := gear.scale.InitMetadata(); err != nil {
		return nil, fmt.Errorf(" gear.scale.InitMetadata failed: %w", err)
	}
	// TODO: add all examples to Readme

	/*
		// TODO: REMOVE (Debug)
		dec := gear_utils.GetExtrinsicDecoderByRawHex(kkk, gear.scale.GetMetadata())
		dec2 := gear_utils.GetExtrinsicDecoderByRawHex(kkk2, gear.scale.GetMetadata())
		var val1, val2 string
		for _, k := range dec.Params {
			if k.Name == "code" {
				val1 = k.Value.(string)
			}
		}
		for _, k := range dec2.Params {
			if k.Name == "code" {
				val2 = k.Value.(string)
			}
		}

	*/
	//option := types.ScaleDecoderOption{Metadata: gear.scale.GetMetadata()}

	//logs := digest["logs"].([]any)
	//m := types.ScaleDecoder{}

	//m.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes("0x05424142450101786b36574e09b48fa1988b5b506a6a4024aeedf65fb415e5ecbe000bc66d301f983902d254fc0764ac2541627250cdeaee15b6d6f7b3cd1e0d8b66221ebc6d83")}, nil)
	//r := m.ProcessAndUpdateData("DigestItem")
	//dec := &types.Vec{}
	//dec.Process()

	//TODO: CLEAN!!

	//-----------------------  Upload Code example ----------------- //

	//a, err := gear.UploadCodeTemp()
	//if err != nil {
	//	return nil, fmt.Errorf(" gear.UploadCodeTemp() failed: %w", err)
	//}
	//
	//var args []string
	//args = append(args, a)
	//gear.client.Subscribe(args, "author_submitAndWatchExtrinsic")

	//storage := gear_storage_methods.NewStorage("GearProgram", "CodeStorage", gear.GetScale().GetMetadata())
	//var vv map[string]any
	//err := storage.DecodeStorage(gear.GetRPC(), &vv, true)
	//if err != nil {
	//	return nil, fmt.Errorf(" gear_rpc.DecodeStorage failed: %w", err)
	//}

	// -----------------------------------------------------------------//
	//str, _ := gear_utils.GetCodeIdFromWasmFile("demo_messenger.opt.wasm")
	//storage example call:
	/*
		storage := gear_storage_methods.NewStorage("GearProgram", "CodeStorage", gear.scale.GetMetadata())

		var vv map[string]any
		err := storage.DecodeStorage(gear.gearRPC, &vv, true)
		if err != nil {
			return nil, fmt.Errorf(" gear.scale failed: %w", err)
		}
	*/
	//code := vv["code"] //large data to write!

	return gear, nil
}

func (gear *Gear) initGearRpc() {
	gearRpc := gear_rpc_method.NewGearRpc(gear.client, gear.config)
	gear.gearRPC = gearRpc
}
func (gear *Gear) initScale() error {
	scale := gear_scale.NewScale(gear.gearRPC, gear.config)
	gear.scale = scale
	return nil
}

func (gear *Gear) initClient() error {
	if gear.config.Client.IsWebSocket {
		client, err := gear_ws.NewWsClient(gear.config)
		if err != nil {
			return fmt.Errorf("ws.Handler failed: %w", err)
		}
		gear.client = client
	} else {
		client := gear_http.NewHttpClient(time.Second*10, gear.config)
		gear.client = client
	}
	return nil
}

func (gear *Gear) preRequests() error {

	vers, err := initVersion()
	if err != nil {
		return fmt.Errorf("initialize version: %w", err)
	}
	gear.version = vers
	if err = config.InitConfig(gear.config); err != nil {
		return fmt.Errorf("failed initialize config: %w", err)
	}
	return nil
}

func initConfig() *config.Scheme {
	return &config.Scheme{}
}

func initVersion() (*version.Version, error) {
	ver, err := version.NewVersion()
	if err != nil {
		return nil, fmt.Errorf("init app version: %w", err)
	}
	return ver, nil
}

func (gear *Gear) GetConfig() *config.Scheme {
	return gear.config
}

func (gear *Gear) GetClient() gear_client.IClient {
	return gear.client
}
func (gear *Gear) GetScale() *gear_scale.Scale {
	return gear.scale
}
func (gear *Gear) GetRPC() gear_rpc.IGearRPC {
	return gear.gearRPC
}
