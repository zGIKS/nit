package nit

import tea "github.com/charmbracelet/bubbletea"

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clamp()
		return m, nil
	case tea.KeyMsg:
		return m.updateBrowse(msg)
	}
	return m, nil
}

func (m model) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch {
	case hasKey(m.keys.Quit, key):
		return m, tea.Quit
	case hasKey(m.keys.TogglePanel, key):
		if m.focus == focusChanges {
			m.focus = focusGraph
		} else {
			m.focus = focusChanges
			m.snapChangesCursorToSelectable(1)
		}
		m.clamp()
	case hasKey(m.keys.Down, key):
		m.moveCursor(1)
	case hasKey(m.keys.Up, key):
		m.moveCursor(-1)
	case m.focus == focusChanges && hasKey(m.keys.ToggleOne, key):
		m.toggleSelectedChange()
	case m.focus == focusChanges && hasKey(m.keys.StageAll, key):
		if err := stageAll(); err == nil {
			m.refreshAfterChangeAction("stage_all")
		}
	case m.focus == focusChanges && hasKey(m.keys.UnstageAll, key):
		if err := unstageAll(); err == nil {
			m.refreshAfterChangeAction("unstage_all")
		}
	}
	return m, nil
}

func (m *model) moveCursor(delta int) {
	if m.focus == focusGraph {
		m.cursor += delta
		m.clamp()
		return
	}
	m.changesCursor += delta
	if delta >= 0 {
		m.snapChangesCursorToSelectable(1)
	} else {
		m.snapChangesCursorToSelectable(-1)
	}
	m.clamp()
}

func (m *model) snapChangesCursorToSelectable(dir int) {
	if len(m.changeRows) == 0 {
		m.changesCursor = 0
		return
	}
	if m.changesCursor < 0 {
		m.changesCursor = 0
	}
	if m.changesCursor >= len(m.changeRows) {
		m.changesCursor = len(m.changeRows) - 1
	}
	if m.changeRows[m.changesCursor].selectable {
		return
	}

	i := m.changesCursor
	for i >= 0 && i < len(m.changeRows) {
		if m.changeRows[i].selectable {
			m.changesCursor = i
			return
		}
		i += dir
	}

	for i = 0; i < len(m.changeRows); i++ {
		if m.changeRows[i].selectable {
			m.changesCursor = i
			return
		}
	}
	m.changesCursor = 0
}

func (m *model) toggleSelectedChange() {
	entry, ok := m.selectedChange()
	if !ok {
		return
	}
	if entry.staged {
		if err := unstagePath(entry.path); err == nil {
			m.refreshAfterChangeAction("unstage_one")
		}
		return
	}
	if err := stagePath(entry.path); err == nil {
		m.refreshAfterChangeAction("stage_one")
	}
}

func (m *model) refreshAfterChangeAction(action string) {
	m.refreshData()
	m.focus = focusChanges
	switch action {
	case "stage_one":
		// Keep workflow fast: continue in unstaged while there are files left.
		if !m.moveToFirstSelectableSection("unstaged") {
			m.moveToFirstSelectableSection("staged")
		}
	case "stage_all":
		m.moveToFirstSelectableSection("staged")
	case "unstage_one", "unstage_all":
		if !m.moveToFirstSelectableSection("staged") {
			m.moveToFirstSelectableSection("unstaged")
		}
	}
	m.snapChangesCursorToSelectable(1)
	m.clamp()
}
