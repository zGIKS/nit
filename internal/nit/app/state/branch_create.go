package state

import "strings"

func (s *AppState) OpenBranchCreate() {
	s.BranchCreateOpen = true
	s.MenuOpen = false
	s.BranchCreateHoverIndex = -1
	s.BranchCreateSelectAll = false
	s.syncBranchCreateSources()
	moveTextInputCursorEnd(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) CloseBranchCreate() {
	s.BranchCreateOpen = false
	s.BranchCreateHoverIndex = -1
}

func (s *AppState) ToggleBranchCreate() {
	if s.BranchCreateOpen {
		s.CloseBranchCreate()
		return
	}
	s.OpenBranchCreate()
}

func (s *AppState) syncBranchCreateSources() {
	sources := make([]string, 0, len(s.Branches.Lines))
	for _, line := range s.Branches.Lines {
		name := strings.TrimSpace(line)
		if name == "" || strings.Contains(name, "No local branches") || strings.Contains(name, "Loading branches") || strings.Contains(name, "Not a git repo") {
			continue
		}
		name = strings.TrimPrefix(name, "●")
		name = strings.TrimPrefix(name, "*")
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		sources = append(sources, name)
	}
	s.BranchCreateSourceList = sources
	if s.BranchCreateSource == "" && s.BranchName != "" {
		s.BranchCreateSource = s.BranchName
	}
	if findStringIndex(s.BranchCreateSourceList, s.BranchCreateSource) < 0 {
		if len(s.BranchCreateSourceList) > 0 {
			s.BranchCreateSource = s.BranchCreateSourceList[0]
		} else if s.BranchCreateSource == "" {
			s.BranchCreateSource = "-"
		}
	}
	s.ensureBranchCreateSourceVisible()
}

func (s *AppState) BranchCreateAppendText(text string) {
	appendTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll, text)
}

func (s *AppState) BranchCreateBackspace() {
	backspaceTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateDelete() {
	deleteTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorLeft() {
	moveTextInputCursorLeft(&s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorRight() {
	moveTextInputCursorRight(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorHome() {
	moveTextInputCursorHome(&s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorEnd() {
	moveTextInputCursorEnd(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateSelectAllText() {
	selectAllTextInput(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s AppState) SelectedBranchCreateText() string {
	if s.BranchCreateSelectAll {
		return s.BranchCreateName
	}
	return ""
}

func (s *AppState) DeleteBranchCreateSelection() {
	clearSelectedText(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s AppState) BranchCreatePanelRect() (x, y, w, h int) {
	totalW := max(40, s.Viewport.Width)
	totalH := max(12, s.Viewport.Height)
	_, _, bx, bw, _, _ := s.topBarBoxRects()
	w = 56
	if w > totalW {
		w = totalW
	}
	if w < 36 {
		w = 36
	}
	listCount := len(s.BranchCreateSourceList)
	if listCount > 6 {
		listCount = 6
	}
	if listCount < 1 {
		listCount = 1
	}
	h = 9 + listCount
	x = bx + bw - w
	if x < 0 {
		x = 0
	}
	if x+w > totalW {
		x = totalW - w
		if x < 0 {
			x = 0
		}
	}
	y = 2
	if y+h > totalH {
		y = max(0, totalH-h)
	}
	return x, y, w, h
}

func (s AppState) BranchCreateNameInputRect() (x, y, w, h int) {
	px, py, pw, _ := s.BranchCreatePanelRect()
	return px + 2, py + 3, max(10, pw-4), 1
}

func (s AppState) BranchCreateSourceListRect() (x, y, w, h int) {
	px, py, pw, ph := s.BranchCreatePanelRect()
	x = px + 2
	y = py + 6
	w = max(10, pw-4)
	h = ph - 7
	if h < 1 {
		h = 1
	}
	return x, y, w, h
}

func (s AppState) BranchCreateSourceIndexAt(x, y int) (int, bool) {
	if !s.BranchCreateOpen {
		return -1, false
	}
	lx, ly, lw, lh := s.BranchCreateSourceListRect()
	if x < lx || x >= lx+lw || y < ly || y >= ly+lh {
		return -1, false
	}
	idx := s.BranchCreateSourceOffset + (y - ly)
	if idx < 0 || idx >= len(s.BranchCreateSourceList) {
		return -1, false
	}
	return idx, true
}

func (s *AppState) BranchCreateSelectSourceIndex(idx int) {
	if idx < 0 || idx >= len(s.BranchCreateSourceList) {
		return
	}
	s.BranchCreateSource = s.BranchCreateSourceList[idx]
	s.ensureBranchCreateSourceVisible()
}

func (s *AppState) BranchCreateMoveSource(delta int) {
	if len(s.BranchCreateSourceList) == 0 || delta == 0 {
		return
	}
	cur := 0
	for i, name := range s.BranchCreateSourceList {
		if name == s.BranchCreateSource {
			cur = i
			break
		}
	}
	cur += delta
	if cur < 0 {
		cur = 0
	}
	if cur >= len(s.BranchCreateSourceList) {
		cur = len(s.BranchCreateSourceList) - 1
	}
	s.BranchCreateSource = s.BranchCreateSourceList[cur]
	s.ensureBranchCreateSourceVisible()
}

func (s *AppState) BranchCreateHoverAt(x, y int) {
	if idx, ok := s.BranchCreateSourceIndexAt(x, y); ok {
		s.BranchCreateHoverIndex = idx
		return
	}
	s.BranchCreateHoverIndex = -1
}

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
	if bx, by, bw, bh := s.BranchButtonRect(); x >= bx && x < bx+bw && y >= by && y < by+bh {
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
