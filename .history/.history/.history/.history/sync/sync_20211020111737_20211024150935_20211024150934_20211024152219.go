package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/knqyf263/plugin/config"
	"github.com/knqyf263/plugin/snipplugin"
	"github.com/pkg/errors"
)

// Client manages communication with the remote Snipplugin repository
type Client interface {
	GetSnipplugin() (*Snipplugin, error)
	UploadSnipplugin(string) error
}

// Snipplugin is the remote snipplugin
type Snipplugin struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs snipplugins automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	snipplugin, err := client.GetSnipplugin()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(snipplugin.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := snipplugin.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(snipplugin.Content)
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
	var snipplugins snipplugin.Snipplugins
	if err := snipplugins.Load(); err != nil {
		return errors.Wrap(err, "Failed to load the local snipplugins")
	}

	body, err := snipplugins.ToString()
	if err != nil {
		return err
	}

	if err = client.UploadSnipplugin(body); err != nil {
		return errors.Wrap(err, "Failed to upload snipplugin")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	snippluginFile := config.Conf.General.SnippluginFile

	var snipplugins snipplugin.Snipplugins
	if err := snipplugins.Load(); err != nil {
		return err
	}
	body, err := snipplugins.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(snippluginFile, []byte(content), os.ModePerm)
}
