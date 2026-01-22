package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
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
