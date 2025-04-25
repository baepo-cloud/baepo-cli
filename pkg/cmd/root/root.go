package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/cmd/auth"
	"github.com/baepo-cloud/baepo-cli/pkg/cmd/machine"
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	"github.com/spf13/cobra"
)

var (
	rootFlagCurrentContext = "default"
	rootJSONOutput         = false
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "baepo <command> <subcommand> [flags]",
		Short:         "Baepo CLI",
		Long:          `Work seamlessly with Baepo from the command line.`,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ baepo auth login --email lou@corp.com --password corp123Corp
			$ baepo machine create \
			  --name web-server \
			  --vcpus 2 \
			  --memory-mb 4096 \
			  --image ubuntu:latest \
			  --env "NODE_ENV=production" --env "PORT=3000" \
			  --command "/usr/bin/startup.sh" 
			$ baepo machine ls
		`),
		Annotations: map[string]string{
			"versionInfo": "0.0.1",
		},
		SilenceUsage: true,
		Version:      "0.0.1",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ios := iostream.New(rootJSONOutput)
			cfg, err := config.LoadConfig(rootFlagCurrentContext)

			a := app.NewApp(cfg, ios)
			cmd.SetContext(app.SaveToContext(a, cmd.Context()))

			if err != nil {
				a.IOStream.Error("failed to load config: %v", err)
				return baepoerrors.ConfigError
			}

			p := cmd.Parent()
			if cfg.CurrentContext.SecretKey == nil && (p == nil || p.Name() != "auth") {
				a.IOStream.Error("No secret key found in the current context. Please login to Baepo using the command: baepo auth login --email <email> --password <password>")
				return baepoerrors.AuthError
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&rootFlagCurrentContext, "context", "x", "default", "Set the current context")
	cmd.PersistentFlags().BoolVarP(&rootJSONOutput, "json", "j", false, "Output in JSON format")

	cmd.AddCommand(auth.NewAuthCmd())
	cmd.AddCommand(machine.NewMachineCmd())

	return cmd
}
