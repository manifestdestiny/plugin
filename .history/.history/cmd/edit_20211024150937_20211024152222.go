package cmd

import (
	"io/ioutil"

	"github.com/256bit-src/plugin/config"
	pluginSync "github.com/256bit-src/plugin/sync"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit cli-plugin-plugin file",
	Long:  `Edit cli-plugin-plugin file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	cli-plugin-pluginFile := config.Conf.General.cli-plugin-pluginFile

	// file content before editing
	before := fileContent(cli-plugin-pluginFile)

	err = editFile(editor, cli-plugin-pluginFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(cli-plugin-pluginFile)

	// return if same file content
	if before == after {
		return nil
	}

	if config.Conf.Gist.AutoSync {
		return pluginSync.AutoSync(cli-plugin-pluginFile)
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := ioutil.ReadFile(fname)
	return string(data)
}

func init() {
	RootCmd.AddCommand(editCmd)
}
