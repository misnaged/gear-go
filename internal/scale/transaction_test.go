package gear_scale

import (
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/utiles"
	"github.com/misnaged/gear-go/config"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func scale() (*Scale, error) {
	clientCfg := &config.Client{
		IsWebSocket: false,
		IsSecured:   false,
	}
	clientCfg.Transport = "http"
	clientCfg.Host = "127.0.0.1"
	clientCfg.Port = 9944
	cfg := &config.Scheme{Client: clientCfg}
	client := gear_http.NewHttpClient(time.Second*10, cfg)
	gearRpc := gear_rpc_method.NewGearRpc(client, cfg)

	scl := NewScale(gearRpc, cfg)
	decoder := &scalecodec.MetadataDecoder{}
	resp, err := scl.gearRpc.StateGetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("post request err: %v", err)
	}
	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	return scl, nil
}

func TestScale_SignTransaction(t *testing.T) {
	scaleT, err := scale()
	assert.NoError(t, err)
	err = scaleT.InitMetadata()
	assert.NoError(t, err)
	var args []any
	assert.NoError(t, err)
	args = append(args, "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", "10000000000000000000", "", true, "1000000")
	params, err := extrinsic_params.InitBuilder("GearVoucher", "issue", scaleT.GetMetadata().Metadata.Modules, args)
	assert.NoError(t, err)
	//Alice
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	aa, err := scaleT.SignTransaction("GearVoucher", "issue", kr, params)
	assert.NoError(t, err)
	if aa == "" {
		assert.FailNow(t, "sign transaction failed")
	}
}
