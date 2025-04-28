package contextcmd

import (
	"fmt"

	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var workspaceID string
	var userID string
	var secretKey string
	var current bool
	var url string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create context",
		Example: `
# Create a blank new context (you will need to login then)
baepo context create mycompany --current
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			if len(args) < 1 {
				a.IOStream.Error("You must provide a name for the context.")
				return baepoerrors.InvalidArgsError
			}

			name := args[0]

			_, exists := a.Config.Contexts[name]
			if exists {
				a.IOStream.Error("Context with the name '%s' already exists.", name)
				return baepoerrors.InvalidArgsError
			}

			// user id and secret key must be provided together if one is provided
			if (userID != "" && secretKey == "") || (userID == "" && secretKey != "") {
				a.IOStream.Error("Both user ID and secret key must be provided together.")
				return baepoerrors.InvalidArgsError
			}

			// Create a new context
			newContext := *config.DefaultContext

			if workspaceID != "" {
				newContext.WorkspaceID = workspaceID
			}

			if userID != "" {
				newContext.UserID = userID
				newContext.SecretKey = secretKey
			}

			if url != "" {
				newContext.URL = url
			}

			a.Config.Contexts[name] = &newContext

			if current {
				a.Config.CurrentContext = &newContext
			}

			err := config.SaveConfig(a.Config)
			if err != nil {
				a.IOStream.Error("Failed to save config: %v", err)
				return baepoerrors.ConfigError
			}

			message := fmt.Sprintf("Context '%s' created.", name)
			if current {
				message = fmt.Sprintf("Context '%s' created and set as current context.", name)
			}
			a.IOStream.Message(message)

			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceID, "workspace-id", "w", "", "Workspace ID")
	cmd.Flags().StringVar(&userID, "user-id", "", "User ID")
	cmd.Flags().StringVarP(&secretKey, "s", "", "", "Secret Key")
	cmd.Flags().StringVarP(&url, "url", "u", "", "Baepo API URL")
	cmd.Flags().BoolVarP(&current, "current", "c", false, "Set this context as the current context")

	return cmd
}
