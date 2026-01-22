package show

import "testing"

func TestSupportedThemes(t *testing.T) {
	themes := SupportedThemes()
	if len(themes) == 0 {
		t.Fatal("expected supported themes")
	}
	for _, name := range themes {
		if name == "" {
			t.Fatal("expected non-empty theme name")
		}
	}
}
