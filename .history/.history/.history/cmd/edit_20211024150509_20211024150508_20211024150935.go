package cmd

import (
	"io/ioutil"

	"github.com/knqyf263/pet/config"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit cli-plugin-snippet file",
	Long:  `Edit cli-plugin-snippet file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	cli-plugin-snippetFile := config.Conf.General.cli-plugin-snippetFile

	// file content before editing
	before := fileContent(cli-plugin-snippetFile)

	err = editFile(editor, cli-plugin-snippetFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(cli-plugin-snippetFile)

	// return if same file content
	if before == after {
		return nil
	}

	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(cli-plugin-snippetFile)
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
