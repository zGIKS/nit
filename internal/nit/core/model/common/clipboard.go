package common

import (
"bytes"
"encoding/base64"
"fmt"
"os"
"os/exec"
"strings"

"github.com/zGIKS/nit/internal/nit/config"
)

// tryFuncs runs functions in order, returning nil on first success.
func tryFuncs(fns ...func() error) error {
var lastErr error
for _, fn := range fns {
if err := fn(); err == nil {
return nil
} else {
lastErr = err
}
}
return lastErr
}

func CopyWithMode(cfg config.ClipboardConfig, text string) error {
if text == "" {
return nil
}
osc := func() error { return copyWithOSC52(text) }
sys := func() error { return copyToSystemClipboard(text, cfg.CopyCmd) }

switch cfg.Mode {
case config.ClipboardInternal:
return nil
case config.ClipboardOSC52:
return tryFuncs(osc)
case config.ClipboardSystem:
return tryFuncs(sys)
case config.ClipboardOnlyCopy:
return tryFuncs(osc, sys)
default:
return tryFuncs(osc, sys)
}
}

func PasteWithMode(cfg config.ClipboardConfig) (string, error) {
switch cfg.Mode {
case config.ClipboardInternal, config.ClipboardOSC52, config.ClipboardOnlyCopy:
return "", nil
default:
return pasteFromSystemClipboard(cfg.PasteCmd)
}
}

func copyWithOSC52(text string) error {
encoded := base64.StdEncoding.EncodeToString([]byte(text))
seq := "\x1b]52;c;" + encoded + "\a"
_, err := os.Stdout.WriteString(seq)
return err
}

func copyToSystemClipboard(text, customCmd string) error {
if customCmd != "" {
cmd := exec.Command("sh", "-lc", customCmd)
cmd.Stdin = strings.NewReader(text)
if out, runErr := cmd.CombinedOutput(); runErr != nil {
return fmt.Errorf("clipboard copy failed: %s", strings.TrimSpace(string(out)))
}
return nil
}
cmd, err := clipboardCopyCommand()
if err != nil {
return err
}
cmd.Stdin = strings.NewReader(text)
if out, runErr := cmd.CombinedOutput(); runErr != nil {
return fmt.Errorf("clipboard copy failed: %s", strings.TrimSpace(string(out)))
}
return nil
}

func pasteFromSystemClipboard(customCmd string) (string, error) {
if customCmd != "" {
cmd := exec.Command("sh", "-lc", customCmd)
out, runErr := cmd.Output()
if runErr != nil {
return "", fmt.Errorf("clipboard paste failed: %w", runErr)
}
return string(bytes.TrimRight(out, "\r\n")), nil
}
cmd, err := clipboardPasteCommand()
if err != nil {
return "", err
}
out, runErr := cmd.Output()
if runErr != nil {
if ee, ok := runErr.(*exec.ExitError); ok && len(ee.Stderr) == 0 {
return "", nil
}
return "", fmt.Errorf("clipboard paste failed: %w", runErr)
}
return string(bytes.TrimRight(out, "\r\n")), nil
}
