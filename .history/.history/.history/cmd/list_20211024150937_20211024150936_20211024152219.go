package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/cli-plugin-snipplugin"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all cli-plugin-snipplugins",
	Long:  `Show all cli-plugin-snipplugins`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var cli-plugin-snipplugins cli-plugin-snipplugin.cli-plugin-snipplugins
	if err := cli-plugin-snipplugins.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, cli-plugin-snipplugin := range cli-plugin-snipplugins.cli-plugin-snipplugins {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(cli-plugin-snipplugin.Description, col, "..."), col)
			command := runewidth.Truncate(cli-plugin-snipplugin.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), cli-plugin-snipplugin.Description)
			if strings.Contains(cli-plugin-snipplugin.Command, "\n") {
				lines := strings.Split(cli-plugin-snipplugin.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), cli-plugin-snipplugin.Command)
			}
			if cli-plugin-snipplugin.Tag != nil {
				tag := strings.Join(cli-plugin-snipplugin.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if cli-plugin-snipplugin.Output != "" {
				output := strings.Replace(cli-plugin-snipplugin.Output, "\n", "\n             ", -1)
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
		`Display cli-plugin-snipplugins in one line`)
}
