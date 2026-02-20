package nit

import (
	"bytes"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func runCommand(title string, name string, args ...string) tea.Cmd {
	return runCommandWithOutputMode(title, true, name, args...)
}

func runCommandWithOutputMode(title string, switchToOutput bool, name string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(name, args...)
		out, err := cmd.CombinedOutput()
		lines := linesFromOutput(out)
		return cmdResultMsg{title: title, output: lines, err: err, switchToOutput: switchToOutput}
	}
}

func runShellCommand(title, command string) tea.Cmd {
	return runShellCommandWithOutputMode(title, true, command)
}

func runShellCommandWithOutputMode(title string, switchToOutput bool, command string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-lc", command)
		out, err := cmd.CombinedOutput()
		lines := linesFromOutput(out)
		return cmdResultMsg{title: title, output: lines, err: err, switchToOutput: switchToOutput}
	}
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

func loadChanges() ([]changeEntry, []string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, []string{"Unable to load changes."}, err
	}

	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return []changeEntry{}, []string{"Working tree clean."}, nil
	}

	rawLines := strings.Split(string(trimmed), "\n")
	entries := make([]changeEntry, 0, len(rawLines))
	lines := make([]string, 0, len(rawLines))
	for _, raw := range rawLines {
		e := parseChangeLine(raw)
		entries = append(entries, e)
		lines = append(lines, formatChangeLine(e))
	}
	return entries, lines, nil
}

func parseChangeLine(raw string) changeEntry {
	e := changeEntry{raw: raw}
	if len(raw) < 3 {
		e.path = raw
		return e
	}
	e.x = raw[0]
	e.y = raw[1]
	e.path = strings.TrimSpace(raw[3:])
	e.staged = e.x != ' ' && e.x != '?'
	e.changed = e.y != ' ' || e.x == '?'
	return e
}

func formatChangeLine(e changeEntry) string {
	stagedMark := " "
	if e.staged {
		stagedMark = "S"
	}
	changedMark := " "
	if e.changed {
		changedMark = "M"
	}
	return "[" + stagedMark + "][" + changedMark + "] " + e.path
}

func linesFromOutput(out []byte) []string {
	trimmed := bytes.TrimSpace(out)
	if len(trimmed) == 0 {
		return []string{"(no output)"}
	}
	return strings.Split(string(trimmed), "\n")
}
