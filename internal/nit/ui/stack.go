package ui

import (
	"strings"
)

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
