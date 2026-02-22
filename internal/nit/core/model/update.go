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
		return m, tea.Batch(cmds.SchedulePoll(), cmds.LoadChangesCmd(m.Git))

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
