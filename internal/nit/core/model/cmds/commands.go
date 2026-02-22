package cmds

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
)

const pollInterval = 700 * time.Millisecond

func SchedulePoll() tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg { return common.PollMsg{} })
}

func LoadChangesCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		entries, err := svc.LoadChanges()
		return common.ChangesLoadedMsg{Entries: entries, Err: err}
	}
}

func LoadGraphCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		lines, err := svc.LoadGraph()
		return common.GraphLoadedMsg{Lines: lines, Err: err}
	}
}

func ExecOpCmd(svc g.Service, op app.Operation, refreshChanges, refreshGraph bool) tea.Cmd {
	return func() tea.Msg {
		cmd, err := ExecOperation(svc, op)
		if err != nil {
			return common.OpDoneMsg{Err: err, Command: cmd}
		}
		return common.OpDoneMsg{RefreshChanges: refreshChanges, RefreshGraph: refreshGraph, Command: cmd}
	}
}
