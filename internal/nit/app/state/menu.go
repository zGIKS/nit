package state

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

var dropdownMenuItems = []string{
	"Pull",
	"Fetch",
}

func (s *AppState) CloseMenu() {
	s.MenuOpen = false
	s.MenuHoverIndex = -1
}

func (s *AppState) ToggleMenu() {
	s.MenuOpen = !s.MenuOpen
	if !s.MenuOpen {
		s.MenuHoverIndex = -1
	}
}

func (s AppState) MenuItems() []string {
	return dropdownMenuItems
}

func (s AppState) topBarBoxes() (fetchX, fetchW, menuX, menuW int) {
	_, branchW, _, _, menuX, menuW := s.topBarBoxRects()
	_ = branchW
	fetchX, fetchW = -1, 0
	return fetchX, fetchW, menuX, menuW
}

func (s AppState) topBarBoxRects() (repoX, repoW, branchX, branchW, menuX, menuW int) {
	totalW := max(40, s.Viewport.Width)
	repoName := s.RepoName
	if repoName == "" {
		repoName = "unknown"
	}
	branchName := s.BranchName
	if branchName == "" {
		branchName = "-"
	}

	repoText := strings.TrimSpace(s.RepoLabel + " " + repoName)
	branchText := strings.TrimSpace(s.BranchLabel + " " + branchName)
	menuText := strings.TrimSpace(s.MenuLabel)

	repoW = max(16, runewidth.StringWidth(repoText)+4)
	branchW = max(16, runewidth.StringWidth(branchText)+4)
	menuW = max(8, runewidth.StringWidth(menuText)+4)
	minRepoW, minBranchW, minMenuW := 14, 12, 8

	totalNeeded := repoW + branchW + menuW + 2
	overflow := totalNeeded - totalW
	shrink := func(w *int, minW int) {
		if overflow <= 0 {
			return
		}
		can := *w - minW
		if can <= 0 {
			return
		}
		d := min(can, overflow)
		*w -= d
		overflow -= d
	}
	shrink(&repoW, minRepoW)
	shrink(&branchW, minBranchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		repoW = max(minRepoW, repoW-overflow)
	}

	rightTopW := branchW + menuW + 1
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}

	repoX = 0
	branchX = repoW + 1 + gapW + 1
	menuX = branchX + branchW + 1
	return repoX, repoW, branchX, branchW, menuX, menuW
}

func (s AppState) MenuButtonRect() (x, y, w, h int) {
	_, _, menuX, menuW := s.topBarBoxes()
	return menuX, 0, menuW, 3
}

func (s AppState) BranchButtonRect() (x, y, w, h int) {
	_, _, x, w, _, _ = s.topBarBoxRects()
	return x, 0, w, 3
}

func (s AppState) FetchButtonRect() (x, y, w, h int) {
	fetchX, fetchW, _, _ := s.topBarBoxes()
	return fetchX, 0, fetchW, 3
}

func (s AppState) MenuPanelRect() (x, y, w, h int) {
	totalW := max(40, s.Viewport.Width)
	_, _, menuX, menuW := s.topBarBoxes()
	w = max(18, menuW+10)
	h = len(dropdownMenuItems) + 2
	// Right-align with the menu button and place under the top bar with one spacer row.
	x = menuX + menuW - w
	if x < 0 {
		x = 0
	}
	if x+w > totalW {
		x = max(0, totalW-w)
	}
	return x, 3, w, h
}

func (s AppState) MenuItemIndexAt(x, y int) (int, bool) {
	if !s.MenuOpen {
		return -1, false
	}
	mx, my, mw, mh := s.MenuPanelRect()
	if y < my || y >= my+mh || x < mx || x >= mx+mw {
		return -1, false
	}
	if y == my || y == my+mh-1 {
		return -1, false
	}
	idx := y - my - 1
	if idx < 0 || idx >= len(dropdownMenuItems) {
		return -1, false
	}
	return idx, true
}
