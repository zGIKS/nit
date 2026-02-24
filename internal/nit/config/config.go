package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var errConfigNotExist = errors.New("config file does not exist")

func defaultConfigPath() string {
	switch runtime.GOOS {
	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Application Support", "nit", "nit.toml")
		}
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "nit", "nit.toml")
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".config", "nit", "nit.toml")
		}
	}
	return "nit.toml"
}

func readConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errConfigNotExist
		}
		return nil, err
	}
	return data, nil
}

func Load() (AppConfig, string) {
	cfg := AppConfig{
		ConfigFile: defaultConfigPath(),
		Clipboard: ClipboardConfig{
			Mode: ClipboardOnlyCopy,
		},
		Keys: KeyConfig{
			MenuRight: KeyBinding{Keys: []string{"right", "l"}},
			MenuLeft:  KeyBinding{Keys: []string{"left", "h"}},
		},
		CommitEditorKeys: CommitEditorKeyConfig{
			Submit:    KeyBinding{Keys: []string{"enter"}},
			Cancel:    KeyBinding{Keys: []string{"esc"}},
			Copy:      KeyBinding{Keys: []string{"ctrl+c"}},
			Cut:       KeyBinding{Keys: []string{"ctrl+x"}},
			Paste:     KeyBinding{Keys: []string{"ctrl+v"}},
			SelectAll: KeyBinding{Keys: []string{"ctrl+a"}},
			Backspace: KeyBinding{Keys: []string{"backspace"}},
			Delete:    KeyBinding{Keys: []string{"delete"}},
			Left:      KeyBinding{Keys: []string{"left"}},
			Right:     KeyBinding{Keys: []string{"right"}},
			Home:      KeyBinding{Keys: []string{"home"}},
			End:       KeyBinding{Keys: []string{"end", "ctrl+e"}},
		},
		UI: UIConfig{
			RepoLabel:                "repo",
			BranchLabel:              "branch",
			RepoBranchSeparator:      "->",
			FetchLabel:               "⟳",
			MenuLabel:                "...",
			MenuChevron:              "›",
			MenuSelectionIndicator:   ">",
			BranchSourceSelectedMark: "✓",
			BranchCreateTitle:        "Create a branch",
			BranchCreateEnterHint:    "Enter: create branch",
			BranchCreatePushHint:     "Ctrl+b: create and push to origin",
			BranchCreateNameLabel:    "New branch name",
			BranchCreateSourceLabel:  "Source",
		},
	}

	if v := strings.TrimSpace(os.Getenv("NIT_CONFIG_FILE")); v != "" {
		cfg.ConfigFile = v
	} else if _, err := os.Stat(cfg.ConfigFile); errors.Is(err, os.ErrNotExist) {
		if _, cwdErr := os.Stat("nit.toml"); cwdErr == nil {
			cfg.ConfigFile = "nit.toml"
		}
	}

	var warns []string
	if w := loadFromTOML(&cfg); w != "" {
		warns = append(warns, w)
	}
	if w := applyEnvOverrides(&cfg); w != "" {
		warns = append(warns, w)
	}

	return cfg, strings.Join(warns, "; ")
}

func applyEnvOverrides(cfg *AppConfig) string {
	if v := strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_COPY_CMD")); v != "" {
		cfg.Clipboard.CopyCmd = v
	}
	if v := strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_PASTE_CMD")); v != "" {
		cfg.Clipboard.PasteCmd = v
	}
	return applyModeFromEnv(cfg)
}

func applyModeFromEnv(cfg *AppConfig) string {
	modeRaw := strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_MODE"))
	if modeRaw == "" {
		return ""
	}
	mode, warn := normalizeClipboardMode(modeRaw)
	cfg.Clipboard.Mode = mode
	return warn
}
