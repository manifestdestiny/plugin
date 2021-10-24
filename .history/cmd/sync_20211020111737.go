package cmd

import (
	"github.com/256bit-src/plugin/config"
	pluginSync "github.com/256bit-src/plugin/sync"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync plugins",
	Long:  `Sync plugins with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	return pluginSync.AutoSync(config.Conf.General.pluginFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
