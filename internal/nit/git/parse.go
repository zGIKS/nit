package git

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseChangeLine(raw string) ChangeEntry {
	e := ChangeEntry{Raw: raw}
	if len(raw) < 3 {
		e.Path = raw
		return e
	}
	e.X = raw[0]
	e.Y = raw[1]
	path := strings.TrimSpace(raw[3:])
	path = decodeCQuoted(path)
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

func decodeCQuoted(s string) string {
	if !strings.HasPrefix(s, "\"") || !strings.HasSuffix(s, "\"") {
		return s
	}
	s = s[1 : len(s)-1]
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+3 < len(s) {
			if octal := s[i+1 : i+4]; isOctal(octal) {
				v, _ := strconv.ParseUint(octal, 8, 8)
				result.WriteByte(byte(v))
				i += 4
				continue
			}
		}
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
				i += 2
				continue
			case 't':
				result.WriteByte('\t')
				i += 2
				continue
			case '\\':
				result.WriteByte('\\')
				i += 2
				continue
			case '"':
				result.WriteByte('"')
				i += 2
				continue
			}
		}
		result.WriteByte(s[i])
		i++
	}
	return result.String()
}

func isOctal(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) || c > '7' {
			return false
		}
	}
	return true
}
