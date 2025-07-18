package calls

import (
	"github.com/misnaged/gear-go/internal/metadata"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"github.com/misnaged/substrate-api-rpc/rpc"
)

type IGearCalls interface {
}

type GearCalls struct {
	GearRpc  gear_rpc.IGearRPC
	Meta     *metadata.Metadata
	customTx rpc.ICustomTranscation
}

func NewGearCalls(meta *metadata.Metadata, gearRpc gear_rpc.IGearRPC) *GearCalls {
	return &GearCalls{
		Meta:    meta,
		GearRpc: gearRpc,
	}
}
