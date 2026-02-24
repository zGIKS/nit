package state

func (s *AppState) BranchCreateWheelAt(x, y, delta int) bool {
	if !s.BranchCreateOpen || delta == 0 {
		return false
	}
	px, py, pw, ph := s.BranchCreatePanelRect()
	if x < px || x >= px+pw || y < py || y >= py+ph {
		return false
	}
	lx, ly, lw, lh := s.BranchCreateSourceListRect()
	if x >= lx && x < lx+lw && y >= ly && y < ly+lh {
		s.BranchCreateMoveSource(delta)
		s.BranchCreateHoverAt(x, y)
		return true
	}
	if len(s.BranchCreateSourceList) > 0 {
		s.BranchCreateMoveSource(delta)
		s.BranchCreateHoverAt(x, y)
		return true
	}
	return true
}

func (s *AppState) BranchCreateClick(x, y int) bool {
	if !s.BranchCreateOpen {
		return false
	}
	if idx, ok := s.BranchCreateSourceIndexAt(x, y); ok {
		s.BranchCreateSelectSourceIndex(idx)
		return true
	}
	px, py, pw, ph := s.BranchCreatePanelRect()
	if x >= px && x < px+pw && y >= py && y < py+ph {
		nx, ny, nw, nh := s.BranchCreateNameInputRect()
		if x >= nx && x < nx+nw && y >= ny && y < ny+nh {
			moveTextInputCursorEnd(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
		}
		return true
	}
	s.CloseBranchCreate()
	return false
}

func (s *AppState) CloseBranchCreateOnOutsideClick(x, y int) {
	if !s.BranchCreateOpen {
		return
	}
	px, py, pw, ph := s.BranchCreatePanelRect()
	if x >= px && x < px+pw && y >= py && y < py+ph {
		return
	}
	s.CloseBranchCreate()
}

func (s *AppState) clampBranchCreateSourceOffset() {
	_, _, _, h := s.BranchCreateSourceListRect()
	maxOffset := len(s.BranchCreateSourceList) - h
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.BranchCreateSourceOffset < 0 {
		s.BranchCreateSourceOffset = 0
	}
	if s.BranchCreateSourceOffset > maxOffset {
		s.BranchCreateSourceOffset = maxOffset
	}
}

func (s *AppState) ensureBranchCreateSourceVisible() {
	idx := findStringIndex(s.BranchCreateSourceList, s.BranchCreateSource)
	if idx < 0 {
		s.clampBranchCreateSourceOffset()
		return
	}
	_, _, _, h := s.BranchCreateSourceListRect()
	if h < 1 {
		h = 1
	}
	if idx < s.BranchCreateSourceOffset {
		s.BranchCreateSourceOffset = idx
	}
	if idx >= s.BranchCreateSourceOffset+h {
		s.BranchCreateSourceOffset = idx - h + 1
	}
	s.clampBranchCreateSourceOffset()
}

func findStringIndex(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
