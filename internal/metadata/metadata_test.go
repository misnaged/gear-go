package metadata

import (
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMetadata(t *testing.T) {
	cfg := &config.Scheme{}
	if err := config.InitConfig(cfg); err != nil {
		assert.NoError(t, err)
	}

	var client gear_client.IClient
	if cfg.Client.IsWebSocket {
		wsCli, err := gear_ws.NewWsClient(cfg)
		if err != nil {
			assert.NoError(t, err)
		}
		client = wsCli
	} else {

		httpCli := gear_http.NewHttpClient(time.Second*10, cfg)
		client = httpCli
	}
	gearGRPC := gear_rpc_method.NewGearRpc(client, cfg)
	_, err := NewMetadata(gearGRPC)
	assert.NoError(t, err)
}
