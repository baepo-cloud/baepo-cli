package contextcmd

import (
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	"github.com/spf13/cobra"
)

func newUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Use context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			if len(args) < 1 {
				a.IOStream.Error("You must provide a name for the context.")
				return baepoerrors.InvalidArgsError
			}

			name := args[0]

			_, exists := a.Config.Contexts[name]
			if !exists {
				a.IOStream.Error("Context with the name '%s' does not exist.", name)
				return baepoerrors.InvalidArgsError
			}

			a.Config.Context = name

			err := config.SaveConfig(a.Config)
			if err != nil {
				a.IOStream.Error("Failed to save config: %v", err)
				return baepoerrors.ConfigError
			}

			a.IOStream.Message("Switched to context '%s'", name)

			return nil
		},
	}

	return cmd
}
