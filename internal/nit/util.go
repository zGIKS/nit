package nit

func normalizeLines(lines []string) []string {
	if len(lines) == 0 {
		return []string{"(empty)"}
	}
	return lines
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
