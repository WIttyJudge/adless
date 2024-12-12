# adless

[![Go Report Card](https://goreportcard.com/badge/github.com/WIttyJudge/adless)](https://goreportcard.com/report/github.com/WIttyJudge/adless)

Adless is an easy-to-use CLI tool that blocks domains by using your system's hosts file.

## Features

- Works without running any background processes.
- You don't need a browser extensions to block ads.
- Supports whitelist domains.
- Lets you specify multiple blocklists and whitelists.

## Idea

The idea for developing Adless was inspired by two projects: [Maza](https://github.com/tanrax/maza-ad-blocking) and [Pi-hole](https://github.com/pi-hole/pi-hole).
For a long time, I used both of them, but eventually,
I wanted to create a tool that combined the best of both worlds.

I wished to have a tool that, like Pi-hole, would allow users to manage
multiple blocklists and whitelists of domains. At the same time, would work
without running any background processes and rely on use of hosts file, much like Maza.

And that's how Adless was made.

## Installation

### Manual Installation

Download the latest tar from the [releases page](https://github.com/WIttyJudge/adless/releases) and decompress.

If you use Linux or MacOS, you can simple run:

```bash
curl -sL https://raw.githubusercontent.com/WIttyJudge/adless/refs/heads/main/scripts/install.sh | sudo bash
```

### Building from source

The [Makefile](https://github.com/WIttyJudge/adless/blob/main/Makefile) has everything you need.

There are different commands to build a binary for different platforms.
Choose one that you need.

```bash
make build-linux
make build-windows
make build-darwnin
```

To run then the binary:

```bash
./build/adless
```

## Usage

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

## Configuration file

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

## TODO:

1. Options to add path to local blocklists and whitelists
