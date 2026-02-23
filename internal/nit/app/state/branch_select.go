package state

import "strings"

func (s AppState) SelectedBranchName() (string, bool) {
	if s.Branches.Cursor < 0 || s.Branches.Cursor >= len(s.Branches.Lines) {
		return "", false
	}
	name := normalizeBranchListLine(s.Branches.Lines[s.Branches.Cursor])
	if name == "" {
		return "", false
	}
	return name, true
}

func normalizeBranchListLine(line string) string {
	name := strings.TrimSpace(line)
	if name == "" {
		return ""
	}
	if strings.Contains(name, "No local branches") || strings.Contains(name, "Loading branches") || strings.Contains(name, "Not a git repo") {
		return ""
	}
	name = strings.TrimPrefix(name, "●")
	name = strings.TrimPrefix(name, "*")
	return strings.TrimSpace(name)
}
