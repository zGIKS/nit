package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	"nit/internal/nit/core/model/cmds"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
)

// We need a way to access Model without circular dependency if possible, 
// but since handlers depends on Model fields, we can pass State/Git/etc.

func HandleKeyMsg(
	state *app.AppState,
	git g.Service,
	clipCfg config.ClipboardConfig,
	pasteHintAlreadySeen *bool,
	msg tea.KeyMsg,
) tea.Cmd {
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
		case tea.KeyCtrlC:
			selected := state.SelectedCommandText()
			if selected == "" {
				state.Clamp()
				return nil
			}
			state.SetCommandClipboard(selected)
			_ = common.CopyWithMode(clipCfg, selected)
			state.SetError("")
			state.Clamp()
			return nil
		case tea.KeyBackspace:
			state.BackspaceCommandText()
			state.Clamp()
			return nil
		case tea.KeyDelete:
			state.DeleteCommandText()
			state.Clamp()
			return nil
		case tea.KeyLeft:
			state.MoveCommandCursorLeft()
			state.Clamp()
			return nil
		case tea.KeyRight:
			state.MoveCommandCursorRight()
			state.Clamp()
			return nil
		case tea.KeyHome:
			state.MoveCommandCursorToStart()
			state.Clamp()
			return nil
		case tea.KeyEnd, tea.KeyCtrlE:
			state.MoveCommandCursorToEnd()
			state.Clamp()
			return nil
		case tea.KeyCtrlA:
			state.SelectAllCommandText()
			state.Clamp()
			return nil
		case tea.KeyCtrlX:
			selected := state.SelectedCommandText()
			if selected == "" {
				state.Clamp()
				return nil
			}
			state.SetCommandClipboard(selected)
			_ = common.CopyWithMode(clipCfg, selected)
			state.DeleteCommandSelection()
			state.SetError("")
			state.Clamp()
			return nil
		case tea.KeyCtrlV:
			pasted, err := common.PasteWithMode(clipCfg)
			if err != nil || pasted == "" {
				pasted = state.CommandClipboard()
			}
			if pasted == "" {
				if !*pasteHintAlreadySeen && clipCfg.Mode == config.ClipboardOnlyCopy {
					state.SetError("paste from OS disabled in only_copy mode")
					*pasteHintAlreadySeen = true
				}
				state.Clamp()
				return nil
			}
			state.AppendCommandText(pasted)
			state.SetError("")
			state.Clamp()
			return nil
		case tea.KeySpace:
			state.AppendCommandText(" ")
			state.Clamp()
			return nil
		case tea.KeyRunes:
			state.AppendCommandText(string(msg.Runes))
			state.Clamp()
			return nil
		}

		action := state.Keys.Match(msg.String())
		if action == app.ActionTogglePanel {
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
