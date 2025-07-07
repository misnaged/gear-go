package remove_target

import (
	"fmt"
	warp "github.com/misnaged/gear-go/pkg/warp_codegen"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-target",
		Short: "remove built target folder",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := warp.RemoveTargerFolder(); err != nil {
				return fmt.Errorf("remove Targer Folder failed: %w", err)
			}

			return nil
		},
	}

	return cmd
}
