package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"

	// nolint:typecheck
	gear_client "github.com/misnaged/gear-go/internal/client"

	"github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/client/ws"

	// nolint:typecheck
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"

	// nolint:typecheck
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"

	// nolint:typecheck
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

// NewGear creates fully functional gear-go API instance
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

	return gear, nil
}

// ************** API Builders ************ //

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

// ************** Helpers **************** //

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
