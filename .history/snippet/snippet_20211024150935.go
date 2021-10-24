package cli-plugin-plugin

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/256bit-src/plugin/config"
)

type cli-plugin-plugins struct {
	cli-plugin-plugins []cli-plugin-pluginInfo `toml:"cli-plugin-plugins"`
}

type cli-plugin-pluginInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (cli-plugin-plugins *cli-plugin-plugins) Load() error {
	cli-plugin-pluginFile := config.Conf.General.cli-plugin-pluginFile
	if _, err := os.Stat(cli-plugin-pluginFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(cli-plugin-pluginFile, cli-plugin-plugins); err != nil {
		return fmt.Errorf("Failed to load cli-plugin-plugin file. %v", err)
	}
	cli-plugin-plugins.Order()
	return nil
}

// Save saves the cli-plugin-plugins to toml file.
func (cli-plugin-plugins *cli-plugin-plugins) Save() error {
	cli-plugin-pluginFile := config.Conf.General.cli-plugin-pluginFile
	f, err := os.Create(cli-plugin-pluginFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save cli-plugin-plugin file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(cli-plugin-plugins)
}

// ToString returns the contents of toml file.
func (cli-plugin-plugins *cli-plugin-plugins) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(cli-plugin-plugins)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order cli-plugin-plugins regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (cli-plugin-plugins *cli-plugin-plugins) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(cli-plugin-plugins.cli-plugin-plugins))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(cli-plugin-plugins.cli-plugin-plugins)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(cli-plugin-plugins.cli-plugin-plugins))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(cli-plugin-plugins.cli-plugin-plugins)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(cli-plugin-plugins.cli-plugin-plugins))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(cli-plugin-plugins.cli-plugin-plugins)))

	case sortBy == "-recency":
		cli-plugin-plugins.reverse()
	}
}

func (cli-plugin-plugins *cli-plugin-plugins) reverse() {
	for i, j := 0, len(cli-plugin-plugins.cli-plugin-plugins)-1; i < j; i, j = i+1, j-1 {
		cli-plugin-plugins.cli-plugin-plugins[i], cli-plugin-plugins.cli-plugin-plugins[j] = cli-plugin-plugins.cli-plugin-plugins[j], cli-plugin-plugins.cli-plugin-plugins[i]
	}
}

type ByCommand []cli-plugin-pluginInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []cli-plugin-pluginInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []cli-plugin-pluginInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
