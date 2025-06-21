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
	cfg := &config.Scheme{}
	err := config.InitConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize config: %v", err)
	}
	client := gear_http.NewHttpClient(time.Second*10, cfg)
	return NewGearRpc(client, cfg), nil
}
func TestGearRpc_NoArgRpcRequest(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)

	for _, v := range gear_rpc.NoArgsMethods {
		_, err = gearRpc.NoArgRpcRequest(gear_rpc.NoArgMethodFromString(v))
		assert.NoError(t, err)
	}

}
