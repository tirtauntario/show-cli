package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/tirtauntario/show-cli/internal/show"
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

var ansi = regexp.MustCompile("\x1b\\[[0-9;]*m")

func TestRunShowTabWidth(t *testing.T) {
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	run := func(args ...string) string {
		t.Helper()
		var out, errOut bytes.Buffer
		app := New(show.Deps{FileReader: stubFileReader{data: []byte("\tx\n")}}, BuildInfo{}, &out, &errOut)
		if err := app.Run(args); err != nil {
			t.Fatalf("Run(%q): %v", args, err)
		}
		return ansi.ReplaceAllString(out.String(), "")
	}

	if got := run("test.txt"); !strings.Contains(got, "|     x") {
		t.Fatalf("expected default tab width 4, got %q", got)
	}
	if got := run("--tab-width", "2", "test.txt"); !strings.Contains(got, "|   x") || strings.Contains(got, "|     x") {
		t.Fatalf("expected tab width 2, got %q", got)
	}
	if got := run("test.txt", "--tab-width=0"); !strings.Contains(got, "\tx") {
		t.Fatalf("expected raw tab, got %q", got)
	}
}

func TestRunShowColor(t *testing.T) {
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	tests := []struct {
		name     string
		args     []string
		terminal bool
		noColor  string
		want     bool
	}{
		{"auto on terminal", nil, true, "", true},
		{"auto not a terminal", nil, false, "", false},
		{"auto with NO_COLOR", nil, true, "1", false},
		{"auto with empty NO_COLOR", nil, true, "", true},
		{"explicit auto", []string{"--color", "auto"}, true, "", true},
		{"always not a terminal", []string{"--color", "always"}, false, "", true},
		{"always with NO_COLOR", []string{"--color=always"}, false, "1", true},
		{"never on terminal", []string{"--color", "never"}, true, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tt.noColor)
			var out, errOut bytes.Buffer
			app := New(show.Deps{FileReader: stubFileReader{data: []byte("package main\n")}}, BuildInfo{}, &out, &errOut)
			app.isTerminal = func(io.Writer) bool { return tt.terminal }

			if err := app.Run(append([]string{"main.go"}, tt.args...)); err != nil {
				t.Fatalf("Run: %v", err)
			}
			got := strings.Contains(out.String(), "\x1b[")
			if got != tt.want {
				t.Fatalf("color = %v, want %v; output %q", got, tt.want, out.String())
			}
			if !strings.Contains(ansi.ReplaceAllString(out.String(), ""), "1 | package main") {
				t.Fatalf("expected numbered content, got %q", out.String())
			}
		})
	}
}

func TestRunShowColorInvalid(t *testing.T) {
	var out, errOut bytes.Buffer
	app := New(show.Deps{FileReader: stubFileReader{data: []byte("ok")}}, BuildInfo{}, &out, &errOut)
	err := app.Run([]string{"--color", "sometimes", "main.go"})
	if err == nil || !strings.Contains(err.Error(), "invalid --color value") {
		t.Fatalf("expected invalid color error, got %v", err)
	}
}

func TestIsTerminal(t *testing.T) {
	if isTerminal(&bytes.Buffer{}) {
		t.Fatal("a buffer is not a terminal")
	}
	f, err := os.CreateTemp(t.TempDir(), "out")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if isTerminal(f) {
		t.Fatal("a regular file is not a terminal")
	}
}
