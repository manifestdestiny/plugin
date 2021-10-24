package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/cli-plugin-plugin"
	"github.com/pkg/errors"
)

// Client manages communication with the remote cli-plugin-plugin repository
type Client interface {
	Getcli-plugin-plugin() (*cli-plugin-plugin, error)
	Uploadcli-plugin-plugin(string) error
}

// cli-plugin-plugin is the remote cli-plugin-plugin
type cli-plugin-plugin struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs cli-plugin-plugins automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	cli-plugin-plugin, err := client.Getcli-plugin-plugin()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(cli-plugin-plugin.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := cli-plugin-plugin.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(cli-plugin-plugin.Content)
	default:
		return nil
	}
}

// NewSyncClient returns Client
func NewSyncClient() (Client, error) {
	if config.Conf.General.Backend == "gitlab" {
		client, err := NewGitLabClient()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize GitLab client")
		}
		return client, nil
	}
	client, err := NewGistClient()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize Gist client")
	}
	return client, nil
}

func upload(client Client) (err error) {
	var cli-plugin-plugins cli-plugin-plugin.cli-plugin-plugins
	if err := cli-plugin-plugins.Load(); err != nil {
		return errors.Wrap(err, "Failed to load the local cli-plugin-plugins")
	}

	body, err := cli-plugin-plugins.ToString()
	if err != nil {
		return err
	}

	if err = client.Uploadcli-plugin-plugin(body); err != nil {
		return errors.Wrap(err, "Failed to upload cli-plugin-plugin")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	cli-plugin-pluginFile := config.Conf.General.cli-plugin-pluginFile

	var cli-plugin-plugins cli-plugin-plugin.cli-plugin-plugins
	if err := cli-plugin-plugins.Load(); err != nil {
		return err
	}
	body, err := cli-plugin-plugins.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(cli-plugin-pluginFile, []byte(content), os.ModePerm)
}
