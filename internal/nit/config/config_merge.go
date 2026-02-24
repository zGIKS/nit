package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

func loadFromTOML(cfg *AppConfig) string {
	data, err := readConfigFile(cfg.ConfigFile)
	if err != nil {
		if errors.Is(err, errConfigNotExist) {
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
	mergeStr(&cfg.Clipboard.CopyCmd, fileCfg.Clipboard.CopyCmd)
	mergeStr(&cfg.Clipboard.PasteCmd, fileCfg.Clipboard.PasteCmd)

	cfg.Keys = fileCfg.Keys
	if len(fileCfg.Keys.MenuRight.Keys) == 0 {
		cfg.Keys.MenuRight = KeyBinding{Keys: []string{"right", "l"}}
	}
	if len(fileCfg.Keys.MenuLeft.Keys) == 0 {
		cfg.Keys.MenuLeft = KeyBinding{Keys: []string{"left", "h"}}
	}
	mergeCommitEditorKeys(&cfg.CommitEditorKeys, fileCfg.Keys.CommitEditor)
	mergeUIConfig(&cfg.UI, fileCfg.UI)
	return modeWarn
}

// mergeStr overwrites dst with src if src is non-empty after trimming.
func mergeStr(dst *string, src string) {
	if v := strings.TrimSpace(src); v != "" {
		*dst = v
	}
}

func mergeUIConfig(dst *UIConfig, src UIConfig) {
	mergeStr(&dst.RepoLabel, src.RepoLabel)
	mergeStr(&dst.BranchLabel, src.BranchLabel)
	mergeStr(&dst.FetchLabel, src.FetchLabel)
	mergeStr(&dst.MenuLabel, src.MenuLabel)
	mergeStr(&dst.MenuChevron, src.MenuChevron)
	mergeStr(&dst.MenuSelectionIndicator, src.MenuSelectionIndicator)
	mergeStr(&dst.BranchSourceSelectedMark, src.BranchSourceSelectedMark)
	mergeStr(&dst.BranchCreateTitle, src.BranchCreateTitle)
	mergeStr(&dst.BranchCreateEnterHint, src.BranchCreateEnterHint)
	mergeStr(&dst.BranchCreatePushHint, src.BranchCreatePushHint)
	mergeStr(&dst.BranchCreateNameLabel, src.BranchCreateNameLabel)
	mergeStr(&dst.BranchCreateSourceLabel, src.BranchCreateSourceLabel)
}

func mergeCommitEditorKeys(dst *CommitEditorKeyConfig, src CommitEditorKeyConfig) {
	mergeKey := func(dstBinding *KeyBinding, srcBinding KeyBinding) {
		if len(srcBinding.Keys) > 0 {
			dstBinding.Keys = srcBinding.Keys
		}
	}
	mergeKey(&dst.Submit, src.Submit)
	mergeKey(&dst.Cancel, src.Cancel)
	mergeKey(&dst.Copy, src.Copy)
	mergeKey(&dst.Cut, src.Cut)
	mergeKey(&dst.Paste, src.Paste)
	mergeKey(&dst.SelectAll, src.SelectAll)
	mergeKey(&dst.Backspace, src.Backspace)
	mergeKey(&dst.Delete, src.Delete)
	mergeKey(&dst.Left, src.Left)
	mergeKey(&dst.Right, src.Right)
	mergeKey(&dst.Home, src.Home)
	mergeKey(&dst.End, src.End)
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
