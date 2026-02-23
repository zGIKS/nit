package git

import (
	"path/filepath"
	"strings"
)

type Service struct {
	runner Runner
}

func NewService(r Runner) Service {
	return Service{runner: r}
}

func (s Service) LoadGraph() ([]string, error) {
	out, _, err := s.runner.Run("--no-optional-locks", "log", "--graph", "--decorate", "--oneline", "--all")
	if err != nil {
		return []string{"Not a git repo or no commits yet."}, err
	}
	if strings.TrimSpace(out) == "" {
		return []string{"No commits to display."}, nil
	}
	lines := strings.Split(out, "\n")
	for i := range lines {
		lines[i] = prettifyGraphLine(lines[i])
	}
	return lines, nil
}

func (s Service) LoadBranches() ([]string, error) {
	out, _, err := s.runner.Run(
		"--no-optional-locks",
		"for-each-ref",
		"--format=%(HEAD) %(refname:short)",
		"refs/heads",
	)
	if err != nil {
		return []string{"Not a git repo."}, err
	}
	if strings.TrimSpace(out) == "" {
		return []string{"No local branches."}, nil
	}
	raw := strings.Split(out, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimRight(line, " \t")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			lines = append(lines, "● "+strings.TrimSpace(strings.TrimPrefix(line, "* ")))
			continue
		}
		if strings.HasPrefix(line, "  ") {
			lines = append(lines, "  "+strings.TrimSpace(line[2:]))
			continue
		}
		lines = append(lines, "  "+strings.TrimSpace(line))
	}
	if len(lines) == 0 {
		return []string{"No local branches."}, nil
	}
	return lines, nil
}

func (s Service) LoadChanges() ([]ChangeEntry, error) {
	out, _, err := s.runner.Run("--no-optional-locks", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(out) == "" {
		return []ChangeEntry{}, nil
	}

	rawLines := strings.Split(out, "\n")
	entries := make([]ChangeEntry, 0, len(rawLines))
	for _, raw := range rawLines {
		e := ParseChangeLine(raw)
		entries = append(entries, e)
	}
	return entries, nil
}

func (s Service) LoadRepoSummary() (string, string, error) {
	root, _, err := s.runner.Run("--no-optional-locks", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}
	repo := filepath.Base(strings.TrimSpace(root))
	branch, _, err := s.runner.Run("--no-optional-locks", "branch", "--show-current")
	if err != nil {
		return repo, "", err
	}
	br := strings.TrimSpace(branch)
	if br == "" {
		br = "(detached)"
	}
	return repo, br, nil
}
