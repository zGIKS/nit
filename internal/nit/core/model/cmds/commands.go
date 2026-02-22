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

const (
	defaultChangesPollInterval = 4 * time.Second
	defaultGraphPollInterval   = 15 * time.Second
)

func pollInterval(envKey string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(envKey))
	if raw == "" {
		if envKey == "NIT_POLL_MS" {
			raw = strings.TrimSpace(os.Getenv("NIT_POLL_MS"))
		}
		if raw == "" {
			return fallback
		}
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 250 {
		return fallback
	}
	return time.Duration(ms) * time.Millisecond
}

func ScheduleChangesPoll() tea.Cmd {
	return tea.Tick(pollInterval("NIT_POLL_CHANGES_MS", defaultChangesPollInterval), func(time.Time) tea.Msg {
		return common.PollMsg{}
	})
}

func ScheduleGraphPoll() tea.Cmd {
	return tea.Tick(pollInterval("NIT_POLL_GRAPH_MS", defaultGraphPollInterval), func(time.Time) tea.Msg {
		return common.GraphPollMsg{}
	})
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

func InitWatchCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		w, err := svc.NewFSWatcher()
		return common.WatchReadyMsg{Watcher: w, Err: err}
	}
}

func WaitWatchCmd(w *g.FSWatcher) tea.Cmd {
	return func() tea.Msg {
		if w == nil {
			return nil
		}
		if _, ok := <-w.Events(); !ok {
			return nil
		}
		return common.WatchTickMsg{}
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
