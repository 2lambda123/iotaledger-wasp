package init

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/iotaledger/wasp/tools/wasp-cli/cli/cliclients"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/iotaledger/wasp/tools/wasp-cli/util"
)

func initCheckVersionsCmd(waspVersion string) *cobra.Command {
	var debug bool
	cmd := &cobra.Command{
		Use:   "check-versions",
		Short: "checks the versions of wasp-cli and wasp nodes match",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// query every wasp node info endpoint and ensure the `Version` matches
			waspSettings := map[string]interface{}{}
			waspKey := viper.Sub("wasp")
			if waspKey != nil {
				waspSettings = waspKey.AllSettings()
			}
			if len(waspSettings) == 0 {
				log.Fatalf("no wasp node configured, you can add a node with `wasp-cli wasp add <name> <api url>`")
			}
			for nodeName := range waspSettings {
				version, _, err := cliclients.WaspClient(nodeName, debug).NodeApi.
					GetVersion(context.Background()).
					Execute()
				log.Check(err)

				if waspVersion == version.Version {
					log.Printf("Wasp-cli version matches Wasp {%s}\n", nodeName)
				} else {
					log.Printf("! -> Version mismatch with Wasp {%s}. cli version: %s, wasp version: %s\n", nodeName, waspVersion, version.Version)
				}
			}
		},
	}
	util.WithDebugFlag(cmd, &debug)
	return cmd
}
