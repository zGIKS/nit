package nit

import "nit/internal/nit/components"

func (m model) View() string {
	changesActive := m.focus == focusChanges
	graphActive := m.focus == focusGraph
	changeSel, changeTotal := m.changePosition()
	graphSel, graphTotal := m.graphPosition()

	changes := components.ChangesView(
		m.width,
		m.changesPaneHeight(),
		normalizeLines(m.changeLines),
		m.changesCursor,
		m.changesOffset,
		changesActive,
		changeSel,
		changeTotal,
	)
	graph := components.GraphView(
		m.width,
		m.graphPaneHeight(),
		normalizeLines(m.graphLines),
		m.cursor,
		m.offset,
		graphActive,
		graphSel,
		graphTotal,
	)
	return changes + "\n" + graph
}

func (m model) changePosition() (int, int) {
	total := len(m.unstagedChanges) + len(m.stagedChanges)
	if total < 1 {
		return 1, 1
	}
	cur := 1
	seen := 0
	for i, row := range m.changeRows {
		if row.selectable {
			seen++
		}
		if i == m.changesCursor {
			if row.selectable {
				cur = seen
			}
			break
		}
	}
	if cur > total {
		cur = total
	}
	return cur, total
}

func (m model) graphPosition() (int, int) {
	total := len(normalizeLines(m.graphLines))
	if total < 1 {
		return 1, 1
	}
	cur := m.cursor + 1
	if cur < 1 {
		cur = 1
	}
	if cur > total {
		cur = total
	}
	return cur, total
}
