package gear_calls

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/calls"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

const testWasmPath = "assets/wasm/test/demo_ping.opt.wasm"

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
	gearGRPC := gear_rpc_method.NewGearRpc(client, cfg)

	return gearGRPC, nil
}

// Warning: for correct work code must be uploaded to chain through upload_code extrinsic
func TestGearCalls_CreateProgram(t *testing.T) {
	err := os.Chdir("../../../")
	assert.NoError(t, err)
	codeId, err := gear_utils.GetCodeIdFromWasmFile(testWasmPath)
	assert.NoError(t, err)
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)

	meta, err := metadata.NewMetadata(gearRpc)
	assert.NoError(t, err)
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	c := calls.NewCalls(meta, gearRpc, kr)
	gearCall := New(c)

	owner := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" //Alice
	payload := gear_utils.TextToHex("PING")
	resp, err := gearRpc.GearCalculateInitCreateGas(owner, codeId, payload, 1, true)
	assert.NoError(t, err)
	gas, err := gear_utils.GetMinimalGasForProgram(resp)
	assert.NoError(t, err)

	p := &extrinsic_params.GearProgram{
		CodeId:      codeId,
		Salt:        "0x1",
		InitPayload: payload,
		GasLimit:    *gas,
		Value:       "1",
		KeepAlive:   true,
	}
	str, err := gearCall.CreateProgram(testWasmPath, p)
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
}
