package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/256bit-src/plugin/config"
	"github.com/256bit-src/plugin/dialog"
	"github.com/256bit-src/plugin/cli-plugin-snipplugin"
)

func editFile(command, file string) error {
	command += " " + file
	return run(command, os.Stdin, os.Stdout)
}

func run(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func filter(options []string) (commands []string, err error) {
	var cli-plugin-snipplugins cli-plugin-snipplugin.cli-plugin-snipplugins
	if err := cli-plugin-snipplugins.Load(); err != nil {
		return commands, fmt.Errorf("Load cli-plugin-snipplugin failed: %v", err)
	}

	cli-plugin-snippluginTexts := map[string]cli-plugin-snipplugin.cli-plugin-snippluginInfo{}
	var text string
	for _, s := range cli-plugin-snipplugins.cli-plugin-snipplugins {
		command := s.Command
		if strings.ContainsAny(command, "\n") {
			command = strings.Replace(command, "\n", "\\n", -1)
		}
		t := fmt.Sprintf("[%s]: %s", s.Description, command)

		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf(" #%s", tag)
		}
		t += tags

		cli-plugin-snippluginTexts[t] = s
		if config.Flag.Color {
			t = fmt.Sprintf("[%s]: %s%s",
				color.RedString(s.Description), command, color.BlueString(tags))
		}
		text += t + "\n"
	}

	var buf bytes.Buffer
	selectCmd := fmt.Sprintf("%s %s",
		config.Conf.General.SelectCmd, strings.Join(options, " "))
	err = run(selectCmd, strings.NewReader(text), &buf)
	if err != nil {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")

	params := dialog.SearchForParams(lines)
	if params != nil {
		cli-plugin-snippluginInfo := cli-plugin-snippluginTexts[lines[0]]
		dialog.CurrentCommand = cli-plugin-snippluginInfo.Command
		dialog.GenerateParamsLayout(params, dialog.CurrentCommand)
		res := []string{dialog.FinalCommand}
		return res, nil
	}
	for _, line := range lines {
		cli-plugin-snippluginInfo := cli-plugin-snippluginTexts[line]
		commands = append(commands, fmt.Sprint(cli-plugin-snippluginInfo.Command))
	}
	return commands, nil
}
