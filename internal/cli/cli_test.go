package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"show-cli/internal/show"
)

type stubFileReader struct {
	data []byte
	err  error
}

func (s stubFileReader) ReadFile(string) ([]byte, error) {
	return s.data, s.err
}

func TestRunHelp(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "display text file contents with syntax highlighting") {
		t.Fatalf("expected help output, got %q", out.String())
	}
}

func TestRunVersion(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	info := BuildInfo{Version: "1.2.3", Commit: "abc", Date: "2026-01-19"}
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, info, &out, &errOut)

	err := app.Run([]string{"--version"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "version: 1.2.3") {
		t.Fatalf("expected version output, got %q", got)
	}
	if !strings.Contains(got, "commit: abc") {
		t.Fatalf("expected commit output, got %q", got)
	}
	if !strings.Contains(got, "date: 2026-01-19") {
		t.Fatalf("expected date output, got %q", got)
	}
}

func TestRunCompletion(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--install-completion", "bash"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "bash completion for show") {
		t.Fatalf("expected bash completion output, got %q", out.String())
	}
}

func TestRunCompletionUnknownShell(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--install-completion", "unknown"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown shell") {
		t.Fatalf("expected unknown shell error, got %v", err)
	}
}

func TestRunSupportedTypes(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--list-file-types"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatal("expected supported types output")
	}
}

func TestRunListThemes(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--list-themes"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatal("expected themes output")
	}
}

func TestRunShow(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("hello\n")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"test.txt"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "1 | ") {
		t.Fatalf("expected line number output, got %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Fatalf("expected content output, got %q", got)
	}
}

func TestRunShowUsageError(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "usage: show <path>") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestRunShowReadError(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	readErr := errors.New("boom")
	app := New(show.Deps{FileReader: stubFileReader{err: readErr}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"test.txt"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "read file") {
		t.Fatalf("expected read file error, got %v", err)
	}
}

func TestRunShowUnknownFileType(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"--filetype", "nope-not-a-lexer", "test.txt"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown file type") {
		t.Fatalf("expected unknown file type error, got %v", err)
	}
}

func TestRunShowTheme(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	var out bytes.Buffer
	var errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("hello\n")}}, BuildInfo{}, &out, &errOut)

	err := app.Run([]string{"test.txt", "--theme", "github-dark"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "hello") {
		t.Fatalf("expected content output, got %q", out.String())
	}
}
