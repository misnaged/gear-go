package gear_storage_methods

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestGearRpc() (gear_rpc.IGearRPC, *metadata.Metadata, error) {
	cfg := &config.Scheme{}
	if err := config.InitConfig(cfg); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize config: %v", err)
	}

	var client gear_client.IClient
	if cfg.Client.IsWebSocket {
		wsCli, err := gear_ws.NewWsClient(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("ws.Handler failed: %w", err)
		}
		client = wsCli
	} else {

		httpCli := gear_http.NewHttpClient(time.Second*10, cfg)
		client = httpCli
	}
	gearGRPC := gear_rpc_method.NewGearRpc(client, cfg)
	meta, err := metadata.NewMetadata(gearGRPC)
	if err != nil {
		return nil, nil, fmt.Errorf("metadata.NewMetadata failed: %v", err)
	}

	return gearGRPC, meta, nil
}
func TestStorage_GetProgramsId(t *testing.T) {
	rpc, meta, err := newTestGearRpc()
	assert.NoError(t, err)

	storage := NewStorage("GearProgram", "ProgramStorage", meta, rpc)

	_, err = storage.GetProgramsId()
	assert.NoError(t, err)
}
