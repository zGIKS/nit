package handlers

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGIKS/nit/internal/nit/app"
	"github.com/zGIKS/nit/internal/nit/config"
	"github.com/zGIKS/nit/internal/nit/core/model/cmds"
	g "github.com/zGIKS/nit/internal/nit/git"
)

func handleBranchCreateKey(
	state *app.AppState,
	git g.Service,
	clipCfg config.ClipboardConfig,
	textKeys config.CommitEditorKeyConfig,
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
) tea.Cmd {
	switch msg.Type {
	case tea.KeyUp:
		state.BranchCreateMoveSource(-1)
	case tea.KeyDown:
		state.BranchCreateMoveSource(1)
	case tea.KeyTab:
		state.BranchCreateMoveSource(1)
	case tea.KeyShiftTab:
		state.BranchCreateMoveSource(-1)
	default:
		switch {
		case matchesConfiguredKey(msg, textKeys.Cancel):
			state.CloseBranchCreate()
		case matchesConfiguredKey(msg, textKeys.Submit):
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
		default:
			if handleSharedTextInputKey(state, clipCfg, textKeys, pasteHintAlreadySeen, msg, textInputKeyOps{
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
	}
	state.Clamp()
	return nil
}
