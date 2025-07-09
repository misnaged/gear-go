package main

import (
	"errors"
	"fmt"
	"github.com/misnaged/gear-go/cmd/remove-rust"
	"github.com/misnaged/gear-go/cmd/remove-target"
	"github.com/misnaged/gear-go/cmd/root"
	"github.com/misnaged/gear-go/pkg/logger"
	warp "github.com/misnaged/gear-go/pkg/warp_codegen"

	"github.com/spf13/cobra"
	"os"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-rust-grpc",
		Short: "Generate rust grpc server file",

		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "true":
				tmpl := warp.NewTemplate(true)

				if err := tmpl.Generate(); err != nil {
					return fmt.Errorf("error generating rust grpc template: %w", err)
				}

			case "false":
				tmpl := warp.NewTemplate(false)

				if err := tmpl.Generate(); err != nil {
					return fmt.Errorf("error generating rust grpc template: %w", err)
				}
			default:

				//nolint:staticcheck
				return fmt.Errorf("%w", errors.New(fmt.Sprintf("unknown command flag: %s should be either `true` or `false`", args[0])))
			}

			return nil
		},
	}

	return cmd
}

func main() {

	rootCmd := root.Cmd()
	rootCmd.AddCommand(Cmd())
	rootCmd.AddCommand(remove_target.Cmd())
	rootCmd.AddCommand(remove_rust.Cmd())
	if err := rootCmd.Execute(); err != nil {
		logger.Log().Errorf("An error occurred %v", err)
		os.Exit(1)
	}
}
