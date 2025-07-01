package gear_scale

import (
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
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

func scale() (*Scale, *types.MetadataStruct, error) {
	cfg := &config.Scheme{}
	client := gear_http.NewHttpClient(time.Second*10, cfg)
	gearRpc := gear_rpc_method.NewGearRpc(client, cfg)

	scl := NewScale(gearRpc, cfg)
	decoder := &scalecodec.MetadataDecoder{}
	resp, err := scl.gearRpc.StateGetMetadataLatest()
	if err != nil {
		return nil, nil, fmt.Errorf("post request err: %v", err)
	}
	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	return scl, &decoder.Metadata, nil
}

func TestScale_SignTransaction(t *testing.T) {

	scaleT, metaT, err := scale()
	assert.NoError(t, err)
	fmt.Printf("%+v\n", scaleT)
	var args []any
	//args = append(args, toHex)
	params, err := extrinsic_params.InitBuilder("GearDebug", "upload_code", metaT.Metadata.Modules, args)
	assert.NoError(t, err)
	if params == nil {
	}
	//Alice
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")
	kr.Sign("")
}
