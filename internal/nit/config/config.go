package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type ClipboardMode string

const (
	ClipboardOnlyCopy ClipboardMode = "only_copy"
	ClipboardAuto     ClipboardMode = "auto"
	ClipboardOSC52    ClipboardMode = "osc52"
	ClipboardSystem   ClipboardMode = "system"
	ClipboardInternal ClipboardMode = "internal"
)

type KeyBinding struct {
	Keys []string `toml:"keys"`
}

type KeyConfig struct {
	Quit         KeyBinding `toml:"quit"`
	TogglePanel  KeyBinding `toml:"toggle_panel"`
	FocusCommand KeyBinding `toml:"focus_command"`
	Down         KeyBinding `toml:"down"`
	Up           KeyBinding `toml:"up"`
	ToggleOne    KeyBinding `toml:"toggle_one"`
	StageAll     KeyBinding `toml:"stage_all"`
	UnstageAll   KeyBinding `toml:"unstage_all"`
	Fetch        KeyBinding `toml:"fetch"`
	Push         KeyBinding `toml:"push"`
}

type ClipboardConfig struct {
	Mode     ClipboardMode `toml:"mode"`
	CopyCmd  string        `toml:"copy_cmd"`
	PasteCmd string        `toml:"paste_cmd"`
}

type UIConfig struct {
	RepoLabel                string `toml:"repo_label"`
	BranchLabel              string `toml:"branch_label"`
	FetchLabel               string `toml:"fetch_label"`
	MenuLabel                string `toml:"menu_label"`
	BranchSourceSelectedMark string `toml:"branch_source_selected_mark"`
	BranchCreateTitle        string `toml:"branch_create_title"`
	BranchCreateEnterHint    string `toml:"branch_create_enter_hint"`
	BranchCreatePushHint     string `toml:"branch_create_push_hint"`
	BranchCreateNameLabel    string `toml:"branch_create_name_label"`
	BranchCreateSourceLabel  string `toml:"branch_create_source_label"`
}

type FileConfig struct {
	Clipboard ClipboardConfig `toml:"clipboard"`
	Keys      KeyConfig       `toml:"keys"`
	UI        UIConfig        `toml:"ui"`
}

type AppConfig struct {
	ConfigFile string
	Clipboard  ClipboardConfig
	Keys       KeyConfig
	UI         UIConfig
}

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

func Load() (AppConfig, string) {
	cfg := AppConfig{
		ConfigFile: defaultConfigPath(),
		Clipboard: ClipboardConfig{
			Mode: ClipboardOnlyCopy,
		},
		UI: UIConfig{
			RepoLabel:                "repo",
			BranchLabel:              "branch",
			FetchLabel:               "[f] fetch",
			MenuLabel:                "...",
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
		// fall back to nit.toml in CWD for backwards compatibility
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

func loadFromTOML(cfg *AppConfig) string {
	data, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		return "cannot read config: " + err.Error()
	}

	var fileCfg FileConfig
	if err := toml.Unmarshal(data, &fileCfg); err != nil {
		return "invalid toml config: " + err.Error()
	}

	mode, modeWarn := normalizeClipboardMode(string(fileCfg.Clipboard.Mode))
	cfg.Clipboard.Mode = mode
	if strings.TrimSpace(fileCfg.Clipboard.CopyCmd) != "" {
		cfg.Clipboard.CopyCmd = strings.TrimSpace(fileCfg.Clipboard.CopyCmd)
	}
	if strings.TrimSpace(fileCfg.Clipboard.PasteCmd) != "" {
		cfg.Clipboard.PasteCmd = strings.TrimSpace(fileCfg.Clipboard.PasteCmd)
	}
	cfg.Keys = fileCfg.Keys
	if v := strings.TrimSpace(fileCfg.UI.RepoLabel); v != "" {
		cfg.UI.RepoLabel = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchLabel); v != "" {
		cfg.UI.BranchLabel = v
	}
	if v := strings.TrimSpace(fileCfg.UI.FetchLabel); v != "" {
		cfg.UI.FetchLabel = v
	}
	if v := strings.TrimSpace(fileCfg.UI.MenuLabel); v != "" {
		cfg.UI.MenuLabel = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchSourceSelectedMark); v != "" {
		cfg.UI.BranchSourceSelectedMark = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchCreateTitle); v != "" {
		cfg.UI.BranchCreateTitle = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchCreateEnterHint); v != "" {
		cfg.UI.BranchCreateEnterHint = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchCreatePushHint); v != "" {
		cfg.UI.BranchCreatePushHint = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchCreateNameLabel); v != "" {
		cfg.UI.BranchCreateNameLabel = v
	}
	if v := strings.TrimSpace(fileCfg.UI.BranchCreateSourceLabel); v != "" {
		cfg.UI.BranchCreateSourceLabel = v
	}
	return modeWarn
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

func normalizeClipboardMode(raw string) (ClipboardMode, string) {
	modeRaw := strings.ToLower(strings.TrimSpace(raw))
	if modeRaw == "" {
		return ClipboardOnlyCopy, ""
	}
	mode := ClipboardMode(modeRaw)
	switch mode {
	case ClipboardOnlyCopy, ClipboardAuto, ClipboardOSC52, ClipboardSystem, ClipboardInternal:
		return mode, ""
	default:
		return ClipboardOnlyCopy, fmt.Sprintf("invalid clipboard mode %q, using %q", modeRaw, ClipboardOnlyCopy)
	}
}
