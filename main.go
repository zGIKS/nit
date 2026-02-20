package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type panelMode int

type uiMode int

const (
	panelTree panelMode = iota
	panelGraph
	panelOutput
)

const (
	uiBrowse uiMode = iota
	uiMenu
	uiPrompt
)

type action struct {
	label string
	kind  string
}

type cmdResultMsg struct {
	title  string
	output []string
	err    error
}

type model struct {
	ui          uiMode
	panel       panelMode
	treeLines   []string
	graphLines  []string
	outputLines []string
	lines       []string

	cursor int
	offset int
	width  int
	height int

	err    string
	status string

	menuItems  []action
	menuCursor int

	promptTitle       string
	promptPlaceholder string
	promptValue       string
	promptKind        string
}

func initialModel() model {
	tree, treeErr := buildTreeLines(".")
	graph, graphErr := loadGraphLines()

	errs := make([]string, 0, 2)
	if treeErr != nil {
		errs = append(errs, fmt.Sprintf("tree error: %v", treeErr))
	}
	if graphErr != nil {
		errs = append(errs, fmt.Sprintf("graph error: %v", graphErr))
	}

	m := model{
		ui:          uiBrowse,
		panel:       panelGraph,
		treeLines:   normalizeLines(tree),
		graphLines:  normalizeLines(graph),
		outputLines: []string{"Run an action from the menu (m) to see command output."},
		height:      24,
		err:         strings.Join(errs, " | "),
		status:      "Ready",
		menuItems: []action{
			{label: "View as Tree", kind: "view_tree"},
			{label: "View as Graph", kind: "view_graph"},
			{label: "Pull", kind: "pull"},
			{label: "Push", kind: "push"},
			{label: "Clone", kind: "clone"},
			{label: "Checkout to...", kind: "checkout"},
			{label: "Fetch", kind: "fetch"},
			{label: "Commit", kind: "commit"},
			{label: "Changes", kind: "changes"},
			{label: "Pull, Push", kind: "pull_push"},
			{label: "Branch", kind: "branch"},
			{label: "Remote", kind: "remote"},
			{label: "Stash", kind: "stash"},
			{label: "Tags", kind: "tags"},
			{label: "Show Git Output", kind: "show_output"},
		},
	}
	m.setActiveLines()
	m.clamp()
	return m
}

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
		m.panel = panelOutput
		m.ui = uiBrowse
		m.setActiveLines()
		m.cursor = 0
		m.offset = 0
		m.refreshTreeAndGraph()
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
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "m":
		m.ui = uiMenu
		return m, nil
	case "tab":
		if m.panel == panelTree {
			m.panel = panelGraph
		} else {
			m.panel = panelTree
		}
		m.setActiveLines()
		m.cursor = 0
		m.offset = 0
		m.clamp()
	case "o":
		m.panel = panelOutput
		m.setActiveLines()
		m.cursor = 0
		m.offset = 0
		m.clamp()
	case "down", "j":
		m.cursor++
		m.clamp()
	case "up", "k":
		m.cursor--
		m.clamp()
	case "pgdown", "f":
		m.cursor += m.pageSize()
		m.clamp()
	case "pgup", "b":
		m.cursor -= m.pageSize()
		m.clamp()
	case "home", "g":
		m.cursor = 0
		m.offset = 0
		m.clamp()
	case "end", "G":
		m.cursor = len(m.lines) - 1
		m.clamp()
	case "r":
		m.refreshTreeAndGraph()
		m.status = "Reloaded"
		m.clamp()
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ui = uiBrowse
		return m, nil
	case "down", "j":
		m.menuCursor++
		if m.menuCursor >= len(m.menuItems) {
			m.menuCursor = len(m.menuItems) - 1
		}
		return m, nil
	case "up", "k":
		m.menuCursor--
		if m.menuCursor < 0 {
			m.menuCursor = 0
		}
		return m, nil
	case "enter":
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
	switch msg.Type {
	case tea.KeyEsc:
		m.ui = uiBrowse
		m.promptValue = ""
		m.promptKind = ""
		return m, nil
	case tea.KeyEnter:
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
	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.promptValue) > 0 {
			m.promptValue = m.promptValue[:len(m.promptValue)-1]
		}
		return m, nil
	case tea.KeySpace:
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

