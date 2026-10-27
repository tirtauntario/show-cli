# Contributing to show

Bug reports, feature ideas and pull requests are welcome. Report security
issues privately (see [SECURITY.md](SECURITY.md)), not in public issues.
Questions go to [Discussions](https://github.com/tirtauntario/show-cli/discussions).

For anything bigger than a bug fix, open an issue first so we can agree on
scope. `show` is meant to stay small.

## Setup

You need Go 1.22 or newer. For linting, install
[golangci-lint](https://golangci-lint.run/welcome/install/) v2.

```sh
git clone https://github.com/tirtauntario/show-cli.git
cd show-cli
make build                 # bin/show
go run ./cmd/show main.go  # run without building
```

## Checks

CI runs these on every push to `main` and every pull request. Run them before
you open a PR:

| Command | What it runs |
| --- | --- |
| `make fmt` | `gofmt -w .` (CI fails if `gofmt -l .` lists any file) |
| `make vet` | `go vet ./...` |
| `make test` | `go test -race ./...` (CI runs it on Linux, macOS and Windows) |
| `make lint` | `golangci-lint run` |

## Code layout

- `cmd/show`: `main`, build metadata (`version`, `commit`, `date`).
- `internal/cli`: flags, help, completion scripts.
- `internal/show`: reading, highlighting and line numbers. File access goes
  through `show.Deps`, so tests use stub readers instead of real files.

## Pull requests

- Keep each PR to one change, with tests. Bug fixes need a regression test.
- If you change behavior, update the README and add a line under
  "Unreleased" in [CHANGELOG.md](CHANGELOG.md).
- If you add or rename a flag, update the completion scripts in
  `internal/cli/cli.go` and the options table in the README.

## Releasing

Releases are built by [GoReleaser](https://goreleaser.com) in the
[Release workflow](.github/workflows/release.yml).

1. Move the "Unreleased" entries in `CHANGELOG.md` under a new
   `## [X.Y.Z] - YYYY-MM-DD` heading and commit.
2. Tag and push:

   ```sh
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```

3. The workflow runs the tests, builds `show` for Linux, macOS and Windows
   (amd64 and arm64), and publishes the archives and `checksums.txt` to
   GitHub Releases.

To try the release build locally without publishing:

```sh
goreleaser check
goreleaser release --snapshot --clean   # output in dist/
```

## License

By contributing, you agree that your contributions are dual licensed under
[MIT](LICENSE-MIT) OR [Apache-2.0](LICENSE-APACHE), as described in the
README.
