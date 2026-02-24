package state

func (s AppState) BranchCreatePanelRect() (x, y, w, h int) {
	totalW := max(40, s.Viewport.Width)
	totalH := max(12, s.Viewport.Height)
	w = 56
	if w > totalW {
		w = totalW
	}
	if w < 36 {
		w = 36
	}
	baseRows := 11 // border/title + helpers + separators + name/input + source + bottom
	maxPanelH := totalH - 2
	if maxPanelH < baseRows+1 {
		maxPanelH = baseRows + 1
	}
	listRows := len(s.BranchCreateSourceList)
	if listRows < 1 {
		listRows = 1
	}
	maxListRows := maxPanelH - baseRows
	if maxListRows < 1 {
		maxListRows = 1
	}
	if listRows > maxListRows {
		listRows = maxListRows
	}
	h = baseRows + listRows
	x = (totalW - w) / 2
	if x < 0 {
		x = 0
	}
	y = (totalH - h) / 2
	if y < 0 {
		y = 0
	}
	if y+h > totalH {
		y = max(0, totalH-h)
	}
	return x, y, w, h
}

func (s AppState) BranchesCreateButtonLabel() string {
	return "+ branch"
}

func (s AppState) BranchesCreateButtonRect() (x, y, w, h int) {
	return s.BranchButtonRect()
}

func (s AppState) BranchCreateNameInputRect() (x, y, w, h int) {
	px, py, pw, _ := s.BranchCreatePanelRect()
	return px + 2, py + 3, max(10, pw-4), 1
}

func (s AppState) BranchCreateSourceListRect() (x, y, w, h int) {
	px, py, pw, ph := s.BranchCreatePanelRect()
	x = px + 2
	y = py + 10
	w = max(10, pw-4)
	h = ph - 11
	if h < 1 {
		h = 1
	}
	return x, y, w, h
}
