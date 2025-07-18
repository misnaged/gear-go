package metadata

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
)

type Metadata struct {
	gearRpc gear_rpc.IGearRPC
	self    *types.MetadataStruct
}

func NewMetadata(rpc gear_rpc.IGearRPC) (*Metadata, error) {
	meta := &Metadata{gearRpc: rpc}
	self, err := meta.initMetadata()
	if err != nil {
		return nil, fmt.Errorf("metadata init error: %v", err)
	}
	meta.self = self
	return meta, nil
}
func (meta *Metadata) initMetadata() (*types.MetadataStruct, error) {
	decoder := &scalecodec.MetadataDecoder{}
	resp, err := meta.gearRpc.StateGetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("post request err: %v", err)
	}
	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	return &decoder.Metadata, nil
}

func (meta *Metadata) GetMetadata() *types.MetadataStruct {
	return meta.self
}

func (meta *Metadata) MetadataCheck() error {
	if meta.self != nil {
		return nil
	}
	return errors.New("metadata is nil")
}
