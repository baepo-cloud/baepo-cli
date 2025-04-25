package machine

import (
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create machine",
		Example: `
# Create a machine with a single container
baepo machine create --name myapp --cpus 2 --memory 2048 --image nginx:latest --env KEY1=value1 --env KEY2=value2 --start

# Create a machine with a container that has a healthcheck
baepo machine create --name myapp --cpus 2 --memory 2048 --image myapp:latest --health-port 8080 --health-path /health

# Create a machine with multiple containers using JSON
baepo machine create --name mydb --cpus 4 --memory 8192 --containers '[{"image":"postgres:14","env":{"POSTGRES_PASSWORD":"secret"}},{"image":"redis:alpine"}]' --start
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx := cmd.Context()
			// a := app.FromContext(ctx)

			// todo

			return nil
		},
	}

	return cmd
}
