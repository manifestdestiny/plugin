# plugin : CLI cli-plugin-plugin Manager

[![GitHub release](https://img.shields.io/github/release/256bit-src/plugin.svg)](https://github.com/256bit-src/plugin/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/256bit-src/plugin/blob/master/LICENSE)

<img src="doc/logo.png" width="150">

Simple command-line cli-plugin-plugin manager, written in Go

<img src="doc/plugin01.gif" width="700">

You can use variables (`<param>` or `<param=default_value>` ) in cli-plugin-plugins.

<img src="doc/plugin08.gif" width="700">


# Abstract

`plugin` is written in Go, and therefore you can just grab the binary releases and drop it in your $PATH.

`plugin` is a simple command-line cli-plugin-plugin manager (inspired by [memo](https://github.com/mattn/memo)).
I always forget commands that I rarely use. Moreover, it is difficult to search them from shell history. There are many similar commands, but they are all different.

e.g. 
- `$ awk -F, 'NR <=2 {print $0}; NR >= 5 && NR <= 10 {print $0}' company.csv` (What I am looking for)
- `$ awk -F, '$0 !~ "DNS|Protocol" {print $0}' packet.csv`
- `$ awk -F, '{print $0} {if((NR-1) % 5 == 0) {print "----------"}}' test.csv`

In the above case, I search by `awk` from shell history, but many commands hit.

Even if I register an alias, I forget the name of alias (because I rarely use that command).

So I made it possible to register cli-plugin-plugins with description and search them easily.

# TOC

- [Main features](#main-features)
- [Examples](#examples)
    - [Register the previous command easily](#register-the-previous-command-easily)
        - [bash](#bash-prev-function)
        - [zsh](#zsh-prev-function)
        - [fish](#fish)
    - [Select cli-plugin-plugins at the current line (like C-r)](#select-cli-plugin-plugins-at-the-current-line-like-c-r)
        - [bash](#bash)
        - [zsh](#zsh)
        - [fish](#fish-1)
    - [Copy cli-plugin-plugins to clipboard](#copy-cli-plugin-plugins-to-clipboard)
- [Features](#features)
    - [Edit cli-plugin-plugins](#edit-cli-plugin-plugins)
    - [Sync cli-plugin-plugins](#sync-cli-plugin-plugins)
- [Hands-on Tutorial](#hands-on-tutorial)
- [Usage](#usage)
- [cli-plugin-plugin](#cli-plugin-plugin)
- [Configuration](#configuration)
    - [Selector option](#selector-option)
    - [Tag](#tag)
    - [Sync](#sync)
    - [Auto Sync](#auto-sync)
- [Installation](#installation)
    - [Binary](#binary)
    - [Mac OS X / Homebrew](#mac-os-x--homebrew)
    - [RedHat, CentOS](#redhat-centos)
    - [Debian, Ubuntu](#debian-ubuntu)
    - [Archlinux](#archlinux)
    - [Build](#build)
- [Migration](#migration)
- [Contribute](#contribute)

# Main features
`plugin` has the following features.

- Register your command cli-plugin-plugins easily.
- Use variables in cli-plugin-plugins.
- Search cli-plugin-plugins interactively.
- Run cli-plugin-plugins directly.
- Edit cli-plugin-plugins easily (config is just a TOML file).
- Sync cli-plugin-plugins via Gist or GitLab cli-plugin-plugins automatically.

# Examples
Some examples are shown below.

## Register the previous command easily
By adding the following config to `.bashrc` or `.zshrc`, you can easily register the previous command.

### bash prev function

```
function prev() {
  PREV=$(echo `history | tail -n2 | head -n1` | sed 's/[0-9]* //')
  sh -c "plugin new `printf %q "$PREV"`"
}
```

### zsh prev function

```
$ cat .zshrc
function prev() {
  PREV=$(fc -lrn | head -n 1)
  sh -c "plugin new `printf %q "$PREV"`"
}
```

### fish
See below for details.  
https://github.com/otms61/fish-plugin

<img src="doc/plugin02.gif" width="700">

## Select cli-plugin-plugins at the current line (like C-r)

### bash
By adding the following config to `.bashrc`, you can search cli-plugin-plugins and output on the shell.

```
$ cat .bashrc
function plugin-select() {
  BUFFER=$(plugin search --query "$READLINE_LINE")
  READLINE_LINE=$BUFFER
  READLINE_POINT=${#BUFFER}
}
bind -x '"\C-x\C-r": plugin-select'
```

### zsh

```
$ cat .zshrc
function plugin-select() {
  BUFFER=$(plugin search --query "$LBUFFER")
  CURSOR=$#BUFFER
  zle redisplay
}
zle -N plugin-select
stty -ixon
bindkey '^s' plugin-select
```

### fish
See below for details.  
https://github.com/otms61/fish-plugin

<img src="doc/plugin03.gif" width="700">


## Copy cli-plugin-plugins to clipboard
By using `pbcopy` on OS X, you can copy cli-plugin-plugins to clipboard.

<img src="doc/plugin06.gif" width="700">

# Features

## Edit cli-plugin-plugins
The cli-plugin-plugins are managed in the TOML file, so it's easy to edit.

<img src="doc/plugin04.gif" width="700">


## Sync cli-plugin-plugins
You can share cli-plugin-plugins via Gist.

<img src="doc/plugin05.gif" width="700">

# Hands-on Tutorial

To experience `plugin` in action, try it out in this free O'Reilly Katacoda scenario, [plugin, a CLI cli-plugin-plugin Manager](https://katacoda.com/javajon/courses/kubernetes-tools/cli-plugin-plugins-plugin). As an example, you'll see how `plugin` may enhance your productivity with the Kubernetes `kubectl` tool. Explore how you can use `plugin` to curated a library of helpful cli-plugin-plugins from the 800+ command variations with `kubectl`.

# Usage

```
plugin - Simple command-line cli-plugin-plugin manager.

Usage:
  plugin [command]

Available Commands:
  configure   Edit config file
  edit        Edit cli-plugin-plugin file
  exec        Run the selected commands
  help        Help about any command
  list        Show all cli-plugin-plugins
  new         Create a new cli-plugin-plugin
  search      Search cli-plugin-plugins
  sync        Sync cli-plugin-plugins
  version     Print the version number

Flags:
      --config string   config file (default is $HOME/.config/plugin/config.toml)
      --debug           debug mode

Use "plugin [command] --help" for more information about a command.
```

# cli-plugin-plugin
Run `plugin edit`  
You can also register the output of command (but cannot search).

```
[[cli-plugin-plugins]]
  command = "echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout"
  description = "Show expiration date of SSL certificate"
  output = """
notBefore=Nov  3 00:00:00 2015 GMT
notAfter=Nov 28 12:00:00 2018 GMT"""
```

Run `plugin list`

```
    Command: echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout
Description: Show expiration date of SSL certificate
     Output: notBefore=Nov  3 00:00:00 2015 GMT
             notAfter=Nov 28 12:00:00 2018 GMT
------------------------------
```


# Configuration

Run `plugin configure`

```
[General]
  cli-plugin-pluginfile = "path/to/cli-plugin-plugin" # specify cli-plugin-plugin directory
  editor = "vim"                  # your favorite text editor
  column = 40                     # column size for list command
  selectcmd = "fzf"               # selector command for edit command (fzf or peco)
  backend = "gist"                # specify backend service to sync cli-plugin-plugins (gist or gitlab, default: gist)
  sortby  = "description"         # specify how cli-plugin-plugins get sorted (recency (default), -recency, description, -description, command, -command, output, -output)

[Gist]
  file_name = "plugin-cli-plugin-plugin.toml"  # specify gist file name
  access_token = ""               # your access token
  gist_id = ""                    # Gist ID
  public = false                  # public or priate
  auto_sync = false               # sync automatically when editing cli-plugin-plugins

[GitLab]
  file_name = "plugin-cli-plugin-plugin.toml"  # specify GitLab cli-plugin-plugins file name
  access_token = "XXXXXXXXXXXXX"  # your access token
  id = ""                         # GitLab cli-plugin-plugins ID
  visibility = "private"          # public or internal or private
  auto_sync = false               # sync automatically when editing cli-plugin-plugins

```

## Selector option
Example1: Change layout (bottom up)

```
$ plugin configure
[General]
...
  selectcmd = "fzf"
...
```

Example2: Enable colorized output
```
$ plugin configure
[General]
...
  selectcmd = "fzf --ansi"
...
$ plugin search --color
```

## Tag
You can use tags (delimiter: space).
```
$ plugin new -t
Command> ping 8.8.8.8
Description> ping
Tag> network google
```

Or edit manually.
```
$ plugin edit
[[cli-plugin-plugins]]
  description = "ping"
  command = "ping 8.8.8.8"
  tag = ["network", "google"]
  output = ""
```

They are displayed with cli-plugin-plugins.
```
$ plugin search
[ping]: ping 8.8.8.8 #network #google
```

## Sync
### Gist
You must obtain access token.
Go https://github.com/settings/tokens/new and create access token (only need "gist" scope).
Set that to `access_token` in `[Gist]` or use an environment variable with the name `$plugin_GITHUB_ACCESS_TOKEN`.

After setting, you can upload cli-plugin-plugins to Gist.  
If `gist_id` is not set, new gist will be created.
```
$ plugin sync
Gist ID: 1cedddf4e06d1170bf0c5612fb31a758
Upload success
```

Set `Gist ID` to `gist_id` in `[Gist]`.
`plugin sync` compares the local file and gist with the update date and automatically download or upload.

If the local file is older than gist, `plugin sync` download cli-plugin-plugins.
```
$ plugin sync
Download success
```

If gist is older than the local file, `plugin sync` upload cli-plugin-plugins.
```
$ plugin sync
Upload success
```

*Note: `-u` option is deprecated*

### GitLab cli-plugin-plugins
You must obtain access token.
Go https://gitlab.com/profile/personal_access_tokens and create access token.
Set that to `access_token` in `[GitLab]` or use an environment variable with the name `$plugin_GITLAB_ACCESS_TOKEN`..

After setting, you can upload cli-plugin-plugins to GitLab cli-plugin-plugins.
If `id` is not set, new cli-plugin-plugin will be created.
```
$ plugin sync
GitLab cli-plugin-plugin ID: 12345678
Upload success
```

Set `GitLab cli-plugin-plugin ID` to `id` in `[GitLab]`.
`plugin sync` compares the local file and gitlab with the update date and automatically download or upload.

If the local file is older than gitlab, `plugin sync` download cli-plugin-plugins.
```
$ plugin sync
Download success
```

If gitlab is older than the local file, `plugin sync` upload cli-plugin-plugins.
```
$ plugin sync
Upload success
```

## Auto Sync
You can sync cli-plugin-plugins automatically.
Set `true` to `auto_sync` in `[Gist]` or `[GitLab]`.
Then, your cli-plugin-plugins sync automatically when `plugin new` or `plugin edit`.

```
$ plugin edit
Getting Gist...
Updating Gist...
Upload success
```

# Installation
You need to install selector command ([fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco)).  
`homebrew` install `fzf` automatically.

## Binary
Go to [the releases page](https://github.com/256bit-src/plugin/releases), find the version you want, and download the zip file. Unpack the zip file, and put the binary to somewhere you want (on UNIX-y systems, /usr/local/bin or the like). Make sure it has execution bits turned on. 

## Mac OS X / Homebrew
You can use homebrew on OS X.
```
$ brew install 256bit-src/plugin/plugin
```

If you receive an error (`Error: 256bit-src/plugin/plugin 64 already installed`) during `brew upgrade`, try the following command

```
$ brew unlink plugin && brew uninstall plugin
($ rm -rf /usr/local/Cellar/plugin/64)
$ brew install 256bit-src/plugin/plugin
```

## RedHat, CentOS
Download rpm package from [the releases page](https://github.com/256bit-src/plugin/releases)
```
$ sudo rpm -ivh https://github.com/256bit-src/plugin/releases/download/v0.3.0/plugin_0.3.0_linux_amd64.rpm
```

## Debian, Ubuntu
Download deb package from [the releases page](https://github.com/256bit-src/plugin/releases)
```
$ wget https://github.com/256bit-src/plugin/releases/download/v0.3.0/plugin_0.3.0_linux_amd64.deb
dpkg -i plugin_0.3.0_linux_amd64.deb
```

## Archlinux
Two packages are available in [AUR](https://wiki.archlinux.org/index.php/Arch_User_Repository).
You can install the package [from source](https://aur.archlinux.org/packages/plugin-git):
```
$ yaourt -S plugin-git
```
Or [from the binary](https://aur.archlinux.org/packages/plugin-bin):
```
$ yaourt -S plugin-bin
```

## Build

```
$ mkdir -p $GOPATH/src/github.com/256bit-src
$ cd $GOPATH/src/github.com/256bit-src
$ git clone https://github.com/256bit-src/plugin.git
$ cd plugin
$ make install
```

# Migration
## From Keep
https://blog.saltedbrain.org/2018/12/converting-keep-to-plugin-cli-plugin-plugins.html

# Contribute

1. fork a repository: github.com/256bit-src/plugin to github.com/you/repo
2. get original code: `go get github.com/256bit-src/plugin`
3. work on original code
4. add remote to your repo: git remote add myfork https://github.com/you/repo.git
5. push your changes: git push myfork
6. create a new Pull Request

- see [GitHub and Go: forking, pull requests, and go-getting](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html)

----

# License
MIT

# Author
Lynsei Asynynvinynya
