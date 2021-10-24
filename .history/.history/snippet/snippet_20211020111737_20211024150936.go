package plugin

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/256bit-src/plugin/config"
	"github.com/BurntSushi/toml"
)

type plugins struct {
	plugins []pluginInfo `toml:"plugins"`
}

type pluginInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (plugins *plugins) Load() error {
	pluginFile := config.Conf.General.pluginFile
	if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(pluginFile, plugins); err != nil {
		return fmt.Errorf("Failed to load plugin file. %v", err)
	}
	plugins.Order()
	return nil
}

// Save saves the plugins to toml file.
func (plugins *plugins) Save() error {
	pluginFile := config.Conf.General.pluginFile
	f, err := os.Create(pluginFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save plugin file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(plugins)
}

// ToString returns the contents of toml file.
func (plugins *plugins) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(plugins)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order plugins regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (plugins *plugins) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(plugins.plugins))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(plugins.plugins)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(plugins.plugins))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(plugins.plugins)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(plugins.plugins))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(plugins.plugins)))

	case sortBy == "-recency":
		plugins.reverse()
	}
}

func (plugins *plugins) reverse() {
	for i, j := 0, len(plugins.plugins)-1; i < j; i, j = i+1, j-1 {
		plugins.plugins[i], plugins.plugins[j] = plugins.plugins[j], plugins.plugins[i]
	}
}

type ByCommand []pluginInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []pluginInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []pluginInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
