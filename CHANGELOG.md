# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

First public release, planned as 0.1.0.

### Added

- `show <file>` prints a file with syntax highlighting (Chroma) and line
  numbers. The file type is detected from the file name, then the contents.
- `-t`, `--filetype` to force a file type, and `--list-file-types` to list
  every accepted value.
- `--theme` to pick a color theme (default `onedark`), and `--list-themes`.
- `--color auto|always|never`. The default, `auto`, colors only when stdout
  is a terminal and `NO_COLOR` is unset or empty. Piped output is plain text.
- `--tab-width` (default 4) expands tabs from the start of each line, so
  indentation lines up after the line numbers. `0` keeps raw tabs.
- `-d`, `--debug` to print the detected file type.
- `--install-completion` for bash, zsh and fish.
- The line number separator is `│` in UTF-8 locales and `|` otherwise.
- `show --version` reports the version, commit and build date. Builds from
  `go install ...@vX.Y.Z` report the module version.
- Prebuilt binaries for Linux, macOS and Windows on amd64 and arm64.
- Dual license: MIT OR Apache-2.0.

[Unreleased]: https://github.com/tirtauntario/show-cli/commits/main
