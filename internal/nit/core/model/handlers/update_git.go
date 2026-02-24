package handlers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGIKS/nit/internal/nit/app"
	"github.com/zGIKS/nit/internal/nit/core/model/cmds"
	"github.com/zGIKS/nit/internal/nit/core/model/common"
	g "github.com/zGIKS/nit/internal/nit/git"
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
	cmdsToRun := make([]tea.Cmd, 0, 2)
	if msg.RefreshChanges {
		cmdsToRun = append(cmdsToRun, cmds.LoadChangesCmd(git))
	}
	if msg.RefreshGraph {
		cmdsToRun = append(cmdsToRun, cmds.LoadGraphCmd(git))
		cmdsToRun = append(cmdsToRun, cmds.LoadBranchesCmd(git))
	}
	if msg.RefreshRepoSummary {
		cmdsToRun = append(cmdsToRun, cmds.LoadRepoSummaryCmd(git))
	}
	state.Clamp()
	if len(cmdsToRun) == 0 {
		return nil
	}
	return tea.Batch(cmdsToRun...)
}
