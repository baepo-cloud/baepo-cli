package auth

import (
	"connectrpc.com/connect"
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	loginEmailFlag    string
	loginPasswordFlag string
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to a Baepo account",
		Example: `baepo auth login --email <email> --password <password>`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if loginEmailFlag == "" || loginPasswordFlag == "" {
				return cmd.Help()
			}

			ctx := cmd.Context()
			a := app.FromContext(ctx)

			login, err := a.AuthClient.Login(ctx, connect.NewRequest(&apiv1pb.AuthLoginRequest{
				Email:    loginEmailFlag,
				Password: loginPasswordFlag,
			}))

			if err != nil {
				a.IOStream.Error("Login failed: %v", err)
				return baepoerrors.AuthError
			}

			a.Config.CurrentContext.SecretKey = &login.Msg.SecretKey
			a.Config.CurrentContext.UserID = &login.Msg.UserId

			me, err := a.UserClient.Me(ctx, connect.NewRequest(&emptypb.Empty{}))
			if err != nil {
				a.IOStream.Error("Failed to get user info: %v", err)
				return baepoerrors.AuthError
			}

			a.Config.CurrentContext.WorkspaceID = &me.Msg.User.WorkspaceId

			err = config.SaveConfig(a.Config)
			if err != nil {
				a.IOStream.Error("Failed to save config: %v", err)
				return baepoerrors.ConfigError
			}

			a.IOStream.Message("Welcome back, %s!", me.Msg.User.FirstName)

			return nil
		},
	}

	cmd.Flags().StringVarP(&loginEmailFlag, "email", "e", "", "Email address")
	cmd.Flags().StringVarP(&loginPasswordFlag, "password", "p", "", "Password")

	return cmd
}
