package components

func GraphView(width, height int, lines []string, cursor, offset int, active bool) string {
	return BoxView("Graph", width, height, lines, cursor, offset, active)
}
