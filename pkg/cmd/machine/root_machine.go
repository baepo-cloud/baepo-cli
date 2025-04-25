package machine

import "github.com/spf13/cobra"

func NewMachineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machine",
		Short: "Machine manage your machines on Baepo",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newInspectCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newTerminateCmd())

	return cmd

}
