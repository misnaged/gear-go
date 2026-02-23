package calls

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	gear_ws "github.com/misnaged/gear-go/internal/client/ws"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func api() (*Calls, error) {
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
	gearRpc := gear_rpc_method.NewGearRpc(client, cfg)

	meta, err := metadata.NewMetadata(gearRpc)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")
	return NewCalls(meta, gearRpc, kr), nil

}

func TestGearCalls_SignTransaction(t *testing.T) {
	apiT, err := api()
	assert.NoError(t, err)

	var args []any
	assert.NoError(t, err)
	args = append(args, "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", "10000000000000000000", "", true, 1000000)
	params, err := extrinsic_params.InitBuilder("GearVoucher", "issue", apiT.Meta.GetMetadata().Metadata.Modules, args)
	assert.NoError(t, err)

	//Alice

	aa, err := apiT.SignTransaction("GearVoucher", "issue", params)
	assert.NoError(t, err)
	if aa == "" {
		assert.FailNow(t, "sign transaction failed")
	}
}
