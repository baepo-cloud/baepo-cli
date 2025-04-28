package contextcmd

import (
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/helper"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			var list []*helper.ContextFmt
			current := a.Config.Context
			for key, context := range a.Config.Contexts {
				list = append(list, &helper.ContextFmt{
					Name:    key,
					Current: key == current,
					Value:   *context,
				})
			}

			if len(list) == 0 {
				a.IOStream.Message("No contexts found.")
				return nil
			}

			a.IOStream.Array(list, helper.ContextFmtMapping(), iostream.ObjectOptions{Full: false})

			return nil
		},
	}

	return cmd
}
