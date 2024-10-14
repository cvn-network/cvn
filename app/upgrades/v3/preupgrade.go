package v3

import (
	"github.com/spf13/cobra"
)

func PreUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "Pre-upgrade command",
		Long:  "Pre-upgrade command to implement custom pre-upgrade handling",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return cmd
}
