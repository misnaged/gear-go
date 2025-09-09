package gear_rpc_method

import (
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testWasmPing = "../../../assets/wasm/test/demo_ping.opt.wasm"

func TestGearRpc_GearCalculateInitCreateGas(t *testing.T) {
	gearRpc, err := newTestGearRpc()
	assert.NoError(t, err)
	codeId, err := gear_utils.GetCodeIdFromWasmFile(testWasmPing)
	assert.NoError(t, err)
	pingPayload := "0x50494e47"
	Alice := "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
	_, err = gearRpc.GearCalculateInitCreateGas(codeId, Alice, pingPayload, 1, true)
	assert.NoError(t, err)
}
