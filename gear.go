package gear_go

import (
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/client/http"
	"github.com/misnaged/gear-go/internal/client/ws"

	gear_scale "github.com/misnaged/gear-go/internal/scale"
	"github.com/misnaged/scriptorium/versioner"
	"time"
)

type Gear struct {
	config  *config.Scheme
	version *version.Version
	client  gear_client.IClient
	scale   *gear_scale.Scale
}

func NewGear() (*Gear, error) {
	gear := &Gear{
		config: initConfig(),
	}
	if err := gear.preRequests(); err != nil {
		return nil, fmt.Errorf(" gear.preRequests failed: %v", err)
	}
	if err := gear.initClient(); err != nil {
		return nil, fmt.Errorf(" gear.initClient failed: %v", err)
	}
	if err := gear.initScale(); err != nil {
		return nil, fmt.Errorf(" gear.initScale failed: %v", err)
	}

	return gear, nil
}

func (gear *Gear) initScale() error {
	scale := gear_scale.NewScale(gear.client)
	gear.scale = scale
	return nil
}

func (gear *Gear) initClient() error {
	if gear.config.Client.IsWebSocket {
		client, err := gear_ws.NewWsClient(gear.config)
		if err != nil {
			return fmt.Errorf("ws.Handler failed: %v", err)
		}
		gear.client = client
	} else {
		client := gear_http.NewHttpClient(time.Second*10, gear.config)
		gear.client = client
	}

	return nil
}

func (gear *Gear) preRequests() error {

	vers, err := initVersion()
	if err != nil {
		return fmt.Errorf("initialize version: %w", err)
	}
	gear.version = vers
	if err = config.InitConfig(gear.config); err != nil {
		return fmt.Errorf("failed initialize config: %w", err)
	}
	return nil
}

func initConfig() *config.Scheme {
	return &config.Scheme{}
}

func initVersion() (*version.Version, error) {
	ver, err := version.NewVersion()
	if err != nil {
		return nil, fmt.Errorf("init app version: %w", err)
	}
	return ver, nil
}

func (gear *Gear) GetConfig() *config.Scheme {
	return gear.config
}

func (gear *Gear) GetClient() gear_client.IClient {
	return gear.client
}
func (gear *Gear) GetScale() *gear_scale.Scale {
	return gear.scale
}
