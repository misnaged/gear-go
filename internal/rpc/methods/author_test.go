package gear_rpc_method

import (
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGearRpc_AuthorSubmitExtrinsic(t *testing.T) {
	gearRpc, api, err := newTestGearRpc()
	assert.NoError(t, err)
	err = api.InitMetadata()
	assert.NoError(t, err)
	f, err := os.ReadFile(testWasmPing)
	assert.NoError(t, err)

	var args []any
	toHex := gear_utils.AddToHex(f)
	args = append(args, toHex)
	params, err := extrinsic_params.InitBuilder("Gear", "upload_code", api.GetMetadata().Metadata.Modules, args)
	assert.NoError(t, err)

	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	signed, err := api.SignTransaction("Gear", "upload_code", kr, params)
	assert.NoError(t, err)

	_, err = gearRpc.AuthorSubmitExtrinsic(signed)
	assert.NoError(t, err)
}
