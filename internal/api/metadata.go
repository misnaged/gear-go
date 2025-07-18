package gear_api

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
)

func (api *Api) InitMetadata() error {
	decoder := &scalecodec.MetadataDecoder{}
	resp, err := api.gearRpc.StateGetMetadataLatest()
	if err != nil {
		return fmt.Errorf("post request err: %v", err)
	}
	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return fmt.Errorf("failed to decode metadata: %w", err)
	}
	api.metadata = &decoder.Metadata
	return nil
}

func (api *Api) GetMetadata() *types.MetadataStruct {
	return api.metadata
}

func (api *Api) MetadataCheck() error {
	if api.metadata != nil {
		return nil
	}
	return errors.New("metadata is nil")
}
