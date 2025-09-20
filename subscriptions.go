package gear_go

import (
	"errors"
	"fmt"
	gear_client "github.com/misnaged/gear-go/internal/client"
	gear_events "github.com/misnaged/gear-go/internal/events"
	"github.com/misnaged/gear-go/internal/models"
	gear_storage_methods "github.com/misnaged/gear-go/internal/storage/methods"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

// SubscriptionFunc is a subscription builder func
type SubscriptionFunc func() error

func (gear *Gear) MergeSubscriptionFunctions(fo ...SubscriptionFunc) {
	gear.subFuncs = append(gear.subFuncs, fo...)
}

func (gear *Gear) initEvents() {
	gear.events = gear_events.NewGearEvents(gear.GetMeta().GetMetadata())
}

// InitSubscriptions is the main func to handle ws subscriptions
// Required both AddResponseTypesAndMakeWsConnectionsPool
// and MergeSubscriptionFunctions had been called before using
// as it is shown in "example/code/example_subscription_upload" example
func (gear *Gear) InitSubscriptions() error {
	if !gear.config.Client.IsWebSocket {
		return errors.New("not a websocket client")
	}

	gear.initEvents()
	if len(gear.subFuncs) == 0 {
		return errors.New("subscription functions were not added. To add subscription function use gear.MergeSubscriptionFunctions")
	}
	for _, sub := range gear.subFuncs {

		go func() {
			if err := sub(); err != nil {
				logger.Log().Errorf("subfunc failed %v", err)
				return
			}
		}()

	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	fmt.Println("Shutdown signal received, starting graceful shutdown...")
	err := gear.wsClient.CloseAllConnection()
	if err != nil {
		return fmt.Errorf(" gear.wsClient.CloseAllConnection failed: %w", err)
	}
	close(gear.stop)

	return nil
}

// Pre-made functions:
// Every pre-made function could be included or excluded from SubscriptionFunc list
// Feel free to substitute these functions on whatever you need

// EventsSubscription implements state_subscribeStorage method
//
// as it is shown in "example/code/example_subscription_upload" (line 23) example
// function is added via MergeSubscriptionFunctions
func (gear *Gear) EventsSubscription() SubscriptionFunc {
	return func() error {
		storage := gear_storage_methods.NewStorage("System", "Events", gear.GetMeta(), gear.GetRPC())
		k, err := storage.GetStorageKey()
		if err != nil {
			logger.Log().Errorf(" storage.GetStorageKeys failed: %v", err)
		}

		storageSub, err := gear.wsClient.NewSubscriptionFunc("state_subscribeStorage", [][]string{{k}}, "subscribeStorage")
		if err != nil {
			return fmt.Errorf(" gear.wsClient.Subscribe failed: %w", err)
		}
		for resp := range storageSub {
			if resp.Params != nil {
				if err = gear.getResponseFromEventsSubscription(resp); err != nil {
					logger.Log().Errorf("gear.GetResponseFromEventsSubscription failed: %v", err)
					return nil
				}
			}
		}
		return nil
	}
}

// getResponseFromEventsSubscription gets change hash from subscription response
// and decodes it to useful payload (ExtrinsicFailed and UserMessageSent at this moment)
func (gear *Gear) getResponseFromEventsSubscription(resp *models.SubscriptionResponse) error {
	changes, err := models.GetChangesFromEvents(resp)
	if err != nil {
		return fmt.Errorf("gear.responsePoolRunner - GetChangesFromEvents failed: %w", err)
	}
	if changes != nil {
		events, err := gear.events.GetEvents(changes.ChangeHash)
		if err != nil {
			return fmt.Errorf("gear.responsePoolRunner - GetEvents failed: %w", err)
		}
		err = gear.events.Handle(events)
		if err != nil {
			return fmt.Errorf("gear.responsePoolRunner - HandleEvents failed: %w", err)
		}
	}
	return nil
}

func (gear *Gear) SubmitAndWatchExtrinsic(args []any, t string) SubscriptionFunc {
	return func() error {
		sSub, err := gear.GetWsClient().NewSubscriptionFunc("author_submitAndWatchExtrinsic", args, gear_client.ResponseType(t))
		if err != nil {
			return fmt.Errorf(" gear.wsClient.Subscribe failed: %w", err)
		}
		for resp := range sSub {
			logger.Log().Info("extrinsic", resp)
			if resp.Error != nil {
				logger.Log().Errorf("gear.wsClient.Subscribe failed: %v", resp.Error)
			}
		}
		return nil
	}
}

// TODO: looking for improvement of EnqueuedSubscriptions

// EnqueuedSubscriptions
func (gear *Gear) EnqueuedSubscriptions(methods []string, rtypes []gear_client.ResponseType, callNames, moduleNames []string, args [][]any) SubscriptionFunc {
	return func() error {
		for i := range rtypes {

			// we need to call CallBuilder firstly, because this function returns Signed transaction
			// with new nonce number
			// if we're calling it in advance - new nonce switching won't be happened
			str, err := gear.calls.CallBuilder(callNames[i], moduleNames[i], args[i])
			if err != nil {
				return fmt.Errorf("gear.wsClient.EnqueuedSubscriptions: %w", err)
			}

			ch, err := gear.wsClient.NewSubscriptionFunc(methods[i], []any{str}, rtypes[i])
			if err != nil {
				return fmt.Errorf("gear.wsClient.EnqueuedSubscriptions: %w", err)
			}
			for resp := range ch {
				logger.Log().Info("extrinsic", resp)
				if resp.Error != nil {
					logger.Log().Errorf("failed to send response: %v", resp.Error)
					gear.wsClient.CloseChannelByResponseType(rtypes[i])
					return errors.New("gear.wsClient.EnqueuedSubscriptions: failed to send response")
				}
				if models.IsFinalized(resp) {
					logger.Log().Info("subscription finalized, moving to next")
					gear.wsClient.CloseChannelByResponseType(rtypes[i])
					break
				}
			}
		}
		return nil
	}
}
