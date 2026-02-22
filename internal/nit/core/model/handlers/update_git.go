package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/core/model/cmds"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
)

func HandleChangesLoaded(state *app.AppState, msg common.ChangesLoadedMsg) tea.Cmd {
	if msg.Err != nil {
		state.SetError(msg.Err.Error())
	} else {
		state.SetError("")
		if !common.SameChanges(state.Changes.Entries, msg.Entries) {
			state.SetChanges(msg.Entries)
		}
	}
	state.Clamp()
	return nil
}

func HandleGraphLoaded(state *app.AppState, msg common.GraphLoadedMsg) tea.Cmd {
	if msg.Err != nil {
		state.SetError(msg.Err.Error())
	} else {
		state.SetError("")
		state.SetGraph(msg.Lines)
	}
	state.Clamp()
	return nil
}

func HandleOpDone(state *app.AppState, git g.Service, msg common.OpDoneMsg) tea.Cmd {
	if msg.Command != "" {
		state.AddCommandLog(msg.Command)
	}
	if msg.Err != nil {
		state.SetError(msg.Err.Error())
		state.Clamp()
		return nil
	}
	state.SetError("")
	cmds_to_run := make([]tea.Cmd, 0, 2)
	if msg.RefreshChanges {
		cmds_to_run = append(cmds_to_run, cmds.LoadChangesCmd(git))
	}
	if msg.RefreshGraph {
		cmds_to_run = append(cmds_to_run, cmds.LoadGraphCmd(git))
	}
	state.Clamp()
	if len(cmds_to_run) == 0 {
		return nil
	}
	return tea.Batch(cmds_to_run...)
}
