package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type clipboardMode string

const (
	clipboardOnlyCopy clipboardMode = "only_copy"
	clipboardAuto     clipboardMode = "auto"
	clipboardOSC52    clipboardMode = "osc52"
	clipboardSystem   clipboardMode = "system"
	clipboardInternal clipboardMode = "internal"
)

type clipboardConfig struct {
	Mode     clipboardMode
	CopyCmd  string
	PasteCmd string
}

func loadClipboardConfig() clipboardConfig {
	mode := clipboardMode(strings.ToLower(strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_MODE"))))
	switch mode {
	case clipboardOnlyCopy, clipboardOSC52, clipboardSystem, clipboardInternal, clipboardAuto:
	default:
		mode = clipboardOnlyCopy
	}
	return clipboardConfig{
		Mode:     mode,
		CopyCmd:  strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_COPY_CMD")),
		PasteCmd: strings.TrimSpace(os.Getenv("NIT_CLIPBOARD_PASTE_CMD")),
	}
}

func copyWithMode(cfg clipboardConfig, text string) error {
	if text == "" {
		return nil
	}
	switch cfg.Mode {
	case clipboardInternal:
		return nil
	case clipboardOnlyCopy:
		if err := copyWithOSC52(text); err == nil {
			return nil
		}
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		}
		return nil
	case clipboardOSC52:
		if err := copyWithOSC52(text); err == nil {
			return nil
		}
		return nil
	case clipboardSystem:
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		}
		return nil
	default:
		if err := copyWithOSC52(text); err == nil {
			return nil
		}
		if err := copyToSystemClipboard(text, cfg.CopyCmd); err == nil {
			return nil
		}
		return nil
	}
}

func pasteWithMode(cfg clipboardConfig) (string, error) {
	switch cfg.Mode {
	case clipboardInternal, clipboardOSC52, clipboardOnlyCopy:
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
		// Some commands fail with non-zero when clipboard is empty.
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
	return nil, fmt.Errorf("no clipboard copy tool found (tried wl-copy/xclip/xsel/pbcopy/clip)")
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
	return nil, fmt.Errorf("no clipboard paste tool found (tried wl-paste/xclip/xsel/pbpaste/powershell)")
}
