package nit

import (
	"fmt"
	"strings"

	"nit/internal/nit/components"
)

func (m model) View() string {
	if m.ui == uiMenu {
		return m.renderMenu()
	}
	if m.ui == uiPrompt {
		return m.renderPrompt() + "\n\n" + m.renderPanels()
	}
	return m.renderPanels()
}

func (m model) renderPanels() string {
	if m.panel == panelOutput {
		total := len(normalizeLines(m.outputLines))
		sel := min(total, m.cursor+1)
		if sel < 1 {
			sel = 1
		}
		return components.BoxView("Git Output", m.width, m.bodyHeight(), normalizeLines(m.outputLines), m.cursor, m.offset, true, fmt.Sprintf("%d of %d", sel, total))
	}

	graphActive := m.focus == focusGraph
	changesActive := m.focus == focusChanges
	commitActive := m.focus == focusCommit
	changeSel, changeTotal := m.changePosition()
	graphSel, graphTotal := m.graphPosition()
	changes := components.ChangesView(m.width, m.changesPaneHeight(), normalizeLines(m.changeLines), m.changesCursor, m.changesOffset, changesActive, changeSel, changeTotal)
	graph := components.GraphView(m.width, m.graphPaneHeight(), normalizeLines(m.graphLines), m.cursor, m.offset, graphActive, graphSel, graphTotal)
	commit := components.CommitView(m.width, m.commitPaneHeight(), m.commitMessage, commitActive, primaryKey(m.keys.CommitSubmit))
	return commit + "\n" + changes + "\n" + graph
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

func (m model) renderMenu() string {
	var b strings.Builder
	b.WriteString("Actions (select/run/cancel via keymap):\n")
	for i, item := range m.menuItems {
		prefix := "  "
		if i == m.menuCursor {
			prefix = "> "
		}
		b.WriteString(prefix + item.label + "\n")
	}
	return b.String()
}

func (m model) renderPrompt() string {
	value := m.promptValue
	if strings.TrimSpace(value) == "" {
		value = "(" + m.promptPlaceholder + ")"
	}
	return fmt.Sprintf("%s: %s\nsubmit/cancel/backspace via keymap", m.promptTitle, value)
}
