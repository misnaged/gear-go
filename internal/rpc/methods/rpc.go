package gear_rpc_method

import (
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
)

type GearRpc struct {
	client gear_client.IClient
	config *config.Scheme
}

func NewGearRpc(client gear_client.IClient, config *config.Scheme) gear_rpc.IGearRPC {
	return &GearRpc{
		client: client,
		config: config,
	}
}
