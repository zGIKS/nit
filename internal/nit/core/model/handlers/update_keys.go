package handlers

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	"nit/internal/nit/core/model/cmds"
	g "nit/internal/nit/git"
)

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
