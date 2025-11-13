package storage_subscribe

import (
	"fmt"
	gear_go "github.com/misnaged/gear-go"
	"github.com/spf13/cobra"
)

func Cmd(app *gear_go.Gear) *cobra.Command {
	return &cobra.Command{
		Use:   "subscribe-storage",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var types = []string{"subscribeStorage"}

			err := app.GetWsClient().AddResponseTypesAndMakeWsConnectionsPool(types...)
			if err != nil {
				return fmt.Errorf("error Adding response type: %v", err)
			}
			app.MergeSubscriptionFunctions(app.EventsSubscription())
			err = app.InitSubscriptions()
			if err != nil {
				return fmt.Errorf(" gear.ProcessEventsSubscription failed: %v", err)
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
		},
	}
}
