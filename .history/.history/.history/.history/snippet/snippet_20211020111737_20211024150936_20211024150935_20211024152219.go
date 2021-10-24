package snipplugin

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/plugin/config"
)

type Snipplugins struct {
	Snipplugins []SnippluginInfo `toml:"snipplugins"`
}

type SnippluginInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (snipplugins *Snipplugins) Load() error {
	snippluginFile := config.Conf.General.SnippluginFile
	if _, err := os.Stat(snippluginFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(snippluginFile, snipplugins); err != nil {
		return fmt.Errorf("Failed to load snipplugin file. %v", err)
	}
	snipplugins.Order()
	return nil
}

// Save saves the snipplugins to toml file.
func (snipplugins *Snipplugins) Save() error {
	snippluginFile := config.Conf.General.SnippluginFile
	f, err := os.Create(snippluginFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save snipplugin file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(snipplugins)
}

// ToString returns the contents of toml file.
func (snipplugins *Snipplugins) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(snipplugins)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order snipplugins regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (snipplugins *Snipplugins) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(snipplugins.Snipplugins))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(snipplugins.Snipplugins)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(snipplugins.Snipplugins))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(snipplugins.Snipplugins)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(snipplugins.Snipplugins))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(snipplugins.Snipplugins)))

	case sortBy == "-recency":
		snipplugins.reverse()
	}
}

func (snipplugins *Snipplugins) reverse() {
	for i, j := 0, len(snipplugins.Snipplugins)-1; i < j; i, j = i+1, j-1 {
		snipplugins.Snipplugins[i], snipplugins.Snipplugins[j] = snipplugins.Snipplugins[j], snipplugins.Snipplugins[i]
	}
}

type ByCommand []SnippluginInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []SnippluginInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []SnippluginInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
