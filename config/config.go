package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General GeneralConfig
	Gist    GistConfig
	GitLab  GitLabConfig
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	cli-plugin-pluginFile string `toml:"cli-plugin-pluginfile"`
	Editor      string `toml:"editor"`
	Column      int    `toml:"column"`
	SelectCmd   string `toml:"selectcmd"`
	Backend     string `toml:"backend"`
	SortBy      string `toml:"sortby"`
}

// GistConfig is a struct of config for Gist
type GistConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	GistID      string `toml:"gist_id"`
	Public      bool   `toml:"public"`
	AutoSync    bool   `toml:"auto_sync"`
}

// GitLabConfig is a struct of config for GitLabcli-plugin-plugin
type GitLabConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	Url         string `toml:"url"`
	ID          string `toml:"id"`
	Visibility  string `toml:"visibility"`
	AutoSync    bool   `toml:"auto_sync"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Query     string
	Command   bool
	Delimiter string
	OneLine   bool
	Color     bool
	Tag       bool
}

// Load loads a config toml
func (cfg *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		cfg.General.cli-plugin-pluginFile = expandPath(cfg.General.cli-plugin-pluginFile)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir, err := GetDefaultConfigDir()
	if err != nil {
		return errors.Wrap(err, "Failed to get the default config directory")
	}
	cfg.General.cli-plugin-pluginFile = filepath.Join(dir, "cli-plugin-plugin.toml")
	_, err = os.Create(cfg.General.cli-plugin-pluginFile)
	if err != nil {
		return errors.Wrap(err, "Failed to create a config file")
	}

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		if isCommandAvailable("sensible-editor") {
			cfg.General.Editor = "sensible-editor"
		} else {
			cfg.General.Editor = "vim"
		}
	}
	cfg.General.Column = 40
	cfg.General.SelectCmd = "fzf"
	cfg.General.Backend = "gist"

	cfg.Gist.FileName = "cli-plugin-plugin.toml"

	cfg.GitLab.FileName = "cli-plugin-plugin.toml"
	cfg.GitLab.Visibility = "private"

	return toml.NewEncoder(f).Encode(cfg)
}

// GetDefaultConfigDir returns the default config directory
func GetDefaultConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "plugin")
		}
		dir = filepath.Join(dir, "cli")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
