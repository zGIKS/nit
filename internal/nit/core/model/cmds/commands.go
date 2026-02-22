package cmds

import (
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/core/model/common"
	g "nit/internal/nit/git"
)

const defaultPollInterval = 1500 * time.Millisecond

func pollInterval() time.Duration {
	raw := strings.TrimSpace(os.Getenv("NIT_POLL_MS"))
	if raw == "" {
		return defaultPollInterval
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 250 {
		return defaultPollInterval
	}
	return time.Duration(ms) * time.Millisecond
}

func SchedulePoll() tea.Cmd {
	return tea.Tick(pollInterval(), func(time.Time) tea.Msg { return common.PollMsg{} })
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

func LoadRepoSummaryCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		repo, branch, err := svc.LoadRepoSummary()
		return common.RepoSummaryLoadedMsg{Repo: repo, Branch: branch, Err: err}
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
