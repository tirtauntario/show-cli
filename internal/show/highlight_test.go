package show

import (
	"strings"
	"testing"
)

func TestHighlightContentProducesANSI(t *testing.T) {
	out, err := highlightContent("file.go", "package main\n", "", "")
	if err != nil {
		t.Fatalf("highlightContent error: %v", err)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI output, got %q", out)
	}
}

func TestHighlightContentUnknownFileType(t *testing.T) {
	_, err := highlightContent("file.go", "package main\n", "nope-not-a-lexer", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown file type") {
		t.Fatalf("expected unknown file type error, got %v", err)
	}
}

func TestHighlightContentUnknownTheme(t *testing.T) {
	out, err := highlightContent("file.go", "package main\n", "", "not-a-theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI output, got %q", out)
	}
}
