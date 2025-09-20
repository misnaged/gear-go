package metadata

import (
	"github.com/misnaged/gear-go/config"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMetadata(t *testing.T) {
	clientCfg := &config.Client{
		IsWebSocket: false,
		IsSecured:   false,
	}
	clientCfg.Transport = "http"
	clientCfg.Host = "127.0.0.1"
	clientCfg.Port = 9944
	cfg := &config.Scheme{Client: clientCfg}

	client := gear_http.NewHttpClient(time.Second*3, cfg)
	gearGRPC := gear_rpc_method.NewGearRpc(client, cfg)
	_, err := NewMetadata(gearGRPC)
	assert.NoError(t, err)
}
