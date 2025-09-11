// events for 'System' module

package gear_events

import (
	"errors"
	"fmt"
	"github.com/itering/scale.go/types"
	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/pkg/logger"
)

func handleExtrinsicFailed(event *models.Event, meta *types.MetadataStruct) error {
	for _, param := range event.Params {
		if param.Name == "dispatch_error" {
			module, err := models.GetFieldFromAny("Module", param.Value)
			if err != nil {
				return fmt.Errorf(" GetFieldFromAny for module failed: %w", err)
			}
			errorIndexHex, err := models.GetFieldFromAny("error", module)
			if err != nil {
				return fmt.Errorf(" GetFieldFromAny for errorIndexHex failed: %w", err)
			}
			moduleIdxHex, err := models.GetFieldFromAny("index", module)
			if err != nil {
				return fmt.Errorf(" GetFieldFromAny for moduleIdx failed: %w", err)
			}
			errorIndex, err := models.ConvertFromHexToInt(errorIndexHex.(string))
			if err != nil {
				return fmt.Errorf(" ConvertFromHexToInt failed: %w", err)
			}
			if moduleIdxHex == nil {
				return fmt.Errorf("%w", errors.New("module index is nil"))
			}

			errorMessage := models.GetMessageByIndex(moduleIdxHex.(int), *errorIndex, meta)
			logger.Log().Error(errorMessage)
		}
	}
	return nil
}
