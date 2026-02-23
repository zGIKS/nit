package git

import (
	"regexp"
	"strings"
)

var graphHashRe = regexp.MustCompile(`[0-9a-f]{7,40}\b`)
var graphPrefixRe = regexp.MustCompile(`^[|\\/*_. ]+`)

func prettifyGraphLine(line string) string {
	if line == "" {
		return line
	}
	prefixEnd := 0
	if loc := graphHashRe.FindStringIndex(line); loc != nil && loc[0] > 0 {
		prefixEnd = loc[0]
	} else if loc := graphPrefixRe.FindStringIndex(line); loc != nil {
		prefixEnd = loc[1]
	}
	if prefixEnd <= 0 {
		return line
	}
	return replaceGraphChars(line[:prefixEnd]) + line[prefixEnd:]
}

func replaceGraphChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '|':
			b.WriteRune('│')
		case '/':
			b.WriteRune('╱')
		case '\\':
			b.WriteRune('╲')
		case '*':
			b.WriteRune('●')
		case '_':
			b.WriteRune('─')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
