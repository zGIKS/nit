package app

func (s *AppState) Clamp() {
	if s.Focus == FocusGraph {
		if len(s.Graph.Lines) == 0 {
			s.Graph.Cursor = 0
			s.Graph.Offset = 0
			return
		}
		if s.Graph.Cursor < 0 {
			s.Graph.Cursor = 0
		}
		if s.Graph.Cursor >= len(s.Graph.Lines) {
			s.Graph.Cursor = len(s.Graph.Lines) - 1
		}
		page := s.graphPageSize()
		if s.Graph.Cursor < s.Graph.Offset {
			s.Graph.Offset = s.Graph.Cursor
		}
		if s.Graph.Cursor >= s.Graph.Offset+page {
			s.Graph.Offset = s.Graph.Cursor - page + 1
		}
		maxOffset := max(0, len(s.Graph.Lines)-page)
		if s.Graph.Offset > maxOffset {
			s.Graph.Offset = maxOffset
		}
		if s.Graph.Offset < 0 {
			s.Graph.Offset = 0
		}
		return
	}

	if len(s.Changes.Rows) == 0 {
		s.Changes.Cursor = 0
		s.Changes.Offset = 0
		return
	}
	if s.Changes.Cursor < 0 {
		s.Changes.Cursor = 0
	}
	if s.Changes.Cursor >= len(s.Changes.Rows) {
		s.Changes.Cursor = len(s.Changes.Rows) - 1
	}
	page := s.changesPageSize()
	if s.Changes.Cursor < s.Changes.Offset {
		s.Changes.Offset = s.Changes.Cursor
	}
	if s.Changes.Cursor >= s.Changes.Offset+page {
		s.Changes.Offset = s.Changes.Cursor - page + 1
	}
	maxOffset := max(0, len(s.Changes.Rows)-page)
	if s.Changes.Offset > maxOffset {
		s.Changes.Offset = maxOffset
	}
	if s.Changes.Offset < 0 {
		s.Changes.Offset = 0
	}
}

func (s AppState) bodyHeight() int {
	h := s.Viewport.Height
	if h < 6 {
		return 6
	}
	return h
}

func (s AppState) GraphPaneHeight() int {
	h := s.bodyHeight()
	gh := (h * 45) / 100
	if gh < 4 {
		gh = 4
	}
	if gh > h-4 {
		gh = h - 4
	}
	return gh
}

func (s AppState) ChangesPaneHeight() int {
	h := s.bodyHeight()
	ch := h - s.GraphPaneHeight()
	if ch < 4 {
		return 4
	}
	return ch
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
