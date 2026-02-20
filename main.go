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

type viewMode int

const (
	viewTree viewMode = iota
	viewGraph
)

type model struct {
	mode      viewMode
	treeLines []string
	graphLines []string
	lines     []string
	cursor    int
	offset    int
	width     int
	height    int
	err       string
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
		mode:       viewTree,
		treeLines:  tree,
		graphLines: graph,
		lines:      tree,
		height:     24,
		err:        strings.Join(errs, " | "),
	}
	if len(m.lines) == 0 {
		m.lines = []string{"(empty)"}
	}
	return m
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clamp()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.mode == viewTree {
				m.mode = viewGraph
				m.lines = m.graphLines
			} else {
				m.mode = viewTree
				m.lines = m.treeLines
			}
			if len(m.lines) == 0 {
				m.lines = []string{"(empty)"}
			}
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
			tree, treeErr := buildTreeLines(".")
			graph, graphErr := loadGraphLines()
			m.treeLines = tree
			m.graphLines = graph
			if m.mode == viewTree {
				m.lines = m.treeLines
			} else {
				m.lines = m.graphLines
			}
			if len(m.lines) == 0 {
				m.lines = []string{"(empty)"}
			}
			errs := make([]string, 0, 2)
			if treeErr != nil {
				errs = append(errs, fmt.Sprintf("tree error: %v", treeErr))
			}
			if graphErr != nil {
				errs = append(errs, fmt.Sprintf("graph error: %v", graphErr))
			}
			m.err = strings.Join(errs, " | ")
			m.clamp()
		}
	}
	return m, nil
}

func (m model) View() string {
	title := "Tree"
	if m.mode == viewGraph {
		title = "Graph"
	}
	header := fmt.Sprintf(
		"Git Navigator | View: %s | Lines: %d | tab switch | r reload | q quit\n",
		title, len(m.lines),
	)
	if m.err != "" {
		header += "Warning: " + m.err + "\n"
	}
	header += strings.Repeat("-", max(10, m.width)) + "\n"

	bodyHeight := m.bodyHeight()
	end := min(len(m.lines), m.offset+bodyHeight)
	var b strings.Builder
	b.WriteString(header)
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
