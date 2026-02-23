package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/core/model/cmds"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
)

func HandleChangesLoaded(state *app.AppState, msg common.ChangesLoadedMsg) tea.Cmd {
	return handleLoadResult(state, msg.Err, func() {
		if !common.SameChanges(state.Changes.Entries, msg.Entries) {
			state.SetChanges(msg.Entries)
		}
	})
}

func HandleGraphLoaded(state *app.AppState, msg common.GraphLoadedMsg) tea.Cmd {
	return handleLoadResult(state, msg.Err, func() { state.SetGraph(msg.Lines) })
}

func HandleBranchesLoaded(state *app.AppState, msg common.BranchesLoadedMsg) tea.Cmd {
	return handleLoadResult(state, msg.Err, func() { state.SetBranches(msg.Lines) })
}

func HandleRepoSummaryLoaded(state *app.AppState, msg common.RepoSummaryLoadedMsg) tea.Cmd {
	if msg.Err != nil {
		if state.RepoName == "" {
			state.SetRepoSummary("not-a-repo", "-")
		}
		state.SetError(msg.Err.Error())
	} else {
		state.SetRepoSummary(msg.Repo, msg.Branch)
		if state.LastErr == "" {
			state.SetError("")
		}
	}
	state.Clamp()
	return nil
}

func handleLoadResult(state *app.AppState, err error, onSuccess func()) tea.Cmd {
	if err != nil {
		state.SetError(err.Error())
	} else {
		state.SetError("")
		if onSuccess != nil {
			onSuccess()
		}
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
		cmds_to_run = append(cmds_to_run, cmds.LoadBranchesCmd(git))
	}
	state.Clamp()
	if len(cmds_to_run) == 0 {
		return nil
	}
	return tea.Batch(cmds_to_run...)
}
