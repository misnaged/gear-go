package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/calls"
	gear_events "github.com/misnaged/gear-go/internal/events"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/gear-go/internal/models"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	"github.com/misnaged/gear-go/pkg/logger"
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
	config       *config.Scheme
	version      *version.Version
	client       gear_client.IClient
	wsClient     gear_client.IWsClient
	gearRPC      gear_rpc.IGearRPC
	events       gear_events.IEvent
	meta         *metadata.Metadata
	calls        *calls.Calls
	keyRing      keyring.IKeyRing
	stop         chan struct{}
	responsePool []<-chan *models.SubscriptionResponse
}

// NewGear creates fully functional gear-go API instance
func NewGear() (*Gear, error) {
	// Keeping subsequence of inits is must!
	gear := &Gear{
		config: initConfig(),
		stop:   make(chan struct{}),
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

	if gear.config.Client.IsWebSocket {
		gear.initEvents()

		if err := gear.addStorageEventsToResponsePool(); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if !gear.config.Subscriptions.HasCustomPoolRunner {
			gear.ResponsePoolRunner()
		}
	}
	return gear, nil
}

func (gear *Gear) initEvents() {
	gear.events = gear_events.NewGearEvents(gear.GetMeta().GetMetadata())
}

// ************** API Builders ************ //

// ResponsePoolRunner is the main function that handles
// and operates subscriptions responses.
//
// You can add new subscription by using Subscribe function of gear.wsClient as shown below:
//
//		 	newSubscription, err := gear.wsClient.Subscribe(args, methodName)
//				...
//				gear.responsePool = append(gear.responsePool, newSubscription) // the part in which new sub is added to pool
//	 --------------------------------------------------------------------------------------------------------------------
//
// This function could be overwritten
// if config value of Subscriptions.HasCustomPoolRunner is set by true.
// In this case, you would have to make CustomRunner function by your own
// in the way it suits you
func (gear *Gear) ResponsePoolRunner() {
	for _, resp := range gear.responsePool {
		select {
		case <-gear.stop:
			return
		case <-resp:
			for e := range resp {
				if e.Params != nil {
					if err := gear.Get(e); err != nil {
						logger.Log().Errorf("gear.Get failed: %v", err)
						return
					}
				}
			}
		}
	}
}

func (gear *Gear) Get(resp *models.SubscriptionResponse) error {
	changes, err := models.GetChangesFromEvents(resp)
	if err != nil {
		return fmt.Errorf("gear.responsePoolRunner - GetChangesFromEvents failed: %w", err)
	}
	if changes[1] != nil {
		if changes[1].(string) != "" {
			events, err := gear.events.GetEvents(changes[1].(string))
			if err != nil {
				return fmt.Errorf("gear.responsePoolRunner - GetEvents failed: %w", err)
			}
			err = gear.events.Handle(events)
			if err != nil {
				return fmt.Errorf("gear.responsePoolRunner - HandleEvents failed: %w", err)
			}
		}
	}
	return nil
}
func (gear *Gear) addStorageEventsToResponsePool() error {
	storage := gear_storage_methods.NewStorage("System", "Events", gear.GetMeta(), gear.GetRPC())
	k, err := storage.GetStorageKey()
	if err != nil {
		logger.Log().Errorf(" storage.GetStorageKeys failed: %v", err)
	}

	storageSub, err := gear.wsClient.Subscribe([][]string{{k}}, "state_subscribeStorage")
	if err != nil {
		return fmt.Errorf(" gear.wsClient.Subscribe2 failed: %w", err)
	}

	gear.responsePool = append(gear.responsePool, storageSub)

	return nil
}

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
		gear.wsClient = client
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

func (gear *Gear) GetWsClient() gear_client.IWsClient {
	if !gear.config.Client.IsWebSocket {
		logger.Log().Error("config setting for gear.client.IsWebSocket is false")
		return nil
	}
	if gear.wsClient == nil {
		logger.Log().Error("ws client not initialized")
		return nil
	}
	return gear.wsClient
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
