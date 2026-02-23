package handlers

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	"nit/internal/nit/core/model/cmds"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
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

// We need a way to access Model without circular dependency if possible,
// but since handlers depends on Model fields, we can pass State/Git/etc.

func HandleKeyMsg(
	state *app.AppState,
	git g.Service,
	clipCfg config.ClipboardConfig,
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
) tea.Cmd {
	if state.BranchCreateOpen {
		switch msg.Type {
		case tea.KeyEsc:
			state.CloseBranchCreate()
		case tea.KeyEnter:
			name := strings.TrimSpace(state.BranchCreateName)
			if name == "" {
				state.SetError("branch name is empty")
				state.Clamp()
				return nil
			}
			source := strings.TrimSpace(state.BranchCreateSource)
			state.CloseBranchCreate()
			state.BranchCreateName = ""
			state.BranchCreateCursor = 0
			state.BranchCreateSelectAll = false
			state.Clamp()
			return cmds.CreateBranchCmd(git, name, source, false)
		case tea.KeyCtrlB:
			name := strings.TrimSpace(state.BranchCreateName)
			if name == "" {
				state.SetError("branch name is empty")
				state.Clamp()
				return nil
			}
			source := strings.TrimSpace(state.BranchCreateSource)
			state.CloseBranchCreate()
			state.BranchCreateName = ""
			state.BranchCreateCursor = 0
			state.BranchCreateSelectAll = false
			state.Clamp()
			return cmds.CreateBranchCmd(git, name, source, true)
		case tea.KeyUp:
			state.BranchCreateMoveSource(-1)
		case tea.KeyDown:
			state.BranchCreateMoveSource(1)
		case tea.KeyTab:
			state.BranchCreateMoveSource(1)
		case tea.KeyShiftTab:
			state.BranchCreateMoveSource(-1)
		default:
			if handleSharedTextInputKey(state, clipCfg, pasteHintAlreadySeen, msg, textInputKeyOps{
				Selected:        state.SelectedBranchCreateText,
				Append:          state.BranchCreateAppendText,
				Backspace:       state.BranchCreateBackspace,
				Delete:          state.BranchCreateDelete,
				MoveLeft:        state.BranchCreateCursorLeft,
				MoveRight:       state.BranchCreateCursorRight,
				MoveHome:        state.BranchCreateCursorHome,
				MoveEnd:         state.BranchCreateCursorEnd,
				SelectAll:       state.BranchCreateSelectAllText,
				DeleteSelection: state.DeleteBranchCreateSelection,
			}) {
				state.Clamp()
				return nil
			}
		}
		state.Clamp()
		return nil
	}

	if state.Focus == app.FocusCommand {
		switch msg.Type {
		case tea.KeyEnter:
			result := state.Apply(app.ActionToggleOne)
			state.Clamp()
			return cmds.HandleResult(git, result)
		case tea.KeyEsc:
			state.ExitCommandFocus()
			state.Clamp()
			return nil
		}
		if handleSharedTextInputKey(state, clipCfg, pasteHintAlreadySeen, msg, textInputKeyOps{
			Selected:        state.SelectedCommandText,
			Append:          state.AppendCommandText,
			Backspace:       state.BackspaceCommandText,
			Delete:          state.DeleteCommandText,
			MoveLeft:        state.MoveCommandCursorLeft,
			MoveRight:       state.MoveCommandCursorRight,
			MoveHome:        state.MoveCommandCursorToStart,
			MoveEnd:         state.MoveCommandCursorToEnd,
			SelectAll:       state.SelectAllCommandText,
			DeleteSelection: state.DeleteCommandSelection,
		}) {
			state.Clamp()
			return nil
		}

		action := state.Keys.Match(msg.String())
		if action == app.ActionTogglePanel || action == app.ActionPush {
			result := state.Apply(action)
			state.Clamp()
			return cmds.HandleResult(git, result)
		}
		// Ignore quit in commit input to avoid accidental exit while typing.
		state.Clamp()
		return nil
	}

	action := state.Keys.Match(msg.String())
	result := state.Apply(action)
	state.Clamp()
	return cmds.HandleResult(git, result)
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
