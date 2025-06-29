/*
	Use this Builder to make []ExtrinsicParams by your own or pick pre-build in the current directory

	Use `gear_utils.GetArgumentsForExtrinsicParams` function to get ExtrinsicParams arguments
	([]types.MetadataModuleCallArgument)

	e.g.

	//assuming we would like to get Gear/SendMessage args

	sendMessageArgs := gear_utils.GetArgumentsForExtrinsicParams("Gear", "send_message", modules)

	// modules is modules []types.MetadataModules is being gathered from types.MetadataStruct
	// example from gear_go.NewGear `gear.scale.GetMetadata().Metadata.Modules`

*/

package extrinsic_params

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
)

func InitBuilder(moduleName, callName string, modules []types.MetadataModules, values []any) ([]scalecodec.ExtrinsicParam, error) {
	args := gear_utils.GetArgumentsForExtrinsicParams(moduleName, callName, modules)
	if len(values) != len(args) {
		errMsg := fmt.Sprintf("invalid number of values for extrinsic params, expected %d, got %d", len(args), len(values))
		return nil, errors.New(errMsg)
	}
	var params []scalecodec.ExtrinsicParam
	for i, v := range args {
		params = append(params, scalecodec.ExtrinsicParam{
			Name:  v.Name,
			Type:  v.Type,
			Value: values[i],
		})
	}
	return params, nil
}
