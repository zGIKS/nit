package components

import "strings"

func BoxView(title string, width, boxHeight int, lines []string, cursor, offset int, active bool) string {
	w := max(30, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	label := title
	if active {
		label += " *"
	}
	top := "+" + fitText(label, innerW, '-') + "+"
	bottom := "+" + strings.Repeat("-", innerW) + "+"

	var b strings.Builder
	b.WriteString(top + "\n")

	maxOffset := max(0, len(lines)-contentHeight)
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	end := min(len(lines), offset+contentHeight)
	for i := 0; i < contentHeight; i++ {
		idx := offset + i
		text := ""
		if idx < end {
			prefix := "  "
			if idx == cursor {
				prefix = "> "
			}
			text = prefix + lines[idx]
		}
		text = fitText(text, innerW, ' ')
		b.WriteString("|" + text + "|\n")
	}

	b.WriteString(bottom)
	return b.String()
}

func fitText(text string, width int, fill rune) string {
	if width <= 0 {
		return ""
	}
	if len(text) > width {
		if width <= 3 {
			return text[:width]
		}
		return text[:width-3] + "..."
	}
	if len(text) == width {
		return text
	}
	return text + strings.Repeat(string(fill), width-len(text))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
