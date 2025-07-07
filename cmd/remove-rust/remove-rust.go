package remove_rust

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-rust-gen",
		Short: "remove generated grpc server file",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.Remove("lib/server_grpc/src/gear_grpc.rs"); err != nil {
				return fmt.Errorf("%w", err)
			}
			return nil
		},
	}

	return cmd
}
