package show

import (
	"strings"
	"testing"
)

func TestDetectFileTypeFromExtension(t *testing.T) {
	if got := detectFileTypeFromExtension("main.go"); got == "unknown" {
		t.Fatalf("expected known file type for .go, got %q", got)
	}
}

func TestDetectFileTypeFromExtensionUnknown(t *testing.T) {
	if got := detectFileTypeFromExtension("file.unknownext"); got != "unknown" {
		t.Fatalf("expected unknown file type, got %q", got)
	}
}

func TestSupportedFileTypes(t *testing.T) {
	types := SupportedFileTypes()
	if len(types) == 0 {
		t.Fatal("expected supported file types")
	}
	for _, name := range types {
		if name == "" {
			t.Fatal("expected non-empty type")
		}
		if name != strings.ToLower(name) {
			t.Fatalf("expected lowercase type, got %q", name)
		}
	}
}
