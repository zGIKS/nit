package state

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func (s AppState) topBarBoxes() (fetchX, fetchW, menuX, menuW int) {
	totalW := max(40, s.Viewport.Width)
	repoName := s.RepoName
	if repoName == "" {
		repoName = "unknown"
	}
	branchName := s.BranchName
	if branchName == "" {
		branchName = "-"
	}

	sep := strings.TrimSpace(s.RepoBranchSeparator)
	if sep == "" {
		sep = "->"
	}
	repoText := strings.TrimSpace(
		strings.TrimSpace(s.RepoLabel+" "+repoName) +
			" " + sep + " " +
			strings.TrimSpace(s.BranchLabel+" "+branchName),
	)
	createText := strings.TrimSpace(s.BranchesCreateButtonLabel())
	fetchText := strings.TrimSpace(s.FetchLabel)
	menuText := strings.TrimSpace(s.MenuLabel)

	repoW := max(16, runewidth.StringWidth(repoText)+4)
	createW := max(12, runewidth.StringWidth(createText)+4)
	fetchW = max(8, runewidth.StringWidth(fetchText)+4)
	menuW = max(8, runewidth.StringWidth(menuText)+4)
	minRepoW, minCreateW, minFetchW, minMenuW := 14, 12, 8, 8

	totalNeeded := repoW + createW + fetchW + menuW + 3
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
	shrink(&createW, minCreateW)
	shrink(&fetchW, minFetchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		repoW = max(minRepoW, repoW-overflow)
	}

	rightTopW := createW + fetchW + menuW + 2
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}

	createX := repoW + 1 + gapW + 1
	fetchX = createX + createW + 1
	menuX = fetchX + fetchW + 1
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

	sep := strings.TrimSpace(s.RepoBranchSeparator)
	if sep == "" {
		sep = "->"
	}
	repoText := strings.TrimSpace(
		strings.TrimSpace(s.RepoLabel+" "+repoName) +
			" " + sep + " " +
			strings.TrimSpace(s.BranchLabel+" "+branchName),
	)
	createText := strings.TrimSpace(s.BranchesCreateButtonLabel())
	fetchText := strings.TrimSpace(s.FetchLabel)
	menuText := strings.TrimSpace(s.MenuLabel)

	repoW = max(16, runewidth.StringWidth(repoText)+4)
	createW := max(12, runewidth.StringWidth(createText)+4)
	fetchW := max(8, runewidth.StringWidth(fetchText)+4)
	menuW = max(8, runewidth.StringWidth(menuText)+4)
	minRepoW, minCreateW, minFetchW, minMenuW := 14, 12, 8, 8

	totalNeeded := repoW + createW + fetchW + menuW + 3
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
	shrink(&createW, minCreateW)
	shrink(&fetchW, minFetchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		repoW = max(minRepoW, repoW-overflow)
	}

	rightTopW := createW + fetchW + menuW + 2
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}

	repoX = 0
	branchX = repoW + 1 + gapW + 1
	branchW = createW
	menuX = branchX + branchW + 1 + fetchW + 1
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

func (s AppState) menuMaxPanelHeight(y int) int {
	totalH := max(6, s.Viewport.Height)
	avail := totalH - y
	if avail < 3 {
		return 3
	}
	return avail
}

func (s AppState) MenuPanelRect() (x, y, w, h int) {
	totalW := max(40, s.Viewport.Width)
	_, _, menuX, menuW := s.topBarBoxes()
	maxItemW := dropdownItemsMaxWidth(dropdownMenuItems)
	w = max(18, max(menuW+10, maxItemW+2))
	h = min(len(dropdownMenuItems)+2, s.menuMaxPanelHeight(3))
	x = menuX + menuW - w
	if x < 0 {
		x = 0
	}
	if x+w > totalW {
		x = max(0, totalW-w)
	}
	return x, 3, w, h
}

func (s AppState) MenuSubmenuRect() (x, y, w, h int) {
	items := s.MenuSubmenuItems()
	if len(items) == 0 {
		return 0, 0, 0, 0
	}
	px, py, pw, _ := s.MenuPanelRect()
	anchorIdx := s.submenuAnchorIndex()
	if anchorIdx < 0 {
		return 0, 0, 0, 0
	}
	totalW := max(40, s.Viewport.Width)
	w = max(24, dropdownItemsMaxWidth(items)+2)
	y = py + 1 + (anchorIdx - s.MenuOffset)
	if y < py+1 {
		y = py + 1
	}
	h = min(len(items)+2, s.menuMaxPanelHeight(y))
	x = px + pw - 1
	if x+w > totalW {
		x = max(0, px-w+1)
	}
	if y < 0 {
		y = 0
	}
	return x, y, w, h
}
