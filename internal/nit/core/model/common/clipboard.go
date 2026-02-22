package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"nit/internal/nit/config"
)

func CopyWithMode(cfg config.ClipboardConfig, text string) error {
	if text == "" {
		return nil
	}
	var lastErr error
	switch cfg.Mode {
	case config.ClipboardInternal:
		return nil
	case config.ClipboardOnlyCopy:
		if err := copyWithOSC52(text); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		} else {
			lastErr = err
		}
		return nil
	case config.ClipboardOSC52:
		if err := copyWithOSC52(text); err == nil {
			return nil
		} else {
			lastErr = err
		}
		return lastErr
	case config.ClipboardSystem:
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		} else {
			lastErr = err
		}
		return lastErr
	default:
		if err := copyWithOSC52(text); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		} else {
			lastErr = err
		}
		return nil
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

func clipboardCopyCommand() (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbcopy"), nil
	case "windows":
		return exec.Command("cmd", "/c", "clip"), nil
	default:
		if p, _ := exec.LookPath("wl-copy"); p != "" {
			return exec.Command(p), nil
		}
		if p, _ := exec.LookPath("xclip"); p != "" {
			return exec.Command(p, "-selection", "clipboard"), nil
		}
		if p, _ := exec.LookPath("xsel"); p != "" {
			return exec.Command(p, "--clipboard", "--input"), nil
		}
	}
	return nil, fmt.Errorf("no clipboard copy tool found (macOS: pbcopy, Windows: clip, Linux/NixOS: wl-clipboard/xclip/xsel)")
}

func clipboardPasteCommand() (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbpaste"), nil
	case "windows":
		if p, _ := exec.LookPath("powershell"); p != "" {
			return exec.Command(p, "-NoProfile", "-Command", "Get-Clipboard"), nil
		}
		if p, _ := exec.LookPath("pwsh"); p != "" {
			return exec.Command(p, "-NoProfile", "-Command", "Get-Clipboard"), nil
		}
	default:
		if p, _ := exec.LookPath("wl-paste"); p != "" {
			return exec.Command(p, "-n"), nil
		}
		if p, _ := exec.LookPath("xclip"); p != "" {
			return exec.Command(p, "-selection", "clipboard", "-o"), nil
		}
		if p, _ := exec.LookPath("xsel"); p != "" {
			return exec.Command(p, "--clipboard", "--output"), nil
		}
	}
	return nil, fmt.Errorf("no clipboard paste tool found (macOS: pbpaste, Windows: powershell Get-Clipboard, Linux/NixOS: wl-clipboard/xclip/xsel)")
}
