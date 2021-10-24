package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/snipplugin"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snipplugins",
	Long:  `Show all snipplugins`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var snipplugins snipplugin.Snipplugins
	if err := snipplugins.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, snipplugin := range snipplugins.Snipplugins {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(snipplugin.Description, col, "..."), col)
			command := runewidth.Truncate(snipplugin.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), snipplugin.Description)
			if strings.Contains(snipplugin.Command, "\n") {
				lines := strings.Split(snipplugin.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), snipplugin.Command)
			}
			if snipplugin.Tag != nil {
				tag := strings.Join(snipplugin.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if snipplugin.Output != "" {
				output := strings.Replace(snipplugin.Output, "\n", "\n             ", -1)
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
		`Display snipplugins in one line`)
}
