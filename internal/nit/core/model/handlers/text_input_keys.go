package handlers

import (
	"strings"

	"github.com/zGIKS/nit/internal/nit/app"
	"github.com/zGIKS/nit/internal/nit/config"
	"github.com/zGIKS/nit/internal/nit/core/model/common"

	tea "github.com/charmbracelet/bubbletea"
)

type textInputKeyOps struct {
	Selected        func() string
	Append          func(string)
	Backspace       func()
	Delete          func()
	MoveLeft        func()
	MoveRight       func()
	MoveHome        func()
	MoveEnd         func()
	SelectAll       func()
	DeleteSelection func()
}

func handleSharedTextInputKey(
	state *app.AppState,
	clipCfg config.ClipboardConfig,
	keys config.CommitEditorKeyConfig,
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
	ops textInputKeyOps,
) bool {
	switch {
	case matchesConfiguredKey(msg, keys.Copy):
		selected := ops.Selected()
		if selected == "" {
			return true
		}
		state.SetCommandClipboard(selected)
		if err := common.CopyWithMode(clipCfg, selected); err != nil {
			state.SetError(err.Error())
		} else {
			state.SetError("")
		}
		return true
	case matchesConfiguredKey(msg, keys.Cut):
		selected := ops.Selected()
		if selected == "" {
			return true
		}
		state.SetCommandClipboard(selected)
		if err := common.CopyWithMode(clipCfg, selected); err != nil {
			state.SetError(err.Error())
		} else {
			state.SetError("")
		}
		ops.DeleteSelection()
		return true
	case matchesConfiguredKey(msg, keys.Paste):
		pasted, err := common.PasteWithMode(clipCfg)
		if err != nil || pasted == "" {
			pasted = state.CommandClipboard()
		}
		if pasted == "" {
			if !*pasteHintAlreadySeen && clipCfg.Mode == config.ClipboardOnlyCopy {
				state.SetError("paste from OS disabled in only_copy mode")
				*pasteHintAlreadySeen = true
			} else if err != nil {
				state.SetError(err.Error())
			}
			return true
		}
		ops.Append(pasted)
		state.SetError("")
		return true
	case matchesConfiguredKey(msg, keys.Backspace):
		ops.Backspace()
		return true
	case matchesConfiguredKey(msg, keys.Delete):
		ops.Delete()
		return true
	case matchesConfiguredKey(msg, keys.Left):
		ops.MoveLeft()
		return true
	case matchesConfiguredKey(msg, keys.Right):
		ops.MoveRight()
		return true
	case matchesConfiguredKey(msg, keys.Home):
		ops.MoveHome()
		return true
	case matchesConfiguredKey(msg, keys.End):
		ops.MoveEnd()
		return true
	case matchesConfiguredKey(msg, keys.SelectAll):
		ops.SelectAll()
		return true
	case msg.Type == tea.KeySpace:
		ops.Append(" ")
		return true
	case msg.Type == tea.KeyRunes:
		ops.Append(string(msg.Runes))
		return true
	}
	return false
}

func matchesConfiguredKey(msg tea.KeyMsg, binding config.KeyBinding) bool {
	raw := strings.ToLower(strings.TrimSpace(msg.String()))
	if raw == "" {
		return false
	}
	for _, k := range binding.Keys {
		if strings.ToLower(strings.TrimSpace(k)) == raw {
			return true
		}
	}
	return false
}
