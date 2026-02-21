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
}

func NewRunner(timeout time.Duration) Runner {
	if timeout <= 0 {
		timeout = 4 * time.Second
	}
	return Runner{Timeout: timeout}
}

func (r Runner) Run(args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	cmdStr := "git " + strings.Join(args, " ")
	cmd := exec.CommandContext(ctx, "git", args...)
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
	return stdout, cmdStr, fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
}
