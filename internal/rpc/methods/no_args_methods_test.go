package gear_rpc_method

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestGearRpc() (gear_rpc.IGearRPC, error) {
	clientCfg := &config.Client{
		IsWebSocket: false,
		IsSecured:   false,
	}
	clientCfg.Transport = "http"
	clientCfg.Host = "127.0.0.1"
	clientCfg.Port = 9944
	cfg := &config.Scheme{Client: clientCfg}

	fmt.Println(cfg.Client.Transport)

	client := gear_http.NewHttpClient(time.Second*3, cfg)
	return NewGearRpc(client, cfg), nil
}
func TestGearRpc_NoArgRpcRequest(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)

	for _, v := range gear_rpc.NoArgsMethods {
		fmt.Println("calling", gear_rpc.NoArgMethodFromString(v))
		_, err = gearRpc.NoArgRpcRequest(gear_rpc.NoArgMethodFromString(v))
		assert.NoError(t, err)
	}

}
