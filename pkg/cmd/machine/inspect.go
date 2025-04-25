package machine

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/helper"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	"github.com/spf13/cobra"
)

func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Short:   "Inspect a machine",
		Example: `baepo machine inspect <id>`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			if len(args) < 1 {
				a.IOStream.Error("For machine inspect, you must provide a machine ID.")
				return baepoerrors.InvalidArgsError
			}

			m, err := a.MachineClient.FindById(ctx, connect.NewRequest(&apiv1pb.MachineFindByIdRequest{
				MachineId: args[0],
			}))

			if err != nil {
				a.IOStream.Error(fmt.Sprintf("Inspecting machines: %v", err))
				return baepoerrors.MachineError
			}

			a.IOStream.Object(m.Msg.Machine, helper.MachineObjectFmt())

			return nil
		},
	}

	return cmd
}
