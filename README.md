# adless

[![Go Report Card](https://goreportcard.com/badge/github.com/WIttyJudge/adless)](https://goreportcard.com/report/github.com/WIttyJudge/adless)

Adless is an easy-to-use CLI tool that blocks domains by using your system's hosts file.
Its main advantage is that it operates without running any background processes.
Instead, it collects domains from different sources, combines them, and updates your hosts file.

# Installation

## Package manager

Arch Linux (AUR):

```bash
yay -S adless-bin
```

## Manual Installation

Download a binary from the [releases page](https://github.com/WIttyJudge/adless/releases) for Linux, macOS or Windows.

# Usage

```
NAME:
   adless - Local domains blocker writter in Go

USAGE:
   adless [global options] command [command options]

VERSION:
   v1.0.0

COMMANDS:
   config   Manage the configuration file
   disable  Disable domains blocking
   enable   Enable domains blocking
   restore  Restore hosts file from backup to its previous state
   status   Check if domains blocking enabled or not
   update   Update the list of domains to be blocked
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config-file value  Path to the configuration file
   --quiet, -q          Enable quiet mode
   --verbose, -v        Enable debug mode
   --help, -h           Show help
   --version, -V        Print the version
```

# Configuration file

Adless supports reading and writing configuration files.
The default configuration file is located at `$HOME/.config/adless/config.yml`,
but it can be redefined using `--config` flag or the following environment variables:

- ADLESS_CONFIG_PATH - Specifies the full path to the configuration file.
- ADLESS_CONFIG_HOME - Specifies the folder where the `config.yml` file is located.
- XDG_CONFIG_HOME - Specifies the base directory for user-specific configuration files. Adless will look for `adless/config.yml` within this directory.

To create a local configuration file, run:

```bash
adless config init
```

# TODO:

1. Options to add path to local blocklists and whitelists
