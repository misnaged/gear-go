package gear_events

import (
	"github.com/misnaged/gear-go/config"
	gear_http "github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc_method "github.com/misnaged/gear-go/internal/rpc/methods"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testEventHash = `0x1c00000000000000220f4040551702000000010000000508d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d10b754aa34000000000000000000000000000100000005076d6f646c70792f7472737279000000000000000000000000000000000000000010b754aa3400000000000000000000000000010000000e0410b754aa3400000000000000000000000000010000000600d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d10b754aa340000000000000000000000000000000000000000000000000000000000010000000001036806000000038e62d286e93700000000020000000000623bdc3700020100`

func newTestGearRpc() (*metadata.Metadata, error) {

	clientCfg := &config.Client{
		IsWebSocket: false,
		IsSecured:   false,
	}
	clientCfg.Transport = "http"
	clientCfg.Host = "127.0.0.1"
	clientCfg.Port = 9944
	cfg := &config.Scheme{Client: clientCfg}

	client := gear_http.NewHttpClient(time.Second*3, cfg)
	gearGRPC := gear_rpc_method.NewGearRpc(client, cfg)
	meta, err := metadata.NewMetadata(gearGRPC)
	if err != nil {
		return nil, err
	}
	return meta, nil
}
func TestGearEvent_GetEvents(t *testing.T) {
	meta, err := newTestGearRpc()
	assert.NoError(t, err)

	event := NewGearEvents(meta.GetMetadata())
	events, err := event.GetEvents(testEventHash)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
}

func TestGearEvent_Handle(t *testing.T) {
	meta, err := newTestGearRpc()
	assert.NoError(t, err)

	event := NewGearEvents(meta.GetMetadata())
	events, err := event.GetEvents(testEventHash)
	assert.NoError(t, err)
	err = event.Handle(events)
	assert.NoError(t, err)
}
