package gear_calls

import (
	"github.com/misnaged/gear-go/internal/calls"
	"github.com/misnaged/gear-go/internal/metadata"
	"github.com/misnaged/substrate-api-rpc/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGearCalls_UploadCode(t *testing.T) {
	gearRpc := newTestGearRpc()
	meta, err := metadata.NewMetadata(gearRpc)
	assert.NoError(t, err)
	kr := keyring.New(keyring.Sr25519Type, "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a")

	c := calls.NewCalls(meta, gearRpc, kr)
	gearCall := New(c)
	str, err := gearCall.UploadCode(testWasmPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
}
