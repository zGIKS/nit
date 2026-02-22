package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/core/model/cmds"
	"nit/internal/nit/core/model/common"
	"nit/internal/nit/core/model/handlers"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.State.Viewport.Width == msg.Width && m.State.Viewport.Height == msg.Height {
			return m, nil
		}
		m.State.SetViewport(msg.Width, msg.Height)
		return m, nil

	case common.PollMsg:
		return m, tea.Batch(cmds.ScheduleChangesPoll(), cmds.LoadChangesCmd(m.Git))

	case common.GraphPollMsg:
		return m, tea.Batch(cmds.ScheduleGraphPoll(), cmds.LoadGraphCmd(m.Git), cmds.LoadRepoSummaryCmd(m.Git))

	case common.WatchReadyMsg:
		if msg.Err != nil {
			// Fallback polling remains active; surface the watcher error once.
			m.State.SetError(msg.Err.Error())
			return m, nil
		}
		m.Watcher = msg.Watcher
		return m, cmds.WaitWatchCmd(m.Watcher)

	case common.WatchTickMsg:
		return m, tea.Batch(
			cmds.WaitWatchCmd(m.Watcher),
			cmds.LoadChangesCmd(m.Git),
			cmds.LoadRepoSummaryCmd(m.Git),
		)

	case common.ChangesLoadedMsg:
		return m, handlers.HandleChangesLoaded(&m.State, msg)

	case common.GraphLoadedMsg:
		return m, handlers.HandleGraphLoaded(&m.State, msg)

	case common.RepoSummaryLoadedMsg:
		return m, handlers.HandleRepoSummaryLoaded(&m.State, msg)

	case common.OpDoneMsg:
		return m, handlers.HandleOpDone(&m.State, m.Git, msg)

	case tea.KeyMsg:
		return m, handlers.HandleKeyMsg(&m.State, m.Git, m.ClipCfg, &m.PasteHintAlreadySeen, msg)

	case tea.MouseMsg:
		return m, handlers.HandleMouseMsg(&m.State, msg)
	}

	return m, nil
}
