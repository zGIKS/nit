package state

func (s AppState) ChangesPosition() (int, int) {
	total := len(s.Changes.Staged) + len(s.Changes.Unstaged)
	if total < 1 {
		return 1, 1
	}
	cur := 1
	seen := 0
	for i, row := range s.Changes.Rows {
		if row.Selectable {
			seen++
		}
		if i == s.Changes.Cursor {
			if row.Selectable {
				cur = seen
			}
			break
		}
	}
	if cur > total {
		cur = total
	}
	return cur, total
}

func (s AppState) GraphPosition() (int, int) {
	total := len(s.Graph.Lines)
	if total < 1 {
		return 1, 1
	}
	cur := s.Graph.Cursor + 1
	if cur < 1 {
		cur = 1
	}
	if cur > total {
		cur = total
	}
	return cur, total
}
