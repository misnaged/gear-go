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
	// TODO: Remove !!
	BobAccountId = "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"
	AliceSeed    = "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

func NewGear() (*Gear, error) {
	// Keeping subsequence of inits is must!
	gear := &Gear{
		config: initConfig(),
	}
	if err := gear.preRequests(); err != nil {
		return nil, fmt.Errorf(" gear.preRequests failed: %v", err)
	}
	if err := gear.initClient(); err != nil {
		return nil, fmt.Errorf(" gear.initClient failed: %v", err)
	}

	gear.initGearRpc()

	if err := gear.initScale(); err != nil {
		return nil, fmt.Errorf(" gear.initScale failed: %v", err)
	}
	//
	if err := gear.scale.InitMetadata(); err != nil {
		return nil, fmt.Errorf(" gear.scale.InitMetadata failed: %v", err)
	}
	/*
		//storage example call:

		storage := gear_storage_methods.NewStorage("System", "Account", gear.scale.GetMetadata())
		if err := storage.BuildParams(AliceSeed); err != nil {
			return nil, fmt.Errorf(" gear.buildParams failed: %v", err)
		}
		var toDecode models.FrameSystemAccountInfo
		err := storage.DecodeStorage(gear.gearRPC, &toDecode)
		if err != nil {
			return nil, fmt.Errorf(" gear.scale failed: %v", err)
		}
	*/

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
			return fmt.Errorf("ws.Handler failed: %v", err)
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
