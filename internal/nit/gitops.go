package nit

import (
	"bytes"
	"os/exec"
	"strings"
)

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

func loadChanges() ([]changeEntry, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimRight(out, "\r\n")
	if len(trimmed) == 0 {
		return []changeEntry{}, nil
	}

	rawLines := strings.Split(string(trimmed), "\n")
	entries := make([]changeEntry, 0, len(rawLines))
	for _, raw := range rawLines {
		e := parseChangeLine(raw)
		entries = append(entries, e)
	}
	return entries, nil
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

func stagePath(path string) error {
	cmd := exec.Command("git", "add", "--", path)
	return cmd.Run()
}

func unstagePath(path string) error {
	// Prefer restore; fallback to reset for older git versions.
	if err := exec.Command("git", "restore", "--staged", "--", path).Run(); err == nil {
		return nil
	}
	return exec.Command("git", "reset", "HEAD", "--", path).Run()
}

func stageAll() error {
	return exec.Command("git", "add", "-A").Run()
}

func unstageAll() error {
	if err := exec.Command("git", "restore", "--staged", ".").Run(); err == nil {
		return nil
	}
	return exec.Command("git", "reset", "HEAD", "--", ".").Run()
}
