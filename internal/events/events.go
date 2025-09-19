package gear_events

import (
	"fmt"
	"github.com/goccy/go-json"
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

func (ev *GearEvent) handleGearUserMessage(event *models.Event) error {
	for _, param := range event.Params {
		if param.Name == "message" {
			b, err := json.Marshal(param.Value)
			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}
			var message models.Message
			if err = json.Unmarshal(b, &message); err != nil {
				return fmt.Errorf("failed to unmarshal message: %w", err)
			}
			if message.Details.Code.Error != nil {
				//TODO: error handler for messages
				logger.Log().Errorf("gear.handleGearUserMessage failed: %v", message.Details.Code.Error.Execution)
				continue
			}
			logger.Log().Infof("gear.handleGearUserMessage - message: %v", message)
		}
	}
	return nil
}

func (ev *GearEvent) Handle(events []*models.Event) error {
	// todo: refactor hardcode
	for _, event := range events {
		switch event.EventID {
		case "ExtrinsicFailed":
			err := handleExtrinsicFailed(event, ev.meta)
			if err != nil {
				return fmt.Errorf(" gear.HandleExtrinsicFailed failed: %w", err)
			}
		case "UserMessageSent":
			err := ev.handleGearUserMessage(event)
			if err != nil {
				return fmt.Errorf(" gear.HandleGearUserMessage failed: %w", err)
			}

		default:
			if event.EventID != "ExtrinsicSuccess" {
				logger.Log().Infof("gear.HandleGearUserMessage - event: %v", event)
			}
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
