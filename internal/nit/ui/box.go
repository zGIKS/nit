package ui

import "strings"

func BoxView(title string, width, boxHeight int, lines []string, cursor, offset int, active bool, footer string) string {
	w := max(8, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	head := title
	if active {
		head = "● " + head
	}
	top := "┌" + fitText(" "+head+" ", innerW, '─') + "┐"

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
				prefix = "▌ "
			}
			text = prefix + lines[idx]
		}
		text = fitText(text, innerW-2, ' ')
		b.WriteString("│ " + text + " │\n")
	}

	bottom := "└" + fitText(" "+footer+" ", innerW, '─') + "┘"
	b.WriteString(bottom)
	return b.String()
}

func HStack(left string, leftWidth int, right string, rightWidth int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	h := max(len(leftLines), len(rightLines))

	var b strings.Builder
	for i := 0; i < h; i++ {
		l := ""
		if i < len(leftLines) {
			l = fitText(leftLines[i], leftWidth, ' ')
		} else {
			l = strings.Repeat(" ", leftWidth)
		}
		r := ""
		if i < len(rightLines) {
			r = fitText(rightLines[i], rightWidth, ' ')
		} else {
			r = strings.Repeat(" ", rightWidth)
		}
		b.WriteString(l + " " + r)
		if i != h-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func HStackMany(parts []string, widths []int) string {
	if len(parts) == 0 || len(parts) != len(widths) {
		return ""
	}
	out := parts[0]
	outW := widths[0]
	for i := 1; i < len(parts); i++ {
		out = HStack(out, outW, parts[i], widths[i])
		outW += 1 + widths[i]
	}
	return out
}
