package gear_events

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/pkg/logger"
)

type IEvent interface {
	GetEvents(hex string) ([]*models.Event, error)
	Handle(events []*models.Event) error
}

type GearEvent struct {
	meta    *types.MetadataStruct
	decoder scalecodec.EventsDecoder
	opts    *types.ScaleDecoderOption
}

func NewGearEvents(metadataStruct *types.MetadataStruct) IEvent {
	return &GearEvent{
		meta:    metadataStruct,
		decoder: scalecodec.EventsDecoder{},
		opts:    &types.ScaleDecoderOption{Metadata: metadataStruct},
	}
}

// GetEvents is
func (ev *GearEvent) GetEvents(hex string) ([]*models.Event, error) {
	ev.decoder.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(hex)}, ev.opts)
	ev.decoder.Process()

	events, err := convertToEvents(ev.decoder.Value.([]any))
	if err != nil {
		return nil, fmt.Errorf(" gear.ConvertToEvents failed: %w", err)
	}
	return events, nil
}

//func (ev *GearEvent) handleExtrinsicSuccess(event *models.Event) error {
//	for _, param := range event.Params {
//		fmt.Println(param.Name, param.Value)
//	}
//	return nil
//}

func (ev *GearEvent) handleExtrinsicFailed(event *models.Event) error {
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

			errorMessage := models.GetMessageByIndex(moduleIdxHex.(int), *errorIndex, ev.meta)
			logger.Log().Error(errorMessage)
		}
	}
	return nil
}

func (ev *GearEvent) Handle(events []*models.Event) error {
	for _, event := range events {
		switch event.EventID {
		case "ExtrinsicFailed":
			err := ev.handleExtrinsicFailed(event)
			if err != nil {
				return fmt.Errorf(" gear.HandleExtrinsicFailed failed: %w", err)
			}
			//case "ExtrinsicSuccess":
			//	err := ev.handleExtrinsicSuccess(event)
			//	if err != nil {
			//		return fmt.Errorf(" gear.HandleExtrinsicSuccess failed: %w", err)
			//	}
		}
	}
	return nil
}
func convertToEvents(src []any) ([]*models.Event, error) {
	var events []*models.Event
	for i, raw := range src {
		m, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("element %d is not a map[string]any", i)
		}

		e := models.Event{}

		if v, ok := m["event_id"].(string); ok {
			e.EventID = v
		}
		if v, ok := m["module_id"].(string); ok {
			e.ModuleID = v
		}
		if v, ok := m["event_idx"].(int); ok {
			e.EventIndex = v
		}
		if v, ok := m["extrinsic_idx"].(int); ok {
			e.ExtrinsicIndex = v
		}
		if v, ok := m["params"].([]scalecodec.EventParam); ok {
			e.Params = v
		}

		events = append(events, &e)
	}

	return events, nil
}
