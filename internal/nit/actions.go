package nit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) applyAction(a action) (model, tea.Cmd) {
	switch a.kind {
	case "view_graph":
		m.panel = panelGraph
		m.refreshGraphAndChanges()
		m.setActiveLines()
		m.status = "Graph + Changes view"
		m.clamp()
		return m, nil
	case "pull":
		m.status = "Running git pull..."
		return m, runCommand("Pull", "git", "pull")
	case "push":
		m.status = "Running git push..."
		return m, runCommand("Push", "git", "push")
	case "fetch":
		m.status = "Running git fetch --all --prune..."
		return m, runCommand("Fetch", "git", "fetch", "--all", "--prune")
	case "changes":
		m.status = "Running git status --short --branch..."
		return m, runCommand("Changes", "git", "status", "--short", "--branch")
	case "branch":
		m.status = "Running git branch -vv..."
		return m, runCommand("Branch", "git", "branch", "-vv")
	case "remote":
		m.status = "Running git remote -v..."
		return m, runCommand("Remote", "git", "remote", "-v")
	case "tags":
		m.status = "Running git tag --list..."
		return m, runCommand("Tags", "git", "tag", "--list")
	case "stash":
		m.status = "Running git stash list..."
		return m, runCommand("Stash", "git", "stash", "list")
	case "pull_push":
		m.status = "Running git pull && git push..."
		return m, runShellCommand("Pull, Push", "git pull && git push")
	case "show_output":
		m.panel = panelOutput
		m.setActiveLines()
		m.cursor, m.offset = 0, 0
		m.status = "Output view"
		m.clamp()
		return m, nil
	case "stage_selected":
		return m.stageSelected()
	case "unstage_selected":
		return m.unstageSelected()
	case "stage_all":
		m.status = "Running git add -A..."
		return m, runCommandWithOutputMode("Stage All", false, "git", "add", "-A")
	case "unstage_all":
		m.status = "Running git restore --staged . ..."
		return m, runShellCommandWithOutputMode("Unstage All", false, "git restore --staged . || git reset HEAD -- .")
	case "checkout":
		m.startPrompt("checkout", "Checkout to...", "branch/ref")
		return m, nil
	case "commit":
		m.startPrompt("commit", "Commit", "message")
		return m, nil
	case "clone":
		m.startPrompt("clone", "Clone", "url [directory]")
		return m, nil
	default:
		m.status = "Unknown action"
		return m, nil
	}
}

func (m model) stageSelected() (model, tea.Cmd) {
	entry, ok := m.selectedChange()
	if !ok {
		m.status = "No file selected in Changes"
		return m, nil
	}
	m.status = "Staging: " + entry.path
	return m, runCommandWithOutputMode("Stage", false, "git", "add", "--", entry.path)
}

func (m model) unstageSelected() (model, tea.Cmd) {
	entry, ok := m.selectedChange()
	if !ok {
		m.status = "No file selected in Changes"
		return m, nil
	}
	m.status = "Unstaging: " + entry.path
	cmd := "git restore --staged -- " + shellQuote(entry.path) + " || git reset HEAD -- " + shellQuote(entry.path)
	return m, runShellCommandWithOutputMode("Unstage", false, cmd)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func (m model) selectedChange() (changeEntry, bool) {
	if len(m.changeEntries) == 0 {
		return changeEntry{}, false
	}
	if m.changesCursor < 0 || m.changesCursor >= len(m.changeEntries) {
		return changeEntry{}, false
	}
	return m.changeEntries[m.changesCursor], true
}

func (m model) runPromptAction(kind, value string) (model, tea.Cmd) {
	switch kind {
	case "checkout":
		m.status = "Running git checkout..."
		return m, runCommand("Checkout", "git", "checkout", value)
	case "commit":
		m.status = "Running git commit -m ..."
		return m, runCommand("Commit", "git", "commit", "-m", value)
	case "clone":
		parts := strings.Fields(value)
		if len(parts) == 0 {
			m.status = "Clone canceled (empty input)"
			return m, nil
		}
		args := append([]string{"clone"}, parts...)
		m.status = "Running git clone..."
		return m, runCommand("Clone", "git", args...)
	default:
		m.status = "Unknown prompt action"
		return m, nil
	}
}

func (m *model) startPrompt(kind, title, placeholder string) {
	m.ui = uiPrompt
	m.promptKind = kind
	m.promptTitle = title
	m.promptPlaceholder = placeholder
	m.promptValue = ""
}

func (m *model) setActiveLines() {
	if m.panel == panelOutput {
		m.lines = normalizeLines(m.outputLines)
		return
	}
	if m.focus == focusGraph {
		m.lines = normalizeLines(m.graphLines)
	} else {
		m.lines = normalizeLines(m.changeLines)
	}
}

func (m *model) refreshGraphAndChanges() {
	graph, graphErr := loadGraphLines()
	changes, changeLines, changesErr := loadChanges()
	m.graphLines = normalizeLines(graph)
	m.changeEntries = changes
	m.changeLines = normalizeLines(changeLines)

	errs := make([]string, 0, 2)
	if graphErr != nil {
		errs = append(errs, "graph error: "+graphErr.Error())
	}
	if changesErr != nil {
		errs = append(errs, "changes error: "+changesErr.Error())
	}
	m.err = strings.Join(errs, " | ")

	if len(m.changeEntries) == 0 {
		m.changesCursor = 0
		m.changesOffset = 0
	} else if m.changesCursor >= len(m.changeEntries) {
		m.changesCursor = len(m.changeEntries) - 1
	}

	m.setActiveLines()
}
