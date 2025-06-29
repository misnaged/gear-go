package root

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "root",
		Short:            "gear command",
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			bindFlags(cmd)
		},
	}
	return cmd
}
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
