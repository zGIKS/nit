package git

import (
	"errors"
	"strconv"
	"strings"
)

type CommitOptions struct {
	All     bool
	Amend   bool
	Signoff bool
}

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

func (s Service) DiscardAll() (string, error) {
	cmdLog := ""
	resetErr := error(nil)
	if _, resetCmd, err := s.runner.Run("reset", "--hard", "HEAD"); err == nil {
		cmdLog = resetCmd
	} else {
		resetErr = err
		if resetCmd != "" {
			cmdLog = resetCmd
		}
	}

	_, cleanCmd, cleanErr := s.runner.Run("clean", "-fd")
	if cleanCmd != "" {
		if cmdLog != "" {
			cmdLog += " && " + cleanCmd
		} else {
			cmdLog = cleanCmd
		}
	}

	// In repos without commits yet, reset HEAD fails. Cleaning untracked files can still
	// fully discard the working tree, so treat that case as success when clean works.
	if resetErr != nil && cleanErr == nil {
		errText := strings.ToLower(resetErr.Error())
		if strings.Contains(errText, "ambiguous argument 'head'") ||
			strings.Contains(errText, "unknown revision") ||
			strings.Contains(errText, "bad revision 'head'") {
			return cmdLog, nil
		}
	}
	if resetErr != nil {
		return cmdLog, resetErr
	}
	return cmdLog, cleanErr
}

func (s Service) Commit(message string) (string, error) {
	return s.CommitWithOptions(message, CommitOptions{})
}

func (s Service) CommitWithOptions(message string, opts CommitOptions) (string, error) {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return "", errors.New("commit message is empty")
	}

	cmdLog := ""
	if opts.All {
		_, addCmd, err := s.runner.Run("add", "-A")
		cmdLog = addCmd
		if err != nil {
			return cmdLog, err
		}
	}

	args := []string{"commit"}
	if opts.Amend {
		args = append(args, "--amend")
	}
	if opts.Signoff {
		args = append(args, "--signoff")
	}
	args = append(args, "-m", msg)
	_, commitCmd, err := s.runner.Run(args...)
	if cmdLog != "" && commitCmd != "" {
		cmdLog += " && " + commitCmd
	} else if commitCmd != "" {
		cmdLog = commitCmd
	}
	return cmdLog, err
}

func (s Service) UndoLastCommit() (string, error) {
	_, cmd, err := s.runner.Run("reset", "--soft", "HEAD~1")
	return cmd, err
}

func (s Service) AbortRebase() (string, error) {
	_, cmd, err := s.runner.Run("rebase", "--abort")
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
