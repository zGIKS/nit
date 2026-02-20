package nit

import (
	"fmt"
	"strings"

	"nit/internal/nit/components"
)

func (m model) View() string {
	head := m.header()
	if m.ui == uiMenu {
		return head + m.renderMenu()
	}
	if m.ui == uiPrompt {
		return head + m.renderPrompt() + "\n\n" + m.renderPanels()
	}
	return head + m.renderPanels()
}

func (m model) header() string {
	panel := "Graph+Changes"
	if m.panel == panelOutput {
		panel = "Output"
	}
	focus := "graph"
	if m.focus == focusChanges {
		focus = "changes"
	}
	header := fmt.Sprintf(
		"nit | View: %s | Focus: %s | %s menu | %s switch focus | %s stage | %s unstage | %s stage all | %s unstage all | %s quit\n",
		panel,
		focus,
		primaryKey(m.keys.OpenMenu),
		primaryKey(m.keys.TogglePanel),
		primaryKey(m.keys.StageSelected),
		primaryKey(m.keys.UnstageSelected),
		primaryKey(m.keys.StageAll),
		primaryKey(m.keys.UnstageAll),
		primaryKey(m.keys.Quit),
	)
	if m.status != "" {
		header += "Status: " + m.status + "\n"
	}
	if m.err != "" {
		header += "Warning: " + m.err + "\n"
	}
	header += strings.Repeat("-", max(10, m.width)) + "\n"
	return header
}

func (m model) renderPanels() string {
	if m.panel == panelOutput {
		return components.BoxView("Git Output", m.width, m.bodyHeight(), normalizeLines(m.outputLines), m.cursor, m.offset, true)
	}

	graphActive := m.focus == focusGraph
	changesActive := m.focus == focusChanges
	changes := components.ChangesView(m.width, m.changesPaneHeight(), normalizeLines(m.changeLines), m.changesCursor, m.changesOffset, changesActive)
	graph := components.GraphView(m.width, m.graphPaneHeight(), normalizeLines(m.graphLines), m.cursor, m.offset, graphActive)
	return changes + "\n" + graph
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
