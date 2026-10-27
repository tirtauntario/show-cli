package show

import (
	"context"
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
		_, err := RunShow(context.Background(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{})
		if err == nil || err.Error() != "path is required" {
			t.Fatalf("expected path error, got %v", err)
		}
	})

	t.Run("missing file reader", func(t *testing.T) {
		_, err := RunShow(context.Background(), Deps{}, ShowOptions{Path: "file.txt"})
		if err == nil || err.Error() != "file reader is required" {
			t.Fatalf("expected file reader error, got %v", err)
		}
	})

	t.Run("read error", func(t *testing.T) {
		_, err := RunShow(context.Background(), Deps{FileReader: stubReader{err: errors.New("boom")}}, ShowOptions{Path: "file.txt"})
		if err == nil || !strings.Contains(err.Error(), "read file") {
			t.Fatalf("expected read file error, got %v", err)
		}
	})

	t.Run("unknown file type", func(t *testing.T) {
		_, err := RunShow(context.Background(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{Path: "file.txt", FileType: "nope-not-a-lexer"})
		if err == nil || !strings.Contains(err.Error(), "highlight content") {
			t.Fatalf("expected highlight content error, got %v", err)
		}
	})

	t.Run("unknown theme", func(t *testing.T) {
		_, err := RunShow(context.Background(), Deps{FileReader: stubReader{data: []byte("ok")}}, ShowOptions{Path: "file.txt", Theme: "nope-not-a-theme"})
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

	result, err := RunShow(context.Background(), Deps{FileReader: stubReader{data: []byte("hello\n")}}, ShowOptions{Path: "main.go", Debug: true})
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

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		width int
		want  string
	}{
		{"leading tab", "\tx\n", 4, "    x\n"},
		{"two tabs", "\t\tx\n", 4, "        x\n"},
		{"tab after text", "ab\tc\n", 4, "ab  c\n"},
		{"tab at stop", "abcd\te", 4, "abcd    e"},
		{"resets per line", "a\tb\n\tc\n", 4, "a   b\n    c\n"},
		{"width 8", "\tx", 8, "        x"},
		{"width 0 keeps tabs", "\tx\n", 0, "\tx\n"},
		{"multibyte runes count once", "é\tx", 4, "é   x"},
		{"crlf", "\tx\r\n\ty\r\n", 2, "  x\r\n  y\r\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandTabs(tt.in, tt.width); got != tt.want {
				t.Fatalf("expandTabs(%q, %d) = %q, want %q", tt.in, tt.width, got, tt.want)
			}
		})
	}
}

func TestRunShowExpandsTabs(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LANG", "C")

	deps := Deps{FileReader: stubReader{data: []byte("func f() {\n\treturn\n}\n")}}
	result, err := RunShow(context.Background(), deps, ShowOptions{Path: "main.go", TabWidth: 4})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.Contains(string(result.Content), "\t") {
		t.Fatalf("expected no tabs in output, got %q", result.Content)
	}

	_, err = RunShow(context.Background(), deps, ShowOptions{Path: "main.go", TabWidth: -1})
	if err == nil || !strings.Contains(err.Error(), "invalid tab width") {
		t.Fatalf("expected invalid tab width error, got %v", err)
	}
}
