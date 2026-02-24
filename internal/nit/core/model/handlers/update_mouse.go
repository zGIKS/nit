package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGIKS/nit/internal/nit/app"
	"github.com/zGIKS/nit/internal/nit/core/model/cmds"
	g "github.com/zGIKS/nit/internal/nit/git"
)

func HandleMouseMsg(state *app.AppState, git g.Service, msg tea.MouseMsg) tea.Cmd {
	if msg.Action == tea.MouseActionMotion {
		state.HandleMouseMove(msg.X, msg.Y)
		state.Clamp()
		return nil
	}
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		if state.BranchCreateOpen {
			if state.BranchCreateClick(msg.X, msg.Y) {
				state.Clamp()
				return nil
			}
		}
		if action, ok, consumed := state.MenuSubmenuClickActionAt(msg.X, msg.Y); consumed {
			if ok {
				result := state.Apply(action)
				state.Clamp()
				return cmds.HandleResult(git, result)
			}
			state.Clamp()
			return nil
		}
		if _, ok := state.MenuItemIndexAt(msg.X, msg.Y); ok {
			if action, ok := state.MenuClickActionAt(msg.X, msg.Y); ok {
				result := state.Apply(action)
				state.Clamp()
				return cmds.HandleResult(git, result)
			}
			state.Clamp()
			return nil
		}
		if action, ok := state.MenuClickActionAt(msg.X, msg.Y); ok {
			result := state.Apply(action)
			state.Clamp()
			return cmds.HandleResult(git, result)
		}
		if state.ToggleMenuClick(msg.X, msg.Y) {
			state.Clamp()
			return nil
		}
		if state.ToggleBranchCreateClick(msg.X, msg.Y) {
			state.Clamp()
			return nil
		}
		if action, ok := state.TopBarActionAt(msg.X, msg.Y); ok {
			result := state.Apply(action)
			state.Clamp()
			return cmds.HandleResult(git, result)
		}
		state.CloseTopMenusOnOutsideClick(msg.X, msg.Y)
		state.HandleMouseClick(msg.X, msg.Y)
		state.Clamp()
		return nil
	}
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelUp {
		if state.BranchCreateWheelAt(msg.X, msg.Y, -1) {
			state.Clamp()
			return nil
		}
		state.HandleMouseWheel(msg.X, msg.Y, -1)
		state.Clamp()
		return nil
	}
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelDown {
		if state.BranchCreateWheelAt(msg.X, msg.Y, 1) {
			state.Clamp()
			return nil
		}
		state.HandleMouseWheel(msg.X, msg.Y, 1)
		state.Clamp()
		return nil
	}
	return nil
}
