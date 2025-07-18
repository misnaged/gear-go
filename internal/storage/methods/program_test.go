package gear_storage_methods

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestGearRpc() (gear_rpc.IGearRPC, *metadata.Metadata, error) {
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
