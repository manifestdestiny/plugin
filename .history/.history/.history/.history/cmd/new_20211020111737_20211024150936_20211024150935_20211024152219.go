package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/snipplugin"
	pluginSync "github.com/knqyf263/plugin/sync"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new COMMAND",
	Short: "Create a new snipplugin",
	Long:  `Create a new snipplugin (default: $HOME/.config/plugin/snipplugin.toml)`,
	RunE:  new,
}

func scan(message string) (string, error) {
	tempFile := "/tmp/plugin.tmp"
	if runtime.GOOS == "windows" {
		tempDir := os.Getenv("TEMP")
		tempFile = filepath.Join(tempDir, "plugin.tmp")
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:          message,
		HistoryFile:     tempFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", errors.New("canceled")
}

func new(cmd *cobra.Command, args []string) (err error) {
	var command string
	var description string
	var tags []string

	var snipplugins snipplugin.Snipplugins
	if err := snipplugins.Load(); err != nil {
		return err
	}

	if len(args) > 0 {
		command = strings.Join(args, " ")
		fmt.Fprintf(color.Output, "%s %s\n", color.YellowString("Command>"), command)
	} else {
		command, err = scan(color.YellowString("Command> "))
		if err != nil {
			return err
		}
	}
	description, err = scan(color.GreenString("Description> "))
	if err != nil {
		return err
	}

	if config.Flag.Tag {
		var t string
		if t, err = scan(color.CyanString("Tag> ")); err != nil {
			return err
		}
		tags = strings.Fields(t)
	}

	for _, s := range snipplugins.Snipplugins {
		if s.Description == description {
			return fmt.Errorf("Snipplugin [%s] already exists", description)
		}
	}

	newSnipplugin := snipplugin.SnippluginInfo{
		Description: description,
		Command:     command,
		Tag:         tags,
	}
	snipplugins.Snipplugins = append(snipplugins.Snipplugins, newSnipplugin)
	if err = snipplugins.Save(); err != nil {
		return err
	}

	snippluginFile := config.Conf.General.SnippluginFile
	if config.Conf.Gist.AutoSync {
		return pluginSync.AutoSync(snippluginFile)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&config.Flag.Tag, "tag", "t", false,
		`Display tag prompt (delimiter: space)`)
}