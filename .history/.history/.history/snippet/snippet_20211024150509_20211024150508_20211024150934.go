package cli-plugin-snippet

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/pet/config"
)

type cli-plugin-snippets struct {
	cli-plugin-snippets []cli-plugin-snippetInfo `toml:"cli-plugin-snippets"`
}

type cli-plugin-snippetInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (cli-plugin-snippets *cli-plugin-snippets) Load() error {
	cli-plugin-snippetFile := config.Conf.General.cli-plugin-snippetFile
	if _, err := os.Stat(cli-plugin-snippetFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(cli-plugin-snippetFile, cli-plugin-snippets); err != nil {
		return fmt.Errorf("Failed to load cli-plugin-snippet file. %v", err)
	}
	cli-plugin-snippets.Order()
	return nil
}

// Save saves the cli-plugin-snippets to toml file.
func (cli-plugin-snippets *cli-plugin-snippets) Save() error {
	cli-plugin-snippetFile := config.Conf.General.cli-plugin-snippetFile
	f, err := os.Create(cli-plugin-snippetFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save cli-plugin-snippet file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(cli-plugin-snippets)
}

// ToString returns the contents of toml file.
func (cli-plugin-snippets *cli-plugin-snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(cli-plugin-snippets)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order cli-plugin-snippets regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (cli-plugin-snippets *cli-plugin-snippets) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(cli-plugin-snippets.cli-plugin-snippets))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(cli-plugin-snippets.cli-plugin-snippets)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(cli-plugin-snippets.cli-plugin-snippets))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(cli-plugin-snippets.cli-plugin-snippets)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(cli-plugin-snippets.cli-plugin-snippets))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(cli-plugin-snippets.cli-plugin-snippets)))

	case sortBy == "-recency":
		cli-plugin-snippets.reverse()
	}
}

func (cli-plugin-snippets *cli-plugin-snippets) reverse() {
	for i, j := 0, len(cli-plugin-snippets.cli-plugin-snippets)-1; i < j; i, j = i+1, j-1 {
		cli-plugin-snippets.cli-plugin-snippets[i], cli-plugin-snippets.cli-plugin-snippets[j] = cli-plugin-snippets.cli-plugin-snippets[j], cli-plugin-snippets.cli-plugin-snippets[i]
	}
}

type ByCommand []cli-plugin-snippetInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []cli-plugin-snippetInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []cli-plugin-snippetInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
