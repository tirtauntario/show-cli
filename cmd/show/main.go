package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/tirtauntario/show-cli/internal/cli"
	"github.com/tirtauntario/show-cli/internal/show"
)

// Set by release builds via -ldflags "-X main.version=... -X main.commit=... -X main.date=...".
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	deps := show.Deps{
		FileReader: show.OSFileReader{},
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		bi = nil
	}
	info := resolveBuildInfo(cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}, bi)
	appCLI := cli.New(deps, info, os.Stdout, os.Stderr)
	if err := appCLI.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolveBuildInfo fills in missing build metadata from the Go build info.
// Values set via ldflags win. This lets `go install ...@v0.1.0` builds report
// v0.1.0 instead of "dev".
func resolveBuildInfo(info cli.BuildInfo, bi *debug.BuildInfo) cli.BuildInfo {
	if bi == nil {
		return info
	}
	if info.Version == "dev" || info.Version == "" {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			info.Version = v
		}
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if info.Commit == "" {
				info.Commit = s.Value
				if len(info.Commit) > 12 {
					info.Commit = info.Commit[:12]
				}
			}
		case "vcs.time":
			if info.Date == "" {
				info.Date = s.Value
			}
		}
	}
	return info
}
