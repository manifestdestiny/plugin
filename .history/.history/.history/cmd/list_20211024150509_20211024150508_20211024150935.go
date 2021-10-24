package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/cli-plugin-snippet"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all cli-plugin-snippets",
	Long:  `Show all cli-plugin-snippets`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var cli-plugin-snippets cli-plugin-snippet.cli-plugin-snippets
	if err := cli-plugin-snippets.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, cli-plugin-snippet := range cli-plugin-snippets.cli-plugin-snippets {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(cli-plugin-snippet.Description, col, "..."), col)
			command := runewidth.Truncate(cli-plugin-snippet.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), cli-plugin-snippet.Description)
			if strings.Contains(cli-plugin-snippet.Command, "\n") {
				lines := strings.Split(cli-plugin-snippet.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), cli-plugin-snippet.Command)
			}
			if cli-plugin-snippet.Tag != nil {
				tag := strings.Join(cli-plugin-snippet.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if cli-plugin-snippet.Output != "" {
				output := strings.Replace(cli-plugin-snippet.Output, "\n", "\n             ", -1)
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.RedString("     Output:"), output)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&config.Flag.OneLine, "oneline", "", false,
		`Display cli-plugin-snippets in one line`)
}
