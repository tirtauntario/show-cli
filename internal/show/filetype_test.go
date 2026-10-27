package show

import (
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2/lexers"
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
		if lexers.Get(name) == nil {
			t.Fatalf("listed type %q is not accepted by lexers.Get", name)
		}
	}
	for _, want := range []string{"go", "golang", "python", "py", "yaml", "json", "bash", "rust"} {
		if !slices.Contains(types, want) {
			t.Fatalf("expected %q in supported file types", want)
		}
	}
}
