package handlers

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGIKS/nit/internal/nit/app"
	"github.com/zGIKS/nit/internal/nit/config"
	"github.com/zGIKS/nit/internal/nit/core/model/cmds"
	g "github.com/zGIKS/nit/internal/nit/git"
)

func HandleKeyMsg(
	state *app.AppState,
	git g.Service,
	clipCfg config.ClipboardConfig,
	textKeys config.CommitEditorKeyConfig,
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
) tea.Cmd {
	if state.BranchCreateOpen {
		return handleBranchCreateKey(state, git, clipCfg, textKeys, pasteHintAlreadySeen, msg)
	}

	if state.MenuOpen {
		action := state.Keys.Match(msg.String())
		switch action {
		case app.ActionMoveUp:
			if state.MenuSubActive && state.MenuSubmenuKind != "" {
				state.MoveMenuSubmenuSelection(-1)
			} else {
				state.MoveMenuSelection(-1)
			}
			state.Clamp()
			return nil
		case app.ActionMoveDown:
			if state.MenuSubActive && state.MenuSubmenuKind != "" {
				state.MoveMenuSubmenuSelection(1)
			} else {
				state.MoveMenuSelection(1)
			}
			state.Clamp()
			return nil
		case app.ActionMenuRight:
			if !state.MenuSubActive && state.MenuHoverHasSubmenu() {
				state.OpenHoveredSubmenu()
			}
			state.Clamp()
			return nil
		case app.ActionMenuLeft:
			if state.MenuSubActive {
				state.MenuSubActive = false
				state.MenuSubHoverIndex = -1
			}
			state.Clamp()
			return nil
		case app.ActionToggleOne:
			if state.MenuSubActive && state.MenuSubmenuKind != "" {
				if action, ok, consumed := state.MenuSubmenuActivateIndex(state.MenuSubHoverIndex); consumed {
					if ok {
						result := state.Apply(action)
						state.Clamp()
						return cmds.HandleResult(git, result)
					}
					state.Clamp()
					return nil
				}
			} else {
				if action, ok := state.MenuActivateIndex(state.MenuHoverIndex); ok {
					result := state.Apply(action)
					state.Clamp()
					return cmds.HandleResult(git, result)
				}
			}
			state.Clamp()
			return nil
		}

		switch msg.Type {
		case tea.KeyEsc:
			if state.MenuSubActive {
				state.MenuSubActive = false
				state.MenuSubHoverIndex = -1
			} else {
				state.CloseMenu()
			}
			state.Clamp()
			return nil
		}
	}

	if state.Focus == app.FocusCommand {
		switch {
		case matchesConfiguredKey(msg, textKeys.Submit):
			result := state.Apply(app.ActionToggleOne)
			state.Clamp()
			return cmds.HandleResult(git, result)
		case matchesConfiguredKey(msg, textKeys.Cancel):
			state.ExitCommandFocus()
			state.Clamp()
			return nil
		}
		if handleSharedTextInputKey(state, clipCfg, textKeys, pasteHintAlreadySeen, msg, textInputKeyOps{
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
	if state.Focus == app.FocusBranches && msg.Type == tea.KeyEnter {
		branch, ok := state.SelectedBranchName()
		if !ok {
			state.Clamp()
			return nil
		}
		if strings.TrimSpace(branch) == strings.TrimSpace(state.BranchName) {
			state.Clamp()
			return nil
		}
		state.Clamp()
		return cmds.SwitchBranchCmd(git, branch)
	}

	action := state.Keys.Match(msg.String())
	result := state.Apply(action)
	state.Clamp()
	return cmds.HandleResult(git, result)
}
