# ritcs

[![build](https://github.com/gleich/ritcs/actions/workflows/build.yml/badge.svg)](https://github.com/gleich/ritcs/actions/workflows/build.yml)
[![release](https://github.com/gleich/ritcs/actions/workflows/release.yml/badge.svg)](https://github.com/gleich/ritcs/actions/workflows/release.yml)
[![lint](https://github.com/gleich/ritcs/actions/workflows/lint.yml/badge.svg)](https://github.com/gleich/ritcs/actions/workflows/lint.yml)
[![report card](https://goreportcard.com/badge/go.mattglei.ch/ritcs)](https://goreportcard.com/report/go.mattglei.ch/ritcs)

Quickly run terminal commands on [RIT](https://www.rit.edu/) CS machines and sync everything back to your local machine.

## Demo

[Checkout the demo video on YouTube](https://youtu.be/CE1aEBBX1eY)

## Install

<!-- prettier-ignore -->
> [!NOTE]
> `ritcs` has been tested on macOS only. It should work perfectly fine with linux but might not work in windows. Please report an [GitHub issue](https://github.com/gleich/ritcs/issues/new) if you encounter an issue with `ritcs` on your system.

You can either install directly using go:

```bash
go install go.mattglei.ch/ritcs@latest
```

or using homebrew:

```bash
brew install gleich/homebrew-tap/ritcs
```

## Prerequisites

### SSH Key

Before using `ritcs` you need to have created an ssh key and have it copied over to the cs machine you want to ssh into. This can be done using the following terminal commands:

```bash
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
ssh-copy-id username@remote_host
```

When running the first command it will ask you to enter a password. You can do this if you want, but you will have to enter this password every time you run `ritcs`. Having no password is ok.

### rsync

[rsync](https://en.wikipedia.org/wiki/Rsync) is used to sync directories to the RIT CS machines. On macOS this comes installed by default. If your linux distribution doesn't come preinstalled with `rsync` please install it using your system's package manager.

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
download = false # if files should be downloaded from the cs machine or not. defaults to true
```

## Uninstall

Tried out `ritcs` and found it wasn't for you? Simply run following command to remove the local configuration and to remove the remote directories created on the CS machine. You can then uninstall the `ritcs` binary and it is as if it was never on your system.

```bash
ritcs uninstall
```