func (m model) applyAction(a action) (tea.Model, tea.Cmd) {
	switch a.kind {
	case "view_tree":
		m.panel = panelTree
		m.setActiveLines()
		m.cursor, m.offset = 0, 0
		m.status = "Tree view"
		m.clamp()
		return m, nil
	case "view_graph":
		m.panel = panelGraph
		m.refreshTreeAndGraph()
		m.setActiveLines()
		m.cursor, m.offset = 0, 0
		m.status = "Graph view"
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

func (m model) runPromptAction(kind, value string) (tea.Model, tea.Cmd) {
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
	switch m.panel {
	case panelTree:
		m.lines = normalizeLines(m.treeLines)
	case panelGraph:
		m.lines = normalizeLines(m.graphLines)
	case panelOutput:
		m.lines = normalizeLines(m.outputLines)
	default:
		m.lines = []string{"(empty)"}
	}
}

func (m *model) refreshTreeAndGraph() {
	tree, treeErr := buildTreeLines(".")
	graph, graphErr := loadGraphLines()
	m.treeLines = normalizeLines(tree)
	m.graphLines = normalizeLines(graph)

	errs := make([]string, 0, 2)
	if treeErr != nil {
		errs = append(errs, fmt.Sprintf("tree error: %v", treeErr))
	}
	if graphErr != nil {
		errs = append(errs, fmt.Sprintf("graph error: %v", graphErr))
	}
	m.err = strings.Join(errs, " | ")
	m.setActiveLines()
}

func runCommand(title string, name string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(name, args...)
		out, err := cmd.CombinedOutput()
		lines := linesFromOutput(out)
		return cmdResultMsg{title: title, output: lines, err: err}
	}
}

func runShellCommand(title, command string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-lc", command)
		out, err := cmd.CombinedOutput()
		lines := linesFromOutput(out)
		return cmdResultMsg{title: title, output: lines, err: err}
	}
}

func linesFromOutput(out []byte) []string {
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return []string{"(no output)"}
	}
	return strings.Split(string(trimmed), "\n")
}

func normalizeLines(lines []string) []string {
	if len(lines) == 0 {
		return []string{"(empty)"}
	}
	return lines
}

func (m model) View() string {
	head := m.header()
	if m.ui == uiMenu {
		return head + m.renderMenu()
	}
	if m.ui == uiPrompt {
		return head + m.renderPrompt() + "\n\n" + m.renderBody()
	}
	return head + m.renderBody()
}

func (m model) header() string {
	panel := "Tree"
	switch m.panel {
	case panelGraph:
		panel = "Graph"
	case panelOutput:
		panel = "Output"
	}
	header := fmt.Sprintf("nit | Panel: %s | Lines: %d | m menu | tab tree/graph | o output | r reload | q quit\n", panel, len(m.lines))
	if m.status != "" {
		header += "Status: " + m.status + "\n"
	}
	if m.err != "" {
		header += "Warning: " + m.err + "\n"
	}
	header += strings.Repeat("-", max(10, m.width)) + "\n"
	return header
}

func (m model) renderMenu() string {
	var b strings.Builder
	b.WriteString("Actions (enter run, esc cancel):\n")
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
	return fmt.Sprintf("%s: %s\nenter submit | esc cancel", m.promptTitle, value)
}

func (m model) renderBody() string {
	end := min(len(m.lines), m.offset+m.bodyHeight())
	var b strings.Builder
	for i := m.offset; i < end; i++ {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		line := m.lines[i]
		if m.width > 4 && len(line) > m.width-4 {
			line = line[:m.width-7] + "..."
		}
		b.WriteString(prefix + line + "\n")
	}
	return b.String()
}

func (m *model) clamp() {
	if len(m.lines) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.lines) {
		m.cursor = len(m.lines) - 1
	}

	page := m.bodyHeight()
	if page < 1 {
		page = 1
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+page {
		m.offset = m.cursor - page + 1
	}
	maxOffset := max(0, len(m.lines)-page)
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m model) bodyHeight() int {
	used := 3
	if m.status != "" {
		used++
	}
	if m.err != "" {
		used++
	}
	h := m.height - used
	if h < 1 {
		return 1
	}
	return h
}

func (m model) pageSize() int {
	h := m.bodyHeight() - 1
	if h < 1 {
		return 1
	}
	return h
}

func buildTreeLines(root string) ([]string, error) {
	type node struct {
		path  string
		isDir bool
	}
	var nodes []node
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		if strings.Contains(path, string(filepath.Separator)+".git") || strings.HasPrefix(path, ".git") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		nodes = append(nodes, node{path: path, isDir: d.IsDir()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].path == nodes[j].path {
			return nodes[i].isDir && !nodes[j].isDir
		}
		return nodes[i].path < nodes[j].path
	})

	lines := make([]string, 0, len(nodes))
	for _, n := range nodes {
		parts := strings.Split(filepath.Clean(n.path), string(filepath.Separator))
		indent := strings.Repeat("  ", max(0, len(parts)-1))
		name := parts[len(parts)-1]
		if n.isDir {
			lines = append(lines, indent+"[D] "+name+"/")
		} else {
			lines = append(lines, indent+"[F] "+name)
		}
	}
	return lines, nil
}

func loadGraphLines() ([]string, error) {
	cmd := exec.Command("git", "log", "--graph", "--decorate", "--oneline", "--all")
	out, err := cmd.Output()
	if err != nil {
		return []string{"Not a git repo or no commits yet."}, err
	}
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return []string{"No commits to display."}, nil
	}
	return strings.Split(string(out), "\n"), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %v\n", err)
	}
}
