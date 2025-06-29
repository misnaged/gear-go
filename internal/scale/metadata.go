package gear_scale

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
)

func (s *Scale) InitMetadata() error {
	decoder := &scalecodec.MetadataDecoder{}
	resp, err := s.gearRpc.StateGetMetadataLatest()
	if err != nil {
		return fmt.Errorf("post request err: %v", err)
	}
	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return fmt.Errorf("failed to decode metadata: %w", err)
	}
	s.metadata = &decoder.Metadata
	return nil
}

func (s *Scale) GetMetadata() *types.MetadataStruct {
	return s.metadata
}

func (s *Scale) MetadataCheck() error {
	if s.metadata != nil {
		return nil
	}
	return errors.New("metadata is nil")
}
