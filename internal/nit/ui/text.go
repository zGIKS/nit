package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"nit/internal/nit/util"
)

func fitText(text string, width int, fill rune) string {
	if width <= 0 {
		return ""
	}
	textW := displayWidth(text)
	if textW > width {
		if width <= 3 {
			return truncateDisplayWidth(text, width)
		}
		return truncateDisplayWidth(text, width-3) + "..."
	}
	if textW == width {
		return text
	}
	return text + strings.Repeat(string(fill), width-textW)
}

func min(a, b int) int {
	return util.Min(a, b)
}

func max(a, b int) int {
	return util.Max(a, b)
}

func truncateDisplayWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	cur := 0
	for i := 0; i < len(s); {
		if end, ok := ansiSeqEnd(s, i); ok {
			b.WriteString(s[i:end])
			i = end
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		rw := runewidth.RuneWidth(r)
		if rw == 0 {
			b.WriteRune(r)
			i += size
			continue
		}
		if cur+rw > width {
			break
		}
		b.WriteRune(r)
		cur += rw
		i += size
	}
	return b.String()
}

func displayWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		if end, ok := ansiSeqEnd(s, i); ok {
			i = end
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		width += runewidth.RuneWidth(r)
		i += size
	}
	return width
}

func ansiSeqEnd(s string, i int) (int, bool) {
	if i+1 >= len(s) || s[i] != 0x1b || s[i+1] != '[' {
		return 0, false
	}
	j := i + 2
	for j < len(s) {
		c := s[j]
		// CSI final byte range.
		if c >= 0x40 && c <= 0x7e {
			return j + 1, true
		}
		j++
	}
	return 0, false
}

func ansiUnderline(s string) string {
	return fmt.Sprintf("\x1b[4m%s\x1b[24m", s)
}
