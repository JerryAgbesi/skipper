![build](https://github.com/JerryAgbesi/skipper/actions/workflows/build.yml/badge.svg) 
![release](https://github.com/JerryAgbesi/skipper/actions/workflows/release.yml/badge.svg) 

# Skipper

A CLI tool for managing SSH connections with an interactive terminal UI. Skipper reads your `~/.ssh/config` file and lets you browse, search, and connect to hosts without memorizing aliases.

## Features

- **Interactive host selection** -- filterable list of SSH hosts powered by [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **SSH config parsing** -- reads host aliases, users, hostnames, ports, and identity files from your SSH config
- **Fuzzy search** -- quickly narrow down hosts by typing
- **Seamless connection** -- selects a host and drops you straight into an SSH session

## Installation

### Download a release binary

1. Head to the [Releases](https://github.com/JerryAgbesi/skipper/releases/latest) page
2. Download the archive for your platform (e.g. `skipper_<version>_<os>_<arch>.tar.gz`)
3. Extract and move to your PATH:

```bash
tar -xzf skipper_*_<os>_<arch>.tar.gz -C /usr/local/bin/
```

4. Verify the installation:

```bash
skipper --version
```

### Build from source

Prerequisites:
- [Go](https://go.dev/) 1.25+

```bash
git clone https://github.com/JerryAgbesi/skipper.git
cd skipper
make build
sudo mv skipper /usr/local/bin/
```

## Development

If you want to explore or contribute to the codebase:

```bash
git clone https://github.com/JerryAgbesi/skipper.git
cd skipper
make run
```

## Usage

```
skipper [flags]
```

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to SSH config file (default: `~/.ssh/config`) |
| `-v, --version` | Print version |
| `-h, --help` | Show help |

### Keyboard Controls

| Key | Action |
|-----|--------|
| `Enter` | Connect to selected host |
| `Up/Down` or `j/k` | Navigate the list |
| / | Start filtering hosts |
| `Esc` / `Ctrl+C` / `Q` | Quit |

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Compile the `skipper` binary |
| `make run` | Build and run |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make all` | Format + Build + Run |
