package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/calls"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/substrate-api-rpc/keyring"

	// nolint:typecheck
	gear_client "github.com/misnaged/gear-go/internal/client"

	//nolint:typecheck
	gear_http "github.com/misnaged/gear-go/internal/client/http"

	//nolint:typecheck
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"

	// nolint:typecheck
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"

	// nolint:typecheck
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"

	"github.com/misnaged/scriptorium/versioner"
	"time"
)

type Gear struct {
	config  *config.Scheme
	version *version.Version
	client  gear_client.IClient
	gearRPC gear_rpc.IGearRPC
	meta    *metadata.Metadata
	calls   *calls.Calls
	keyRing keyring.IKeyRing
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

	// Client (http/ws) initialization
	if err := gear.initClient(); err != nil {
		return nil, fmt.Errorf(" gear.initClient failed: %w", err)
	}

	// Keyring initialization
	gear.initKeyRing()

	// RPC initialization
	gear.initGearRpc()

	// Metadata initialization
	if err := gear.initMetadata(); err != nil {
		return nil, fmt.Errorf(" gear.Metadata failed: %w", err)
	}

	// Calls initialization
	gear.initCalls()

	return gear, nil
}

// ************** API Builders ************ //

func (gear *Gear) initGearRpc() {
	gearRpc := gear_rpc_method.NewGearRpc(gear.client, gear.config)
	gear.gearRPC = gearRpc
}

func (gear *Gear) initMetadata() error {
	meta, err := metadata.NewMetadata(gear.gearRPC)
	if err != nil {
		return fmt.Errorf(" gear.initMetadata failed: %w", err)
	}
	gear.meta = meta
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
func (gear *Gear) initKeyRing() {
	kr := keyring.New(gear.config.Keyring.Category, gear.config.Keyring.Seed)
	gear.keyRing = kr
}

func (gear *Gear) initCalls() {
	cs := calls.NewCalls(gear.meta, gear.gearRPC, gear.keyRing)
	gear.calls = cs
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

func (gear *Gear) GetRPC() gear_rpc.IGearRPC {
	return gear.gearRPC
}
func (gear *Gear) GetMeta() *metadata.Metadata {
	return gear.meta
}
