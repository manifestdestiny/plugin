package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/plugin/config"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all plugins",
	Long:  `Show all plugins`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var plugins plugin.plugins
	if err := plugins.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, plugin := range plugins.plugins {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(plugin.Description, col, "..."), col)
			command := runewidth.Truncate(plugin.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), plugin.Description)
			if strings.Contains(plugin.Command, "\n") {
				lines := strings.Split(plugin.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), plugin.Command)
			}
			if plugin.Tag != nil {
				tag := strings.Join(plugin.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if plugin.Output != "" {
				output := strings.Replace(plugin.Output, "\n", "\n             ", -1)
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
		`Display plugins in one line`)
}
