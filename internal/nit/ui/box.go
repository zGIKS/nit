package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

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

func TopBarView(width int, left, right string) string {
	w := max(8, width)
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	leftLen := runewidth.StringWidth(left)
	rightLen := runewidth.StringWidth(right)

	if rightLen >= w {
		return fitText(right, w, ' ')
	}

	space := w - leftLen - rightLen
	if space < 1 {
		maxLeft := w - rightLen - 1
		if maxLeft < 0 {
			maxLeft = 0
		}
		left = fitText(left, maxLeft, ' ')
		left = strings.TrimRight(left, " ")
		space = w - runewidth.StringWidth(left) - rightLen
		if space < 1 {
			space = 1
		}
	}

	return left + strings.Repeat(" ", space) + right
}

func MiniBoxView(text string, width int) string {
	w := max(8, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	mid := "│" + fitText(" "+strings.TrimSpace(text)+" ", innerW, ' ') + "│"
	bot := "└" + strings.Repeat("─", innerW) + "┘"
	return top + "\n" + mid + "\n" + bot
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

func fitText(text string, width int, fill rune) string {
	if width <= 0 {
		return ""
	}
	textW := runewidth.StringWidth(text)
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

func truncateDisplayWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	cur := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if rw == 0 {
			b.WriteRune(r)
			continue
		}
		if cur+rw > width {
			break
		}
		b.WriteRune(r)
		cur += rw
	}
	return b.String()
}
