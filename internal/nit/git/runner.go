package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Runner struct {
	Timeout time.Duration
	GitPath string
}

func NewRunner(timeout time.Duration) Runner {
	if timeout <= 0 {
		timeout = 4 * time.Second
	}
	gitPath, _ := exec.LookPath("git")
	return Runner{Timeout: timeout, GitPath: gitPath}
}

func (r Runner) Run(args ...string) (string, string, error) {
	cmdStr := "git " + strings.Join(args, " ")
	if strings.TrimSpace(r.GitPath) == "" {
		return "", cmdStr, fmt.Errorf("git executable not found in PATH")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, r.GitPath, args...)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	err := cmd.Run()
	stdout := strings.TrimRight(out.String(), "\r\n")
	stderr := strings.TrimSpace(errBuf.String())
	if err == nil {
		return stdout, cmdStr, nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return stdout, cmdStr, fmt.Errorf("git %s timeout after %s", strings.Join(args, " "), r.Timeout)
	}
	if stderr != "" {
		return stdout, cmdStr, fmt.Errorf("git %s failed: %s", strings.Join(args, " "), stderr)
	}
	if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
		return stdout, cmdStr, fmt.Errorf("git executable not found in PATH")
	}
	return stdout, cmdStr, fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
}
