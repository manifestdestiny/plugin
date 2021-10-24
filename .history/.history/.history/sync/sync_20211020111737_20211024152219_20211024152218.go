package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/knqyf263/plugin/config"
	"github.com/pkg/errors"
)

// Client manages communication with the remote plugin repository
type Client interface {
	Getplugin() (*plugin, error)
	Uploadplugin(string) error
}

// plugin is the remote plugin
type plugin struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs plugins automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	plugin, err := client.Getplugin()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(plugin.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := plugin.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(plugin.Content)
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
	var plugins plugin.plugins
	if err := plugins.Load(); err != nil {
		return errors.Wrap(err, "Failed to load the local plugins")
	}

	body, err := plugins.ToString()
	if err != nil {
		return err
	}

	if err = client.Uploadplugin(body); err != nil {
		return errors.Wrap(err, "Failed to upload plugin")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	pluginFile := config.Conf.General.pluginFile

	var plugins plugin.plugins
	if err := plugins.Load(); err != nil {
		return err
	}
	body, err := plugins.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(pluginFile, []byte(content), os.ModePerm)
}
