package nit

import (
	"fmt"
	"strings"
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
		return m.renderBox("Git Output", m.bodyHeight(), normalizeLines(m.outputLines), m.cursor, m.offset, true)
	}

	graphActive := m.focus == focusGraph
	changesActive := m.focus == focusChanges
	graph := m.renderBox("Graph", m.graphPaneHeight(), normalizeLines(m.graphLines), m.cursor, m.offset, graphActive)
	changes := m.renderBox("Changes", m.changesPaneHeight(), normalizeLines(m.changeLines), m.changesCursor, m.changesOffset, changesActive)
	return graph + "\n" + changes
}

func (m model) renderBox(title string, boxHeight int, lines []string, cursor, offset int, active bool) string {
	w := max(30, m.width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	label := title
	if active {
		label += " *"
	}
	top := "+" + fitText(label, innerW, '-') + "+"
	bottom := "+" + strings.Repeat("-", innerW) + "+"

	var b strings.Builder
	b.WriteString(top + "\n")

	maxOffset := max(0, len(lines)-contentHeight)
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	end := min(len(lines), offset+contentHeight)
	for i := 0; i < contentHeight; i++ {
		idx := offset + i
		text := ""
		if idx < end {
			prefix := "  "
			if idx == cursor {
				prefix = "> "
			}
			text = prefix + lines[idx]
		}
		text = fitText(text, innerW, ' ')
		b.WriteString("|" + text + "|\n")
	}

	b.WriteString(bottom)
	return b.String()
}

func fitText(text string, width int, fill rune) string {
	if width <= 0 {
		return ""
	}
	if len(text) > width {
		if width <= 3 {
			return text[:width]
		}
		return text[:width-3] + "..."
	}
	if len(text) == width {
		return text
	}
	return text + strings.Repeat(string(fill), width-len(text))
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
