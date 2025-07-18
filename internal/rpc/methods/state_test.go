package gear_rpc_method

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGearRpc_GetMetadataLatest(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)

	_, err = gearRpc.StateGetMetadataLatest()
	assert.NoError(t, err)
}

func TestGearRpc_GetRuntimeVersionLatest(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)

	_, err = gearRpc.StateGetRuntimeVersionLatest()
	assert.NoError(t, err)
}
