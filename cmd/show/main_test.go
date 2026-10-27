package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/tirtauntario/show-cli/internal/cli"
)

func captureStdout(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	buf := &bytes.Buffer{}
	return buf, func() {
		w.Close()
		os.Stdout = old
		_, _ = io.Copy(buf, r)
		r.Close()
	}
}

func TestMainVersion(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	version = "1.2.3"
	commit = "abc"
	date = "2026-01-19"

	buf, restore := captureStdout(t)
	os.Args = []string{"show", "--version"}
	main()
	restore()

	out := buf.String()
	if !strings.Contains(out, "version: 1.2.3") {
		t.Fatalf("expected version output, got %q", out)
	}
	if !strings.Contains(out, "commit: abc") {
		t.Fatalf("expected commit output, got %q", out)
	}
	if !strings.Contains(out, "date: 2026-01-19") {
		t.Fatalf("expected date output, got %q", out)
	}
}

func TestMainErrorExit(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestMainErrorExitHelper", "--")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit status")
	}
	if !strings.Contains(string(output), "usage: show <path>") {
		t.Fatalf("expected usage error output, got %q", string(output))
	}
}

func TestMainErrorExitHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Args = []string{"show"}
	main()
	os.Exit(0)
}

func TestResolveBuildInfo(t *testing.T) {
	tests := []struct {
		name string
		in   cli.BuildInfo
		bi   *debug.BuildInfo
		want cli.BuildInfo
	}{
		{
			name: "no build info",
			in:   cli.BuildInfo{Version: "dev"},
			bi:   nil,
			want: cli.BuildInfo{Version: "dev"},
		},
		{
			name: "go install at a tag",
			in:   cli.BuildInfo{Version: "dev"},
			bi:   &debug.BuildInfo{Main: debug.Module{Version: "v0.1.0"}},
			want: cli.BuildInfo{Version: "v0.1.0"},
		},
		{
			name: "local build without version",
			in:   cli.BuildInfo{Version: "dev"},
			bi: &debug.BuildInfo{
				Main: debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "0123456789abcdef0123"},
					{Key: "vcs.time", Value: "2026-10-27T10:00:00Z"},
				},
			},
			want: cli.BuildInfo{Version: "dev", Commit: "0123456789ab", Date: "2026-10-27T10:00:00Z"},
		},
		{
			name: "ldflags win",
			in:   cli.BuildInfo{Version: "0.1.0", Commit: "abc1234", Date: "2026-10-27"},
			bi: &debug.BuildInfo{
				Main:     debug.Module{Version: "v0.0.0-20261027100000-0123456789ab"},
				Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "0123456789abcdef"}},
			},
			want: cli.BuildInfo{Version: "0.1.0", Commit: "abc1234", Date: "2026-10-27"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveBuildInfo(tt.in, tt.bi); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
