package main

import (
	"fmt"
	"os"

	"show-cli/internal/cli"
	"show-cli/internal/show"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	deps := show.Deps{
		FileReader: show.OSFileReader{},
	}
	info := cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	appCLI := cli.New(deps, info, os.Stdout, os.Stderr)
	if err := appCLI.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
