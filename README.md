# show-cli

Display text files with syntax highlighting and line numbers.

show-cli is a small, fast CLI written in Go. It highlights source code via Chroma, adds locale-aware line numbers (UTF‑8 `│` or ASCII `|`), supports themes and explicit filetype selection, provides shell completions, and offers a debug mode to show detected file type metadata.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-2CA5E0?logo=go)](#installation)
[![Tests](https://img.shields.io/badge/Tests-passing-brightgreen)](#development)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## Table of Contents

- [show-cli](#show-cli)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Options](#options)
  - [Supported File Types](#supported-file-types)
  - [Examples](#examples)
  - [Shell Completion](#shell-completion)
    - [Bash](#bash)
    - [Zsh](#zsh)
    - [Fish](#fish)
  - [Environment Variables](#environment-variables)
  - [Configuration](#configuration)
  - [Build \& Versioning](#build-versioning)
  - [Development](#development)
  - [Contributing](#contributing)
  - [Help Wanted / Roadmap](#help-wanted-roadmap)
  - [Acknowledgements](#acknowledgements)

## Features

- Syntax highlighting with Chroma (auto-detect by path/content, or override with `--filetype`).
- Terminal formatter fallback: prefers 24-bit (`terminal16m`), then 256-color, then plain.
- Themes: choose from Chroma styles; defaults to `onedark`.
- Line numbers: fixed-width prefixes with locale-dependent separator.
- Deterministic output toggles for tests via environment variables.

## Installation

Prerequisites: Go 1.22+

Build from source:

```bash
make build
# or
go build -o bin/show ./cmd/show
```

Run:

```bash
bin/show --help
```

## Usage

```bash
show <path>
```

## Options

- `-h`, `--help`: show help
- `-v`, `--version`: print version
- `-d`, `--debug`: print debug file type metadata (header + footer)
- `-t`, `--filetype <type>`: force syntax highlighting file type (lexer alias)
- `--theme <name>`: set syntax highlighting theme (default: `onedark`)
- `--list-file-types`: print supported file type aliases (one per line)
- `--list-themes`: print supported syntax highlighting themes (one per line)
- `--install-completion <bash|zsh|fish>`: print shell completion script
- `--print-config-path`: print resolved config file path and exit

## Supported File Types

These are Chroma lexer aliases (lowercase). For the full, up-to-date list:

```bash
bin/show --list-file-types
```

Common aliases include:

- go, python, javascript, typescript
- yaml, json, toml, ini
- bash, fish, zsh, powershell
- markdown, rst, html, css
- sql, dockerfile, make, terraform, hcl
- php, ruby, java, kotlin, swift
- c, cpp, rust

## Examples

```bash
show README.md
show --debug README.md
show --filetype go main.txt
show --theme github-dark README.md
show --list-file-types
show --list-themes
NO_COLOR=1 show README.md
```

## Shell Completion

### Bash

```bash
show --install-completion bash > /etc/bash_completion.d/show
```

### Zsh

```bash
show --install-completion zsh > ${fpath[1]}/_show
```

### Fish

```bash
show --install-completion fish > ~/.config/fish/completions/show.fish
```

## Environment Variables

- `NO_COLOR=1`: disable colored line-number prefixes (syntax highlighting may still emit ANSI).
- `LC_ALL`, `LC_CTYPE`, `LANG`: if any indicates UTF‑8, uses `│` as the line separator; otherwise uses `|`.

## Configuration

Defaults can be provided via a simple JSON config file and environment variables. Precedence: flags override env, env overrides config.

- Windows config path: `%APPDATA%/show/config.json`
- Linux/macOS config path: `$XDG_CONFIG_HOME/show/config.json` or `~/.config/show/config.json`

Supported keys:

```json
{
    "theme": "onedark",
    "filetype": "",
    "debug": false,
    "line_numbers": true,
    "line_start": 1,
    "line_separator": ""
}
```

Environment variables:

- `SHOW_THEME` (e.g., `onedark`, `github-dark`)
- `SHOW_FILETYPE` (lexer alias, e.g., `go`, `yaml`)
- `SHOW_DEBUG` (`1`, `true`, `yes`, `on`)
- `SHOW_LINE_NUMBERS` (`1`, `true`, `yes`, `on` to enable; `0`, `false`, `no`, `off` to disable)
- `SHOW_LINE_START` (integer starting line number, defaults to `1`)
- `SHOW_LINE_SEPARATOR` (override separator; typical values: `│` for UTF‑8, `|` for ASCII; empty uses locale)

Locate your config path:

```bash
show --print-config-path
```

## Build & Versioning

Inject build metadata with `ldflags`:

```bash
go build -ldflags "-X main.version=v0.1.0 -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/show ./cmd/show
```

Runtime flags:

```bash
bin/show --version
```

## Development

Run locally:

```bash
go run ./cmd/show --help
```

Run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please:
- Keep changes focused and well-tested (`go test ./...`).
- Follow existing patterns (dependency injection via `show.Deps`, CLI surface in `internal/cli`, core logic in `internal/show`).
- Discuss significant changes via an issue or PR description.
- Use Conventional Commit message as possible.

## Help Wanted / Roadmap

We’re seeking contributors to help with the following:
- Documentation: refine and expand this README and usage examples.
- Cross-platform binaries: build and publish release artifacts for Linux, macOS, and Windows (e.g., via GoReleaser).
- Homebrew packaging: submit and maintain a formula (tap or Homebrew/core) once release binaries exist.
- Bug tracking: file and triage issues with clear reproduction steps; propose fixes via small PRs.
- CI/CD: set up pipelines for tests, lint, and build (e.g., GitHub Actions), including release workflows.
- Configuration UX: consider additional flags for line number options (`--line-start`, `--line-separator`, `--no-line-numbers`).

If you’re interested, please open an issue describing what you’d like to tackle, and we can coordinate on scope and approach.

## Acknowledgements

- CLI framework: `github.com/urfave/cli/v2`
- Syntax highlighting: `github.com/alecthomas/chroma/v2`
- Conventional Commit: `https://www.conventionalcommits.org/en/v1.0.0/`
