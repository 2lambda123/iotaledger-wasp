package util

import "github.com/spf13/cobra"

func WithDebugFlag(cmd *cobra.Command, debug *bool) {
	cmd.Flags().BoolVar(debug, "debug", false, "enable debug logging (default: false)")
}
