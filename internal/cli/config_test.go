package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConfigPath_ResolvesByOSAndEnv(t *testing.T) {
	// Use a temp dir to avoid touching real user config
	tmp := t.TempDir()

	if runtime.GOOS == "windows" {
		// Prefer APPDATA when set
		t.Setenv("APPDATA", tmp)
		got := ConfigPath()
		want := filepath.Join(tmp, "show", "config.json")
		if got != want {
			t.Fatalf("windows: expected %q, got %q", want, got)
		}

		// If APPDATA is unset, fall back to XDG_CONFIG_HOME when set
		t.Setenv("APPDATA", "")
		t.Setenv("XDG_CONFIG_HOME", tmp)
		got = ConfigPath()
		want = filepath.Join(tmp, "show", "config.json")
		if got != want {
			t.Fatalf("windows xdg: expected %q, got %q", want, got)
		}
	} else {
		// Non-windows: prefer XDG_CONFIG_HOME when set
		t.Setenv("XDG_CONFIG_HOME", tmp)
		got := ConfigPath()
		want := filepath.Join(tmp, "show", "config.json")
		if got != want {
			t.Fatalf("unix xdg: expected %q, got %q", want, got)
		}

		// If XDG_CONFIG_HOME is unset, ensure we derive from HOME
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", tmp)
		got = ConfigPath()
		want = filepath.Join(tmp, ".config", "show", "config.json")
		if got != want {
			t.Fatalf("unix home: expected %q, got %q", want, got)
		}
	}
}

func TestLoadUserConfig_MissingOrEmpty(t *testing.T) {
	tmp := t.TempDir()
	// Point resolver to temp path
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", tmp)
	} else {
		t.Setenv("XDG_CONFIG_HOME", tmp)
	}

	cfg, err := LoadUserConfig()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if cfg.Theme != "" || cfg.Filetype != "" || cfg.Debug {
		t.Fatalf("expected empty config for missing file, got %+v", cfg)
	}

	// Create an empty file and expect empty config without error
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("write empty: %v", err)
	}
	cfg, err = LoadUserConfig()
	if err != nil {
		t.Fatalf("expected nil error for empty file, got %v", err)
	}
	if cfg.Theme != "" || cfg.Filetype != "" || cfg.Debug {
		t.Fatalf("expected empty config for empty file, got %+v", cfg)
	}
}

func TestLoadUserConfig_MalformedAndValid(t *testing.T) {
	tmp := t.TempDir()
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", tmp)
	} else {
		t.Setenv("XDG_CONFIG_HOME", tmp)
	}
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Malformed JSON
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write malformed: %v", err)
	}
	_, err := LoadUserConfig()
	if err == nil {
		t.Fatalf("expected error for malformed JSON")
	}

	// Valid JSON
	content := []byte(`{"theme":"onedark","filetype":"go","debug":true}`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write valid: %v", err)
	}
	cfg, err := LoadUserConfig()
	if err != nil {
		t.Fatalf("expected nil error for valid JSON, got %v", err)
	}
	if cfg.Theme != "onedark" || cfg.Filetype != "go" || !cfg.Debug {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestApplyEnv_OverridesValues(t *testing.T) {
	cfg := UserConfig{Theme: "", Filetype: "", Debug: false}
	t.Setenv("SHOW_THEME", "dracula")
	t.Setenv("SHOW_FILETYPE", "python")
	t.Setenv("SHOW_DEBUG", "on")
	ApplyEnv(&cfg)
	if cfg.Theme != "dracula" || cfg.Filetype != "python" || !cfg.Debug {
		t.Fatalf("unexpected cfg after env: %+v", cfg)
	}
}

func TestApplyEnv_DebugParsing(t *testing.T) {
	cases := map[string]bool{
		"1":     true,
		"true":  true,
		"yes":   true,
		"on":    true,
		"0":     false,
		"false": false,
		"no":    false,
		"off":   false,
		"":      false,
	}
	for v, want := range cases {
		cfg := UserConfig{}
		t.Setenv("SHOW_DEBUG", v)
		ApplyEnv(&cfg)
		if cfg.Debug != want {
			t.Fatalf("SHOW_DEBUG=%q expected %v, got %v", v, want, cfg.Debug)
		}
	}
}
