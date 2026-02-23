package ui

import "github.com/zGIKS/nit/internal/nit/app"

func commitContentWidth(totalWidth int) int {
	totalW := max(40, totalWidth)
	pushW := max(18, totalW/4)
	commitW := totalW - pushW - 1
	if commitW < 20 {
		commitW = 20
	}
	// BoxView visible width for content line is (w-4), but it also prepends
	// a 2-char cursor prefix ("▌ " or "  "), so the user text gets (w-6).
	return commitW - 6
}

func commandLineViewport(state app.AppState, width int) string {
	return textInputViewport(state.Command.Input, state.Command.Cursor, state.Command.SelectAll, width)
}

func textInputViewport(value string, cursor int, selectAll bool, width int) string {
	full := textInputLineWithCaret(value, cursor, selectAll)
	if width < 4 {
		return full
	}
	caret := cursor
	if selectAll {
		caret = len([]rune(value))
	}
	r := []rune(full)
	if len(r) <= width {
		return full
	}

	start := caret - width/2
	if start < 0 {
		start = 0
	}
	if start+width > len(r) {
		start = len(r) - width
	}
	if start < 0 {
		start = 0
	}
	end := min(len(r), start+width)
	return string(r[start:end])
}

func textInputLineWithCaret(value string, cursor int, selectAll bool) string {
	if selectAll && value != "" {
		return "[" + value + "]"
	}
	r := []rune(value)
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(r) {
		cursor = len(r)
	}
	out := make([]rune, 0, len(r)+1)
	out = append(out, r[:cursor]...)
	out = append(out, '|')
	out = append(out, r[cursor:]...)
	return string(out)
}
