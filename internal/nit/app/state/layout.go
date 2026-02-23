package state

func (s AppState) GraphBranchesPaneWidths() (graphW, branchW int) {
	totalW := max(40, s.Viewport.Width)
	branchW = max(24, totalW/3)
	if branchW > totalW-20 {
		branchW = max(18, totalW-20)
	}
	graphW = totalW - branchW - 1
	if graphW < 20 {
		graphW = 20
		branchW = max(18, totalW-graphW-1)
	}
	return graphW, branchW
}

func (s *AppState) Clamp() {
	if s.Focus == FocusGraph {
		clampScrollView(len(s.Graph.Lines), &s.Graph.Cursor, &s.Graph.Offset, s.graphPageSize())
		return
	}

	if s.Focus == FocusBranches {
		clampScrollView(len(s.Branches.Lines), &s.Branches.Cursor, &s.Branches.Offset, s.branchesPageSize())
		return
	}

	if s.Focus == FocusCommandLog {
		clampScrollView(len(s.CommandLog), &s.CommandLogView.Cursor, &s.CommandLogView.Offset, s.commandLogPageSize())
		return
	}

	clampScrollView(len(s.Changes.Rows), &s.Changes.Cursor, &s.Changes.Offset, s.changesPageSize())
}

func (s AppState) bodyHeight() int {
	h := s.Viewport.Height
	if h < 6 {
		return 6
	}
	return h
}

func (s AppState) GraphPaneHeight() int {
	content := max(8, s.bodyHeight()-s.CommandPaneHeight()-s.CommandLogPaneHeight())
	gh := (content * 45) / 100
	if gh < 4 {
		gh = 4
	}
	if gh > content-4 {
		gh = content - 4
	}
	return gh
}

func (s AppState) ChangesPaneHeight() int {
	content := max(8, s.bodyHeight()-s.CommandPaneHeight()-s.CommandLogPaneHeight())
	ch := content - s.GraphPaneHeight()
	if ch < 4 {
		return 4
	}
	return ch
}

func (s AppState) CommandPaneHeight() int {
	return 7
}

func (s AppState) CommandLogPaneHeight() int {
	return 5
}

func (s AppState) graphPageSize() int {
	h := s.GraphPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}

func (s AppState) changesPageSize() int {
	h := s.ChangesPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}

func (s AppState) commandLogPageSize() int {
	h := s.CommandLogPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}

func (s AppState) branchesPageSize() int {
	h := s.GraphPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}

func clampScrollView(total int, cursor, offset *int, page int) {
	if cursor == nil || offset == nil {
		return
	}
	if total <= 0 {
		*cursor = 0
		*offset = 0
		return
	}
	if *cursor < 0 {
		*cursor = 0
	}
	if *cursor >= total {
		*cursor = total - 1
	}
	if page < 1 {
		page = 1
	}
	if *cursor < *offset {
		*offset = *cursor
	}
	if *cursor >= *offset+page {
		*offset = *cursor - page + 1
	}
	maxOffset := max(0, total-page)
	if *offset > maxOffset {
		*offset = maxOffset
	}
	if *offset < 0 {
		*offset = 0
	}
}
