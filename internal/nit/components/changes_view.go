package components

import "fmt"

func ChangesView(width, height int, lines []string, cursor, offset int, active bool, selected, total int) string {
	footer := fmt.Sprintf("%d of %d", selected, total)
	return BoxView("Changes - Staged - Unstaged", width, height, lines, cursor, offset, active, footer)
}
