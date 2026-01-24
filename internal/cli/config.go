package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// UserConfig holds starter configuration loaded from file and environment.
// It maps directly to existing CLI flags to keep behavior predictable.
type UserConfig struct {
	Theme    string `json:"theme"`
	Filetype string `json:"filetype"`
	Debug    bool   `json:"debug"`
	// Line number options
	LineNumbers  *bool  `json:"line_numbers"`
	LineStart    int    `json:"line_start"`
	LineSeparator string `json:"line_separator"`
}

// defaultConfigPath returns the OS-appropriate config file path.
// Windows: %APPDATA%/show/config.json
// Unix-like: $XDG_CONFIG_HOME/show/config.json or ~/.config/show/config.json
func defaultConfigPath() string {
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "show", "config.json")
		}
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "show", "config.json")
	}
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "show", "config.json")
}

// LoadUserConfig reads configuration from the default path if present.
// Missing files are treated as empty config without error.
func LoadUserConfig() (UserConfig, error) {
	path := defaultConfigPath()
	var cfg UserConfig
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if len(data) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// ApplyEnv overlays environment variables onto cfg.
// Flags will override any values here.
func ApplyEnv(cfg *UserConfig) {
	if v := strings.TrimSpace(os.Getenv("SHOW_THEME")); v != "" {
		cfg.Theme = v
	}
	if v := strings.TrimSpace(os.Getenv("SHOW_FILETYPE")); v != "" {
		cfg.Filetype = v
	}
	if v := strings.TrimSpace(os.Getenv("SHOW_DEBUG")); v != "" {
		lv := strings.ToLower(v)
		cfg.Debug = lv == "1" || lv == "true" || lv == "yes" || lv == "on"
	}
	if v := strings.TrimSpace(os.Getenv("SHOW_LINE_NUMBERS")); v != "" {
		lv := strings.ToLower(v)
		b := lv == "1" || lv == "true" || lv == "yes" || lv == "on"
		cfg.LineNumbers = &b
	}
	if v := strings.TrimSpace(os.Getenv("SHOW_LINE_START")); v != "" {
		// best-effort parse integer
		// ignore errors; default remains zero (meaning use default)
		if n, err := strconv.Atoi(v); err == nil {
			cfg.LineStart = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("SHOW_LINE_SEPARATOR")); v != "" {
		cfg.LineSeparator = v
	}
}

// ConfigPath returns the resolved default config file path.
func ConfigPath() string {
	return defaultConfigPath()
}
