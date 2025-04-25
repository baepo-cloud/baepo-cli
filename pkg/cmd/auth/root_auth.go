package auth

import (
	"github.com/spf13/cobra"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Auth manage connection to a Baepo account",
	}

	cmd.AddCommand(newLoginCmd())

	return cmd

}
