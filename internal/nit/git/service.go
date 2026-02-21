package git

import (
	"errors"
	"strings"
)

type Service struct {
	runner Runner
}

func NewService(r Runner) Service {
	return Service{runner: r}
}

func (s Service) LoadGraph() ([]string, error) {
	out, err := s.runner.Run("log", "--graph", "--decorate", "--oneline", "--all")
	if err != nil {
		return []string{"Not a git repo or no commits yet."}, err
	}
	if strings.TrimSpace(out) == "" {
		return []string{"No commits to display."}, nil
	}
	return strings.Split(out, "\n"), nil
}

func (s Service) LoadChanges() ([]ChangeEntry, error) {
	out, err := s.runner.Run("status", "--porcelain")
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

func (s Service) StagePath(path string) error {
	_, err := s.runner.Run("add", "--", path)
	return err
}

func (s Service) UnstagePath(path string) error {
	if _, err := s.runner.Run("restore", "--staged", "--", path); err == nil {
		return nil
	}
	_, err := s.runner.Run("reset", "HEAD", "--", path)
	return err
}

func (s Service) StageAll() error {
	_, err := s.runner.Run("add", "-A")
	return err
}

func (s Service) UnstageAll() error {
	if _, err := s.runner.Run("restore", "--staged", "."); err == nil {
		return nil
	}
	_, err := s.runner.Run("reset", "HEAD", "--", ".")
	return err
}

func (s Service) Commit(message string) error {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return errors.New("commit message is empty")
	}
	_, err := s.runner.Run("commit", "-m", msg)
	return err
}

func (s Service) Push() error {
	_, err := s.runner.Run("push")
	return err
}
