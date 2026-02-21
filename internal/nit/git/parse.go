package git

import "strings"

func ParseChangeLine(raw string) ChangeEntry {
	e := ChangeEntry{Raw: raw}
	if len(raw) < 3 {
		e.Path = raw
		return e
	}
	e.X = raw[0]
	e.Y = raw[1]
	path := strings.TrimSpace(raw[3:])
	// For rename/copy entries porcelain uses: "old/path -> new/path".
	// Git path operations must target the destination path.
	if (e.X == 'R' || e.X == 'C' || e.Y == 'R' || e.Y == 'C') && strings.Contains(path, " -> ") {
		parts := strings.Split(path, " -> ")
		path = strings.TrimSpace(parts[len(parts)-1])
	}
	e.Path = path
	e.Staged = e.X != ' ' && e.X != '?'
	e.Changed = e.Y != ' ' || e.X == '?'
	return e
}
