package git

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
	runner Runner
}

var graphHashRe = regexp.MustCompile(`[0-9a-f]{7,40}\b`)
var graphPrefixRe = regexp.MustCompile(`^[|\\/*_. ]+`)

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

func (s Service) StagePath(path string) (string, error) {
	_, cmd, err := s.runner.Run("add", "--", path)
	return cmd, err
}

func (s Service) UnstagePath(path string) (string, error) {
	if _, cmd, err := s.runner.Run("restore", "--staged", "--", path); err == nil {
		return cmd, nil
	}
	_, cmd, err := s.runner.Run("reset", "HEAD", "--", path)
	return cmd, err
}

func (s Service) StageAll() (string, error) {
	_, cmd, err := s.runner.Run("add", "-A")
	return cmd, err
}

func (s Service) UnstageAll() (string, error) {
	if _, cmd, err := s.runner.Run("restore", "--staged", "."); err == nil {
		return cmd, nil
	}
	_, cmd, err := s.runner.Run("reset", "HEAD", "--", ".")
	return cmd, err
}

func (s Service) Commit(message string) (string, error) {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return "", errors.New("commit message is empty")
	}
	_, cmd, err := s.runner.Run("commit", "-m", msg)
	return cmd, err
}

func (s Service) Pull() (string, error) {
	_, cmd, err := s.runner.Run("pull")
	return cmd, err
}

func (s Service) Push() (string, error) {
	if err := s.ensureHasOutgoingCommits(); err != nil {
		return "", err
	}
	_, cmd, err := s.runner.Run("push")
	return cmd, err
}

func (s Service) Fetch() (string, error) {
	_, cmd, err := s.runner.Run("fetch")
	return cmd, err
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

func (s Service) ensureHasOutgoingCommits() error {
	// If there is no upstream configured, allow push so users can publish/set upstream.
	if _, _, err := s.runner.Run("--no-optional-locks", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"); err != nil {
		return nil
	}
	out, _, err := s.runner.Run("--no-optional-locks", "rev-list", "--count", "@{u}..HEAD")
	if err != nil {
		return nil
	}
	n, convErr := strconv.Atoi(strings.TrimSpace(out))
	if convErr != nil {
		return nil
	}
	if n == 0 {
		return errors.New("nothing to push")
	}
	return nil
}

func prettifyGraphLine(line string) string {
	if line == "" {
		return line
	}
	prefixEnd := 0
	if loc := graphHashRe.FindStringIndex(line); loc != nil && loc[0] > 0 {
		prefixEnd = loc[0]
	} else if loc := graphPrefixRe.FindStringIndex(line); loc != nil {
		prefixEnd = loc[1]
	}
	if prefixEnd <= 0 {
		return line
	}
	return replaceGraphChars(line[:prefixEnd]) + line[prefixEnd:]
}

func replaceGraphChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '|':
			b.WriteRune('│')
		case '/':
			b.WriteRune('╱')
		case '\\':
			b.WriteRune('╲')
		case '*':
			b.WriteRune('●')
		case '_':
			b.WriteRune('─')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
