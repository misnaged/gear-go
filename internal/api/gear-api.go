package gear_api

import (
	"github.com/misnaged/gear-go/config"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"github.com/misnaged/substrate-api-rpc/rpc"

	"github.com/itering/scale.go/types"
)

const (
	VaraPrefix      uint16 = 137
	SubstratePrefix uint16 = 42
)

type Api struct {
	gearRpc  gear_rpc.IGearRPC
	metadata *types.MetadataStruct
	config   *config.Scheme
	customTx rpc.ICustomTranscation
}

func NewApi(gearRpc gear_rpc.IGearRPC, config *config.Scheme) *Api {
	return &Api{
		gearRpc: gearRpc,
		config:  config,
	}
}
