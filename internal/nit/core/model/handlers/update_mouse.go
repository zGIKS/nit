package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/core/model/cmds"
	g "nit/internal/nit/git"
)

func HandleMouseMsg(state *app.AppState, git g.Service, msg tea.MouseMsg) tea.Cmd {
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		if action, ok := state.TopBarActionAt(msg.X, msg.Y); ok {
			result := state.Apply(action)
			state.Clamp()
			return cmds.HandleResult(git, result)
		}
		state.HandleMouseClick(msg.X, msg.Y)
		state.Clamp()
		return nil
	}
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelUp {
		state.HandleMouseWheel(msg.X, msg.Y, -1)
		state.Clamp()
		return nil
	}
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelDown {
		state.HandleMouseWheel(msg.X, msg.Y, 1)
		state.Clamp()
		return nil
	}
	return nil
}
