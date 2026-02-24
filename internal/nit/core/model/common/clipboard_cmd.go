package common

import (
	"fmt"
	"os/exec"
	"runtime"
)

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
