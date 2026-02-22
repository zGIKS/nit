package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
)

func HandleMouseMsg(state *app.AppState, msg tea.MouseMsg) tea.Cmd {
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
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
