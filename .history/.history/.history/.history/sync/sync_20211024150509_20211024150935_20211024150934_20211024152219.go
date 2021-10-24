package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/cli-plugin-snipplugin"
	"github.com/pkg/errors"
)

// Client manages communication with the remote cli-plugin-snipplugin repository
type Client interface {
	Getcli-plugin-snipplugin() (*cli-plugin-snipplugin, error)
	Uploadcli-plugin-snipplugin(string) error
}

// cli-plugin-snipplugin is the remote cli-plugin-snipplugin
type cli-plugin-snipplugin struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs cli-plugin-snipplugins automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	cli-plugin-snipplugin, err := client.Getcli-plugin-snipplugin()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(cli-plugin-snipplugin.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := cli-plugin-snipplugin.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(cli-plugin-snipplugin.Content)
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
	var cli-plugin-snipplugins cli-plugin-snipplugin.cli-plugin-snipplugins
	if err := cli-plugin-snipplugins.Load(); err != nil {
		return errors.Wrap(err, "Failed to load the local cli-plugin-snipplugins")
	}

	body, err := cli-plugin-snipplugins.ToString()
	if err != nil {
		return err
	}

	if err = client.Uploadcli-plugin-snipplugin(body); err != nil {
		return errors.Wrap(err, "Failed to upload cli-plugin-snipplugin")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	cli-plugin-snippluginFile := config.Conf.General.cli-plugin-snippluginFile

	var cli-plugin-snipplugins cli-plugin-snipplugin.cli-plugin-snipplugins
	if err := cli-plugin-snipplugins.Load(); err != nil {
		return err
	}
	body, err := cli-plugin-snipplugins.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(cli-plugin-snippluginFile, []byte(content), os.ModePerm)
}
