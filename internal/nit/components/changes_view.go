package components

func ChangesView(width, height int, lines []string, cursor, offset int, active bool) string {
	return BoxView("Changes", width, height, lines, cursor, offset, active)
}
