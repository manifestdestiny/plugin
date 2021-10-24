package cmd

import (
	"github.com/knqyf263/plugin/config"
	pluginSync "github.com/knqyf263/plugin/sync"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync cli-plugin-snipplugins",
	Long:  `Sync cli-plugin-snipplugins with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	return pluginSync.AutoSync(config.Conf.General.cli - plugin - snippluginFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
