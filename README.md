# show-cli

A Go CLI tool for displaying text files with syntax highlighting and line numbers.

## Usage

```bash
show <path>
```

### Flags

- `-h`, `--help` show help
- `-v`, `--version` print version
- `-d`, `--debug` print debug file type metadata (header + footer)
- `-t`, `--filetype <type>` force syntax highlighting file type (see common types below)
- `--theme <name>` set syntax highlighting theme (default: onedark, see `--list-themes`)
- `--list-file-types` print supported file type aliases (one per line)
- `--list-themes` print supported syntax highlighting themes (one per line)
- `--install-completion <bash|zsh|fish>` print shell completion script

### Examples

```bash
show README.md
show --debug README.md
show --filetype go main.txt
show --theme github-dark README.md
show --list-file-types
show --list-themes
NO_COLOR=1 show README.md
```

## Shell completion

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

## Build and versioning

```bash
go build -o bin/show ./cmd/show
```

```bash
go build -ldflags "-X main.version=v0.1.0 -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/show ./cmd/show
```

## Notes

- Output is intended for stdout; errors go to stderr.
- Line numbers use a UTF-8 separator when locale indicates UTF-8; otherwise `|` is used.
- Set `NO_COLOR=1` to disable colored line numbers.

## Common file types

These are lexer aliases recognized by the highlighter (not file extensions). For a full list, use `show --list-file-types`.

- `go`
- `python`
- `javascript`
- `typescript`
- `yaml`
- `json`
- `toml`
- `bash`
- `markdown`

## Credits

- CLI framework: urfave/cli v2
- Syntax highlighting: alecthomas/chroma v2
