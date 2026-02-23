package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/calls"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const AliceSeed = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"

func NewGearTest() (*Gear, error) {
	gear := &Gear{}
	cfg := &config.Scheme{}
	if err := config.InitConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %v", err)
	}

	var client gear_client.IClient
	if cfg.Client.IsWebSocket {
		wsCli, err := gear_ws.NewWsClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("ws.Handler failed: %w", err)
		}
		client = wsCli
		gear.wsClient = wsCli
	} else {
		httpCli := gear_http.NewHttpClient(time.Second*10, cfg)
		client = httpCli
	}
	gear.client = client
	kr := keyring.New(keyring.Sr25519Type, AliceSeed)
	gear.keyRing = kr
	gear.gearRPC = gear_rpc_method.NewGearRpc(client, cfg)
	meta, err := metadata.NewMetadata(gear.gearRPC)
	if err != nil {
		return nil, fmt.Errorf("new metadata failed: %w", err)
	}
	gear.meta = meta
	gear.calls = calls.NewCalls(gear.meta, gear.gearRPC, gear.keyRing)
	return gear, nil
}
func TestNewGear(t *testing.T) {
	gear, err := NewGearTest()
	assert.NoError(t, err)
	assert.NotNil(t, gear)
}
