package nit

func initialModel(keys keyConfig) model {
	graph, _ := loadGraphLines()
	changes, _ := loadChanges()

	m := model{
		focus:         focusChanges,
		graphLines:    normalizeLines(graph),
		changeEntries: changes,
		height:        24,
		keys:          keys,
	}
	m.rebuildChangeRows()
	m.snapChangesCursorToSelectable(1)
	m.clamp()
	return m
}

func (m *model) refreshData() {
	graph, _ := loadGraphLines()
	changes, _ := loadChanges()
	m.graphLines = normalizeLines(graph)
	m.changeEntries = changes
	m.rebuildChangeRows()
}

func (m *model) rebuildChangeRows() {
	m.stagedChanges = make([]changeEntry, 0, len(m.changeEntries))
	m.unstagedChanges = make([]changeEntry, 0, len(m.changeEntries))
	for _, e := range m.changeEntries {
		if e.staged {
			m.stagedChanges = append(m.stagedChanges, e)
		}
		if e.changed || !e.staged {
			m.unstagedChanges = append(m.unstagedChanges, e)
		}
	}

	rows := make([]changeRow, 0, len(m.changeEntries)+4)
	if len(m.stagedChanges) > 0 {
		rows = append(rows, changeRow{text: "Staged Changes"})
		for i, e := range m.stagedChanges {
			rows = append(rows, changeRow{
				text:       "  " + changeCodeForStaged(e) + "  " + e.path,
				selectable: true,
				section:    "staged",
				index:      i,
			})
		}
	}
	if len(m.unstagedChanges) > 0 {
		rows = append(rows, changeRow{text: "Changes"})
		for i, e := range m.unstagedChanges {
			rows = append(rows, changeRow{
				text:       "  " + changeCodeForUnstaged(e) + "  " + e.path,
				selectable: true,
				section:    "unstaged",
				index:      i,
			})
		}
	}

	m.changeRows = rows
	m.changeLines = make([]string, 0, len(rows))
	for _, r := range rows {
		m.changeLines = append(m.changeLines, r.text)
	}
	if len(m.changeLines) == 0 {
		m.changeLines = []string{"Working tree clean."}
	}
}

func changeCodeForStaged(e changeEntry) string {
	if e.x == '?' {
		return "U"
	}
	if e.x != ' ' {
		return string(e.x)
	}
	if e.staged {
		return "M"
	}
	return "-"
}

func changeCodeForUnstaged(e changeEntry) string {
	if e.x == '?' {
		return "U"
	}
	if e.y != ' ' {
		return string(e.y)
	}
	if !e.staged {
		return "M"
	}
	return "-"
}

func (m model) selectedChange() (changeEntry, bool) {
	if len(m.changeRows) == 0 || m.changesCursor < 0 || m.changesCursor >= len(m.changeRows) {
		return changeEntry{}, false
	}
	row := m.changeRows[m.changesCursor]
	if !row.selectable {
		return changeEntry{}, false
	}
	if row.section == "staged" {
		if row.index < 0 || row.index >= len(m.stagedChanges) {
			return changeEntry{}, false
		}
		return m.stagedChanges[row.index], true
	}
	if row.index < 0 || row.index >= len(m.unstagedChanges) {
		return changeEntry{}, false
	}
	return m.unstagedChanges[row.index], true
}

func (m *model) moveToFirstSelectableSection(section string) bool {
	for i, row := range m.changeRows {
		if row.selectable && row.section == section {
			m.changesCursor = i
			m.changesOffset = 0
			return true
		}
	}
	return false
}
