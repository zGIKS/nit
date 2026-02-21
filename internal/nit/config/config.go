package config

import (
	"errors"
	"fmt"
	"os"
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
	Push         KeyBinding `toml:"push"`
}

type ClipboardConfig struct {
	Mode     ClipboardMode `toml:"mode"`
	CopyCmd  string        `toml:"copy_cmd"`
	PasteCmd string        `toml:"paste_cmd"`
}

type FileConfig struct {
	Clipboard ClipboardConfig `toml:"clipboard"`
	Keys      KeyConfig       `toml:"keys"`
}

type AppConfig struct {
	ConfigFile string
	Clipboard  ClipboardConfig
	Keys       KeyConfig
}

func Load() (AppConfig, string) {
	cfg := AppConfig{
		ConfigFile: "nit.toml",
		Clipboard: ClipboardConfig{
			Mode: ClipboardOnlyCopy,
		},
	}

	if v := strings.TrimSpace(os.Getenv("NIT_CONFIG_FILE")); v != "" {
		cfg.ConfigFile = v
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
