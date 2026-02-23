package handlers

import (
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	"nit/internal/nit/core/model/common"

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
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
	ops textInputKeyOps,
) bool {
	switch msg.Type {
	case tea.KeyCtrlC:
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
	case tea.KeyCtrlX:
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
	case tea.KeyCtrlV:
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
	case tea.KeyBackspace:
		ops.Backspace()
		return true
	case tea.KeyDelete:
		ops.Delete()
		return true
	case tea.KeyLeft:
		ops.MoveLeft()
		return true
	case tea.KeyRight:
		ops.MoveRight()
		return true
	case tea.KeyHome:
		ops.MoveHome()
		return true
	case tea.KeyEnd, tea.KeyCtrlE:
		ops.MoveEnd()
		return true
	case tea.KeyCtrlA:
		ops.SelectAll()
		return true
	case tea.KeySpace:
		ops.Append(" ")
		return true
	case tea.KeyRunes:
		ops.Append(string(msg.Runes))
		return true
	default:
		return false
	}
}
