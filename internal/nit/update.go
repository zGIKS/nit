package nit

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clamp()
		return m, nil
	case cmdResultMsg:
		m.outputLines = normalizeLines(msg.output)
		if msg.err != nil {
			m.status = fmt.Sprintf("%s failed: %v", msg.title, msg.err)
		} else {
			m.status = fmt.Sprintf("%s completed", msg.title)
		}
		if msg.switchToOutput {
			m.panel = panelOutput
		}
		m.ui = uiBrowse
		m.cursor = 0
		m.offset = 0
		m.refreshGraphAndChanges()
		m.setActiveLines()
		m.clamp()
		return m, nil
	case tea.KeyMsg:
		switch m.ui {
		case uiPrompt:
			return m.updatePrompt(msg)
		case uiMenu:
			return m.updateMenu(msg)
		default:
			return m.updateBrowse(msg)
		}
	}
	return m, nil
}

func (m model) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch {
	case hasKey(m.keys.Quit, key):
		return m, tea.Quit
	case hasKey(m.keys.OpenMenu, key):
		m.ui = uiMenu
		return m, nil
	case hasKey(m.keys.TogglePanel, key):
		if m.panel == panelOutput {
			m.panel = panelGraph
		}
		if m.focus == focusGraph {
			m.focus = focusChanges
		} else {
			m.focus = focusGraph
		}
		m.setActiveLines()
		m.clamp()
	case hasKey(m.keys.ShowOutput, key):
		if m.panel == panelOutput {
			m.panel = panelGraph
		} else {
			m.panel = panelOutput
		}
		m.setActiveLines()
		m.cursor = 0
		m.offset = 0
		m.clamp()
	case hasKey(m.keys.Reload, key):
		m.refreshGraphAndChanges()
		m.status = "Reloaded"
		m.clamp()
	case hasKey(m.keys.StageSelected, key):
		if m.panel == panelGraph && m.focus == focusChanges {
			return m.stageSelected()
		}
	case hasKey(m.keys.UnstageSelected, key):
		if m.panel == panelGraph && m.focus == focusChanges {
			return m.unstageSelected()
		}
	case hasKey(m.keys.StageAll, key):
		m.status = "Running git add -A..."
		return m, runCommandWithOutputMode("Stage All", false, "git", "add", "-A")
	case hasKey(m.keys.UnstageAll, key):
		m.status = "Running git restore --staged . ..."
		return m, runShellCommandWithOutputMode("Unstage All", false, "git restore --staged . || git reset HEAD -- .")
	case hasKey(m.keys.Down, key):
		m.moveCursor(1)
	case hasKey(m.keys.Up, key):
		m.moveCursor(-1)
	case hasKey(m.keys.PageDown, key):
		m.moveCursor(m.pageSize())
	case hasKey(m.keys.PageUp, key):
		m.moveCursor(-m.pageSize())
	case hasKey(m.keys.Home, key):
		m.moveHome()
	case hasKey(m.keys.End, key):
		m.moveEnd()
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch {
	case hasKey(m.keys.MenuClose, key):
		m.ui = uiBrowse
		return m, nil
	case hasKey(m.keys.MenuDown, key):
		m.menuCursor++
		if m.menuCursor >= len(m.menuItems) {
			m.menuCursor = len(m.menuItems) - 1
		}
		return m, nil
	case hasKey(m.keys.MenuUp, key):
		m.menuCursor--
		if m.menuCursor < 0 {
			m.menuCursor = 0
		}
		return m, nil
	case hasKey(m.keys.MenuSelect, key):
		if m.menuCursor < 0 || m.menuCursor >= len(m.menuItems) {
			m.ui = uiBrowse
			return m, nil
		}
		selected := m.menuItems[m.menuCursor]
		m.ui = uiBrowse
		return m.applyAction(selected)
	}
	return m, nil
}

func (m model) updatePrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch {
	case hasKey(m.keys.PromptCancel, key):
		m.ui = uiBrowse
		m.promptValue = ""
		m.promptKind = ""
		return m, nil
	case hasKey(m.keys.PromptSubmit, key):
		value := strings.TrimSpace(m.promptValue)
		kind := m.promptKind
		m.ui = uiBrowse
		m.promptValue = ""
		m.promptKind = ""
		if value == "" {
			m.status = "Prompt canceled (empty input)"
			return m, nil
		}
		return m.runPromptAction(kind, value)
	case hasKey(m.keys.PromptBackspace, key):
		if len(m.promptValue) > 0 {
			m.promptValue = m.promptValue[:len(m.promptValue)-1]
		}
		return m, nil
	case msg.Type == tea.KeySpace:
		m.promptValue += " "
		return m, nil
	default:
		r := msg.Runes
		if len(r) > 0 {
			m.promptValue += string(r)
		}
		return m, nil
	}
}

func (m *model) moveCursor(delta int) {
	if m.panel == panelOutput || m.focus == focusGraph {
		m.cursor += delta
		m.clamp()
		return
	}
	m.changesCursor += delta
	m.clamp()
}

func (m *model) moveHome() {
	if m.panel == panelOutput || m.focus == focusGraph {
		m.cursor = 0
		m.offset = 0
		m.clamp()
		return
	}
	m.changesCursor = 0
	m.changesOffset = 0
	m.clamp()
}

func (m *model) moveEnd() {
	if m.panel == panelOutput || m.focus == focusGraph {
		m.cursor = len(m.lines) - 1
		m.clamp()
		return
	}
	m.changesCursor = len(m.changeLines) - 1
	m.clamp()
}
