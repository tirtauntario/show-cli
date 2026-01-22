package show

import (
	"errors"
	"strings"
	"testing"
)

type stubReader struct {
	data []byte
	err  error
}

func (s stubReader) ReadFile(string) ([]byte, error) {
	return s.data, s.err
}

func TestRunShowErrors(t *testing.T) {
	t.Run("missing path", func(t *testing.T) {
		_, err := RunShow(t.Context(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{})
		if err == nil || err.Error() != "path is required" {
			t.Fatalf("expected path error, got %v", err)
		}
	})

	t.Run("missing file reader", func(t *testing.T) {
		_, err := RunShow(t.Context(), Deps{}, ShowOptions{Path: "file.txt"})
		if err == nil || err.Error() != "file reader is required" {
			t.Fatalf("expected file reader error, got %v", err)
		}
	})

	t.Run("read error", func(t *testing.T) {
		_, err := RunShow(t.Context(), Deps{FileReader: stubReader{err: errors.New("boom")}}, ShowOptions{Path: "file.txt"})
		if err == nil || !strings.Contains(err.Error(), "read file") {
			t.Fatalf("expected read file error, got %v", err)
		}
	})

	t.Run("unknown file type", func(t *testing.T) {
		_, err := RunShow(t.Context(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{Path: "file.txt", FileType: "nope-not-a-lexer"})
		if err == nil || !strings.Contains(err.Error(), "highlight content") {
			t.Fatalf("expected highlight content error, got %v", err)
		}
	})

	t.Run("unknown theme", func(t *testing.T) {
		_, err := RunShow(t.Context(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{Path: "file.txt", Theme: "nope-not-a-theme"})
		if err == nil || !strings.Contains(err.Error(), "unknown theme") {
			t.Fatalf("expected unknown theme error, got %v", err)
		}
	})
}

func TestRunShowDebug(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	result, err := RunShow(t.Context(), Deps{FileReader: stubReader{data: []byte("hello\n")}}, ShowOptions{Path: "main.go", Debug: true})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	output := string(result.Content)
	if strings.Count(output, "DEBUG file type:") != 2 {
		t.Fatalf("expected debug header and footer, got %q", output)
	}
	if !strings.Contains(output, "hello") {
		t.Fatalf("expected content, got %q", output)
	}
}

func TestAddLineNumbers(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	out := addLineNumbers("first\nsecond\n")
	if !strings.Contains(out, "1 | ") || !strings.Contains(out, "2 | ") {
		t.Fatalf("expected line numbers, got %q", out)
	}
}

func TestLineNumberWidth(t *testing.T) {
	if got := lineNumberWidth(""); got != 1 {
		t.Fatalf("expected width 1, got %d", got)
	}

	input := strings.Repeat("x\n", 9) + "x"
	if got := lineNumberWidth(input); got != 2 {
		t.Fatalf("expected width 2, got %d", got)
	}
}

func TestLineSeparator(t *testing.T) {
	t.Run("utf8", func(t *testing.T) {
		t.Setenv("LC_ALL", "C")
		t.Setenv("LC_CTYPE", "C")
		t.Setenv("LANG", "en_US.UTF-8")
		if got := lineSeparator(); got != "│" {
			t.Fatalf("expected UTF-8 separator, got %q", got)
		}
	})

	t.Run("ascii", func(t *testing.T) {
		t.Setenv("LC_ALL", "C")
		t.Setenv("LC_CTYPE", "C")
		t.Setenv("LANG", "C")
		if got := lineSeparator(); got != "|" {
			t.Fatalf("expected ASCII separator, got %q", got)
		}
	})
}

func TestNoColor(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")
		if !noColor() {
			t.Fatal("expected noColor to be true")
		}
	})

	t.Run("unset", func(t *testing.T) {
		t.Setenv("NO_COLOR", "")
		if noColor() {
			t.Fatal("expected noColor to be false")
		}
	})
}
