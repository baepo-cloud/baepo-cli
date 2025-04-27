package contextcmd

import (
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/helper"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	"github.com/spf13/cobra"
)

func NewContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage your contexts",

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			if len(args) < 1 {
				c := &helper.ContextFmt{
					Name:  a.Config.Context,
					Value: *a.Config.CurrentContext,
				}

				a.IOStream.Object(c, helper.ContextFmtMapping(), iostream.ObjectOptions{Full: true})
			} else {
				c, ok := a.Config.Contexts[args[0]]
				if !ok {
					a.IOStream.Error("Context not found")
					return baepoerrors.InvalidArgsError
				}

				cf := &helper.ContextFmt{
					Name:  args[0],
					Value: *c,
				}

				a.IOStream.Object(cf, helper.ContextFmtMapping(), iostream.ObjectOptions{Full: true})
			}

			return nil
		},
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUseCmd())

	return cmd

}
