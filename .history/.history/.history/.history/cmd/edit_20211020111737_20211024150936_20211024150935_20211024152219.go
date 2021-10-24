package cmd

import (
	"io/ioutil"

	"github.com/knqyf263/plugin/config"
	pluginSync "github.com/knqyf263/plugin/sync"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snipplugin file",
	Long:  `Edit snipplugin file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	snippluginFile := config.Conf.General.SnippluginFile

	// file content before editing
	before := fileContent(snippluginFile)

	err = editFile(editor, snippluginFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(snippluginFile)

	// return if same file content
	if before == after {
		return nil
	}

	if config.Conf.Gist.AutoSync {
		return pluginSync.AutoSync(snippluginFile)
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
