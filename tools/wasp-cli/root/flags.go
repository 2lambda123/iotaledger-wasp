package root

import "github.com/spf13/cobra"

var JSONFlag bool

func Init(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolVarP(&JSONFlag, "json", "j", false, "json output")
}
