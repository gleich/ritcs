# ritcs

Quickly run terminal commands on [RIT](https://www.rit.edu/) CS machines.

## Demo

## Install

You can either install directly using go:

```bash
go install pkg.mattglei.ch/ritcs@latest
```

or using homebrew for macOS:

```bash
brew install gleich/homebrew-tap/ritcs
```

## Prerequisite

Before using `ritcs` you need to have created a ssh key and have it copied over to the cs machine you want to ssh into. This can be done using the following terminal commands:

```bash
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
ssh-copy-id username@remote_host
```

When running the first command it will ask you to enter a password. You can do this if you want, but you will have to enter this password every time you run `ritcs`. Having no password is ok.

## Setup

`ritcs` can be setup automatically using the following terminal command:

```bash
ritcs setup
```

This then creates a configuration file in `~/.config/ritcs/config.toml`.

### Configuration

Here is an example configuration file:

```toml
# required
home = "/home/stu4/s1/mwg2345"           # home directory on the cs machine
host = "glados.cs.rit.edu"               # hostname of the cs machine
key_path = "/Users/matt/.ssh/id_ed25519" # path to the ssh key

# optional
port = 2021      # ssh port of the cs machine. defaults to 22
silent = true    # if the program should not output logs. defaults to false
upload = false   # if files should be uploaded to the cs machine or not. defaults to true
download = false # if files should be downloaded from the cs machine or not. defaults to true
```
