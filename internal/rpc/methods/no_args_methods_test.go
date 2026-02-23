package gear_rpc_method

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestGearRpc() (gear_rpc.IGearRPC, error) {

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
	} else {

		httpCli := gear_http.NewHttpClient(time.Second*10, cfg)
		client = httpCli
	}
	gearGRPC := NewGearRpc(client, cfg)

	return gearGRPC, nil
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
