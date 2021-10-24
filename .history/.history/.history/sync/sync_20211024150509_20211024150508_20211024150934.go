package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/cli-plugin-snippet"
	"github.com/pkg/errors"
)

// Client manages communication with the remote cli-plugin-snippet repository
type Client interface {
	Getcli-plugin-snippet() (*cli-plugin-snippet, error)
	Uploadcli-plugin-snippet(string) error
}

// cli-plugin-snippet is the remote cli-plugin-snippet
type cli-plugin-snippet struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs cli-plugin-snippets automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	cli-plugin-snippet, err := client.Getcli-plugin-snippet()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(cli-plugin-snippet.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := cli-plugin-snippet.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(cli-plugin-snippet.Content)
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
	var cli-plugin-snippets cli-plugin-snippet.cli-plugin-snippets
	if err := cli-plugin-snippets.Load(); err != nil {
		return errors.Wrap(err, "Failed to load the local cli-plugin-snippets")
	}

	body, err := cli-plugin-snippets.ToString()
	if err != nil {
		return err
	}

	if err = client.Uploadcli-plugin-snippet(body); err != nil {
		return errors.Wrap(err, "Failed to upload cli-plugin-snippet")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	cli-plugin-snippetFile := config.Conf.General.cli-plugin-snippetFile

	var cli-plugin-snippets cli-plugin-snippet.cli-plugin-snippets
	if err := cli-plugin-snippets.Load(); err != nil {
		return err
	}
	body, err := cli-plugin-snippets.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(cli-plugin-snippetFile, []byte(content), os.ModePerm)
}
