package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/calls"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
)

const AliceSeed = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"

func NewGearTest() (*Gear, error) {
	gear := &Gear{}
	clientCfg := &config.Client{
		IsWebSocket: true,
		IsSecured:   false,
	}
	clientCfg.Transport = "ws"
	clientCfg.Host = "127.0.0.1"
	clientCfg.Port = 9944
	cfg := &config.Scheme{Client: clientCfg}
	gear.config = cfg
	client, err := gear_ws.NewWsClient(gear.config)
	if err != nil {
		return nil, fmt.Errorf("ws.Handler failed: %w", err)
	}
	gear.client = client
	gear.wsClient = client
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
