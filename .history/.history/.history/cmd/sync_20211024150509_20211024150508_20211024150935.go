package cmd

import (
	"github.com/knqyf263/pet/config"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync cli-plugin-snippets",
	Long:  `Sync cli-plugin-snippets with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	return petSync.AutoSync(config.Conf.General.cli - plugin - snippetFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
