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

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List machines",
		Example: `baepo machine list`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			list, err := a.MachineClient.List(ctx, connect.NewRequest(&apiv1pb.MachineListRequest{
				WorkspaceId: *a.Config.CurrentContext.WorkspaceID,
			}))

			if err != nil {
				a.IOStream.Error(fmt.Sprintf("Listing machines: %v", err))
				return baepoerrors.MachineError
			}

			if len(list.Msg.Machines) == 0 {
				a.IOStream.Message("No machines found.")
				return nil
			}

			a.IOStream.Array(list.Msg.Machines, helper.MachineArrayFmt())

			return nil
		},
	}

	return cmd
}
