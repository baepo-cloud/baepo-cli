package machine

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/helper"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	"github.com/spf13/cobra"
)

func newTerminateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "terminate <id>",
		Aliases: []string{"stop", "rm"},
		Short:   "Terminate a machine",
		Example: `# Terminate a machine
baepo machine terminate ID

# Terminate multiple machines
baepo machine terminate ID1 ID2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			if len(args) < 1 {
				a.IOStream.Error("You must provide at least one machine ID.")
				return baepoerrors.InvalidArgsError
			}

			machines := make([]*apiv1pb.Machine, 0)
			for _, machineID := range args {
				req := connect.NewRequest(&apiv1pb.MachineTerminateRequest{
					MachineId: machineID,
				})

				res, err := a.MachineClient.Terminate(ctx, req)
				if err != nil {
					a.IOStream.Error(fmt.Sprintf("Terminating machine %s: %v", machineID, err))
					continue
				}

				machines = append(machines, res.Msg.Machine)
			}

			if len(machines) > 1 {
				a.IOStream.Array(machines, helper.MachineMapping(), iostream.ObjectOptions{Full: false})
			} else {
				a.IOStream.Object(machines[0], helper.MachineMapping(), iostream.ObjectOptions{Full: true})
			}

			return nil
		},
	}

	return cmd
}
