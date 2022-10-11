package root

import "github.com/spf13/cobra"

var JsonFlag bool

func Init(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolVarP(&JsonFlag, "json", "j", false, "json output")
}
