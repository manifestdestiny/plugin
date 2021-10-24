package cli-plugin-snipplugin

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/256bit-src/plugin/config"
)

type cli-plugin-snipplugins struct {
	cli-plugin-snipplugins []cli-plugin-snippluginInfo `toml:"cli-plugin-snipplugins"`
}

type cli-plugin-snippluginInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (cli-plugin-snipplugins *cli-plugin-snipplugins) Load() error {
	cli-plugin-snippluginFile := config.Conf.General.cli-plugin-snippluginFile
	if _, err := os.Stat(cli-plugin-snippluginFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(cli-plugin-snippluginFile, cli-plugin-snipplugins); err != nil {
		return fmt.Errorf("Failed to load cli-plugin-snipplugin file. %v", err)
	}
	cli-plugin-snipplugins.Order()
	return nil
}

// Save saves the cli-plugin-snipplugins to toml file.
func (cli-plugin-snipplugins *cli-plugin-snipplugins) Save() error {
	cli-plugin-snippluginFile := config.Conf.General.cli-plugin-snippluginFile
	f, err := os.Create(cli-plugin-snippluginFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save cli-plugin-snipplugin file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(cli-plugin-snipplugins)
}

// ToString returns the contents of toml file.
func (cli-plugin-snipplugins *cli-plugin-snipplugins) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(cli-plugin-snipplugins)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order cli-plugin-snipplugins regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (cli-plugin-snipplugins *cli-plugin-snipplugins) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(cli-plugin-snipplugins.cli-plugin-snipplugins))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(cli-plugin-snipplugins.cli-plugin-snipplugins)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(cli-plugin-snipplugins.cli-plugin-snipplugins))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(cli-plugin-snipplugins.cli-plugin-snipplugins)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(cli-plugin-snipplugins.cli-plugin-snipplugins))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(cli-plugin-snipplugins.cli-plugin-snipplugins)))

	case sortBy == "-recency":
		cli-plugin-snipplugins.reverse()
	}
}

func (cli-plugin-snipplugins *cli-plugin-snipplugins) reverse() {
	for i, j := 0, len(cli-plugin-snipplugins.cli-plugin-snipplugins)-1; i < j; i, j = i+1, j-1 {
		cli-plugin-snipplugins.cli-plugin-snipplugins[i], cli-plugin-snipplugins.cli-plugin-snipplugins[j] = cli-plugin-snipplugins.cli-plugin-snipplugins[j], cli-plugin-snipplugins.cli-plugin-snipplugins[i]
	}
}

type ByCommand []cli-plugin-snippluginInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []cli-plugin-snippluginInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []cli-plugin-snippluginInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
