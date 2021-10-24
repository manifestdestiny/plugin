package cmd

import (
	"github.com/256bit-src/pet/config"
	petSync "github.com/256bit-src/pet/sync"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	return petSync.AutoSync(config.Conf.General.SnippetFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
