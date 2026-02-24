package git

import (
	"errors"
	"strconv"
	"strings"
)

func (s Service) CreateBranch(name, source string) (string, error) {
	branch := strings.TrimSpace(name)
	if branch == "" {
		return "", errors.New("branch name is empty")
	}
	src := strings.TrimSpace(source)
	if src == "" || src == "-" {
		return s.runWithFallback(
			[]string{"switch", "-c", branch},
			[]string{"checkout", "-b", branch},
		)
	}
	return s.runWithFallback(
		[]string{"switch", "-c", branch, src},
		[]string{"checkout", "-b", branch, src},
	)
}

func (s Service) SwitchBranch(name string) (string, error) {
	branch := strings.TrimSpace(name)
	if branch == "" {
		return "", errors.New("branch name is empty")
	}
	return s.runWithFallback(
		[]string{"switch", branch},
		[]string{"checkout", branch},
	)
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
