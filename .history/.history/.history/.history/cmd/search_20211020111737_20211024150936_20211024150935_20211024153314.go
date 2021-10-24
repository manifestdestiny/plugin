package cmd

import (
	"fmt"
	"strings"

	"github.com/256bit-src/plugin/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alessio/shellescape.v1"
)

var delimiter string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search snipplugins",
	Long:  `Search snipplugins interactively (default filtering tool: peco)`,
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}
	commands, err := filter(options)
	if err != nil {
		return err
	}

	fmt.Print(strings.Join(commands, flag.Delimiter))
	if terminal.IsTerminal(1) {
		fmt.Print("\n")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	searchCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	searchCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}
