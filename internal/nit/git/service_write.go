package git

import (
"errors"
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
return s.runWithFallback(
[]string{"restore", "--staged", "--", path},
[]string{"reset", "HEAD", "--", path},
)
}

func (s Service) StageAll() (string, error) {
_, cmd, err := s.runner.Run("add", "-A")
return cmd, err
}

func (s Service) UnstageAll() (string, error) {
return s.runWithFallback(
[]string{"restore", "--staged", "--", "."},
[]string{"reset", "HEAD", "--", "."},
)
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

// runWithFallback tries primaryArgs first, then falls back to fallbackArgs.
func (s Service) runWithFallback(primaryArgs, fallbackArgs []string) (string, error) {
if _, cmd, err := s.runner.Run(primaryArgs...); err == nil {
return cmd, nil
}
_, cmd, err := s.runner.Run(fallbackArgs...)
return cmd, err
}
