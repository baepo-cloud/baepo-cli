package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	"github.com/baepo-cloud/baepo-proto/go/baepo/api/v1/apiv1pbconnect"
)

const (
	ctxKey = "baepo-cli"
)

type App struct {
	Config   *config.Config
	IOStream *iostream.IOStream

	AuthClient    apiv1pbconnect.AuthServiceClient
	UserClient    apiv1pbconnect.UserServiceClient
	MachineClient apiv1pbconnect.MachineServiceClient
}

func NewApp(cfg *config.Config, ioStream *iostream.IOStream) *App {
	return &App{
		Config:   cfg,
		IOStream: ioStream,

		AuthClient:    apiv1pbconnect.NewAuthServiceClient(http.DefaultClient, cfg.CurrentContext.URL),
		UserClient:    apiv1pbconnect.NewUserServiceClient(http.DefaultClient, cfg.CurrentContext.URL, AuthenticatedClientOption(cfg)),
		MachineClient: apiv1pbconnect.NewMachineServiceClient(http.DefaultClient, cfg.CurrentContext.URL, AuthenticatedClientOption(cfg)),
	}
}

// AuthenticatedClientOption returns a connect.ClientOption that automatically adds
// the authentication header to all requests using the provided config.
func AuthenticatedClientOption(cfg *config.Config) connect.ClientOption {
	return connect.WithInterceptors(
		connect.UnaryInterceptorFunc(
			func(next connect.UnaryFunc) connect.UnaryFunc {
				return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
					if cfg.CurrentContext.UserID != "" && cfg.CurrentContext.SecretKey != "" {
						token := base64.StdEncoding.EncodeToString(
							[]byte(fmt.Sprintf("%s:%s", cfg.CurrentContext.UserID, cfg.CurrentContext.SecretKey)),
						)

						req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
					}
					return next(ctx, req)
				}
			},
		),
	)
}

func SaveToContext(a *App, ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, a)
}

func FromContext(ctx context.Context) *App {
	if ctx == nil {
		return nil
	}
	app, ok := ctx.Value(ctxKey).(*App)
	if !ok {
		return nil
	}
	return app
}
