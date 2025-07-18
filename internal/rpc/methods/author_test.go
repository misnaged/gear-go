package gear_rpc_method

import (
	"github.com/misnaged/gear-go/internal/calls"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/gear-go/internal/models/extrinsic_params"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGearRpc_AuthorSubmitExtrinsic(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)
	meta, err := metadata.NewMetadata(gearRpc)
	assert.NoError(t, err)
	f, err := os.ReadFile(testWasmPing)
	assert.NoError(t, err)
	call := calls.NewGearCalls(meta, gearRpc)

	var args []any
	toHex := gear_utils.AddToHex(f)
	args = append(args, toHex)
	params, err := extrinsic_params.InitBuilder("Gear", "upload_code", meta.GetMetadata().Metadata.Modules, args)
	assert.NoError(t, err)

	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	signed, err := call.SignTransaction("Gear", "upload_code", kr, params)
	assert.NoError(t, err)

	_, err = gearRpc.AuthorSubmitExtrinsic(signed)
	assert.NoError(t, err)
}
