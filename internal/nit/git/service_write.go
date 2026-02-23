package git

import (
	"errors"
	"strconv"
	"strings"
)

func (s Service) StagePath(path string) (string, error) {
	_, cmd, err := s.runner.Run("add", "--", path)
	return cmd, err
}

func (s Service) UnstagePath(path string) (string, error) {
	return s.unstageWithFallback(path)
}

func (s Service) StageAll() (string, error) {
	_, cmd, err := s.runner.Run("add", "-A")
	return cmd, err
}

func (s Service) UnstageAll() (string, error) {
	return s.unstageWithFallback(".")
}

func (s Service) Commit(message string) (string, error) {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return "", errors.New("commit message is empty")
	}
	_, cmd, err := s.runner.Run("commit", "-m", msg)
	return cmd, err
}

func (s Service) CreateBranch(name, source string) (string, error) {
	branch := strings.TrimSpace(name)
	if branch == "" {
		return "", errors.New("branch name is empty")
	}
	src := strings.TrimSpace(source)
	if src == "" || src == "-" {
		if _, cmd, err := s.runner.Run("switch", "-c", branch); err == nil {
			return cmd, nil
		}
		_, cmd, err := s.runner.Run("checkout", "-b", branch)
		return cmd, err
	}
	if _, cmd, err := s.runner.Run("switch", "-c", branch, src); err == nil {
		return cmd, nil
	}
	_, cmd, err := s.runner.Run("checkout", "-b", branch, src)
	return cmd, err
}

func (s Service) SwitchBranch(name string) (string, error) {
	branch := strings.TrimSpace(name)
	if branch == "" {
		return "", errors.New("branch name is empty")
	}
	if _, cmd, err := s.runner.Run("switch", branch); err == nil {
		return cmd, nil
	}
	_, cmd, err := s.runner.Run("checkout", branch)
	return cmd, err
}

func (s Service) PushCurrentBranchUpstream() (string, error) {
	_, cmd, err := s.runner.Run("push", "-u", "origin", "HEAD")
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

func (s Service) unstageWithFallback(target string) (string, error) {
	if _, cmd, err := s.runner.Run("restore", "--staged", "--", target); err == nil {
		return cmd, nil
	}
	_, cmd, err := s.runner.Run("reset", "HEAD", "--", target)
	return cmd, err
}
