package serve

import (
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/pkg/logger"

	"github.com/spf13/cobra"
)

func Cmd(app *gear_go.Gear) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run Application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			logger.Log().Info()
		},
	}
}
