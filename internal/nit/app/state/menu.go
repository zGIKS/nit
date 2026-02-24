package state

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type dropdownMenuItem struct {
	Label      string
	HasChevron bool
	Separator  bool
}

var dropdownMenuItems = []dropdownMenuItem{
	{Label: "Pull"},
	{Label: "Fetch"},
	{Separator: true},
	{Label: "Commit", HasChevron: true},
	{Label: "Changes", HasChevron: true},
	{Label: "Pull, Push", HasChevron: true},
	{Label: "Branch", HasChevron: true},
	{Label: "Remote", HasChevron: true},
	{Label: "Stash", HasChevron: true},
	{Label: "Tags", HasChevron: true},
	{Label: "Worktrees", HasChevron: true},
}

var commitDropdownMenuItems = []dropdownMenuItem{
	{Label: "Commit"},
	{Label: "Commit Staged"},
	{Label: "Commit All"},
	{Label: "Undo Last Commit"},
	{Label: "Abort Rebase"},
	{Separator: true},
	{Label: "Commit (Amend)"},
	{Label: "Commit Staged (Amend)"},
	{Label: "Commit All (Amend)"},
	{Separator: true},
	{Label: "Commit (Signed Off)"},
	{Label: "Commit Staged (Signed Off)"},
	{Label: "Commit All (Signed Off)"},
}

var changesDropdownMenuItems = []dropdownMenuItem{
	{Label: "Stage All Changes"},
	{Label: "Unstage All Changes"},
	{Label: "Discard All Changes"},
}

func (s *AppState) CloseMenu() {
	s.MenuOpen = false
	s.MenuHoverIndex = -1
	s.MenuOffset = 0
	s.MenuSubmenuKind = ""
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
}

func (s *AppState) ToggleMenu() {
	if s.MenuOpen {
		s.CloseMenu()
		return
	}
	s.MenuOpen = true
	s.MenuHoverIndex = s.firstSelectableMenuIndex()
	s.MenuOffset = 0
	s.MenuSubmenuKind = ""
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
}

func (s AppState) MenuItems() []dropdownMenuItem {
	return dropdownMenuItems
}

func dropdownItemsMaxWidth(items []dropdownMenuItem) int {
	maxItemW := 0
	for _, item := range items {
		if item.Separator {
			continue
		}
		itemW := runewidth.StringWidth(item.Label) + 2 // left/right padding
		if item.HasChevron {
			itemW += 2 // space + chevron
		}
		if itemW > maxItemW {
			maxItemW = itemW
		}
	}
	return maxItemW
}

func (s AppState) MenuSubmenuItems() []dropdownMenuItem {
	switch s.MenuSubmenuKind {
	case "commit":
		return commitDropdownMenuItems
	case "changes":
		return changesDropdownMenuItems
	default:
		return nil
	}
}

func firstSelectableIndex(items []dropdownMenuItem) int {
	for i, item := range items {
		if !item.Separator {
			return i
		}
	}
	return -1
}

func (s AppState) firstSelectableMenuIndex() int {
	return firstSelectableIndex(dropdownMenuItems)
}

func (s AppState) firstSelectableSubmenuIndex() int {
	return firstSelectableIndex(s.MenuSubmenuItems())
}

func menuPageSizeForRectHeight(h int) int {
	page := h - 2
	if page < 1 {
		return 1
	}
	return page
}

func (s AppState) menuMaxPanelHeight(y int) int {
	totalH := max(6, s.Viewport.Height)
	avail := totalH - y
	if avail < 3 {
		return 3
	}
	return avail
}

func clampScrollSelection(items []dropdownMenuItem, hover *int, offset *int, page int) {
	if len(items) == 0 {
		*hover = -1
		*offset = 0
		return
	}
	if page < 1 {
		page = 1
	}
	if *hover < 0 || *hover >= len(items) || items[*hover].Separator {
		*hover = firstSelectableIndex(items)
	}
	if *hover < 0 {
		*offset = 0
		return
	}
	if *hover < *offset {
		*offset = *hover
	}
	if *hover >= *offset+page {
		*offset = *hover - page + 1
	}
	maxOffset := max(0, len(items)-page)
	if *offset > maxOffset {
		*offset = maxOffset
	}
	if *offset < 0 {
		*offset = 0
	}
}

func nextSelectableIndex(items []dropdownMenuItem, start, delta int) int {
	if len(items) == 0 || delta == 0 {
		return start
	}
	if start < 0 || start >= len(items) {
		start = firstSelectableIndex(items)
	}
	if start < 0 {
		return -1
	}
	i := start
	for {
		i += delta
		if i < 0 {
			i = len(items) - 1
		}
		if i >= len(items) {
			i = 0
		}
		if !items[i].Separator {
			return i
		}
		if i == start {
			return start
		}
	}
}

func (s AppState) commitMenuIndex() int {
	return s.menuIndexByLabel("Commit")
}

func (s AppState) changesMenuIndex() int {
	return s.menuIndexByLabel("Changes")
}

func (s AppState) menuIndexByLabel(label string) int {
	for i, item := range dropdownMenuItems {
		if item.Label == label {
			return i
		}
	}
	return -1
}

func (s AppState) MenuHoverIsCommit() bool {
	return s.MenuHoverIndex == s.commitMenuIndex()
}

func (s AppState) MenuHoverHasSubmenu() bool {
	if s.MenuHoverIndex < 0 || s.MenuHoverIndex >= len(dropdownMenuItems) {
		return false
	}
	return dropdownMenuItems[s.MenuHoverIndex].HasChevron
}

func (s AppState) submenuAnchorIndex() int {
	switch s.MenuSubmenuKind {
	case "commit":
		return s.commitMenuIndex()
	case "changes":
		return s.changesMenuIndex()
	default:
		return -1
	}
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
	maxItemW := dropdownItemsMaxWidth(dropdownMenuItems)
	w = max(18, max(menuW+10, maxItemW+2))
	h = min(len(dropdownMenuItems)+2, s.menuMaxPanelHeight(3))
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
	idx := s.MenuOffset + (y - my - 1)
	if idx < 0 || idx >= len(dropdownMenuItems) {
		return -1, false
	}
	if dropdownMenuItems[idx].Separator {
		return -1, false
	}
	return idx, true
}

func (s AppState) MenuSubmenuItemIndexAt(x, y int) (int, bool) {
	items := s.MenuSubmenuItems()
	if !s.MenuOpen || len(items) == 0 {
		return -1, false
	}
	mx, my, mw, mh := s.MenuSubmenuRect()
	if mw <= 0 || mh <= 0 {
		return -1, false
	}
	if y < my || y >= my+mh || x < mx || x >= mx+mw {
		return -1, false
	}
	if y == my || y == my+mh-1 {
		return -1, false
	}
	idx := s.MenuSubOffset + (y - my - 1)
	if idx < 0 || idx >= len(items) {
		return -1, false
	}
	if items[idx].Separator {
		return -1, false
	}
	return idx, true
}

func (s *AppState) ensureMenuScrollVisible() {
	_, _, _, h := s.MenuPanelRect()
	clampScrollSelection(dropdownMenuItems, &s.MenuHoverIndex, &s.MenuOffset, menuPageSizeForRectHeight(h))
	if s.MenuSubmenuKind != "" {
		_, _, _, sh := s.MenuSubmenuRect()
		items := s.MenuSubmenuItems()
		clampScrollSelection(items, &s.MenuSubHoverIndex, &s.MenuSubOffset, menuPageSizeForRectHeight(sh))
	}
}

func (s *AppState) OpenCommitSubmenu() {
	s.MenuSubmenuKind = "commit"
	if s.MenuSubHoverIndex < 0 {
		s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
	}
	s.ensureMenuScrollVisible()
}

func (s *AppState) OpenChangesSubmenu() {
	s.MenuSubmenuKind = "changes"
	if s.MenuSubHoverIndex < 0 {
		s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
	}
	s.ensureMenuScrollVisible()
}

func (s *AppState) OpenSubmenuForMenuIndex(idx int) bool {
	if idx < 0 || idx >= len(dropdownMenuItems) {
		return false
	}
	s.MenuHoverIndex = idx
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
	switch dropdownMenuItems[idx].Label {
	case "Commit":
		s.OpenCommitSubmenu()
		return true
	case "Changes":
		s.OpenChangesSubmenu()
		return true
	default:
		s.CloseSubmenu()
		return false
	}
}

func (s *AppState) OpenHoveredSubmenu() bool {
	return s.OpenSubmenuForMenuIndex(s.MenuHoverIndex)
}

func (s *AppState) CloseSubmenu() {
	s.MenuSubmenuKind = ""
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
}

func (s *AppState) MoveMenuSelection(delta int) {
	if !s.MenuOpen || delta == 0 {
		return
	}
	s.MenuHoverIndex = nextSelectableIndex(dropdownMenuItems, s.MenuHoverIndex, delta)
	if s.MenuHoverIndex >= 0 && s.MenuHoverIndex < len(dropdownMenuItems) {
		s.OpenSubmenuForMenuIndex(s.MenuHoverIndex)
	}
	s.ensureMenuScrollVisible()
}

func (s *AppState) MoveMenuSubmenuSelection(delta int) {
	items := s.MenuSubmenuItems()
	if !s.MenuOpen || len(items) == 0 || delta == 0 {
		return
	}
	s.MenuSubHoverIndex = nextSelectableIndex(items, s.MenuSubHoverIndex, delta)
	s.ensureMenuScrollVisible()
}

func (s *AppState) MenuWheelAt(x, y, delta int) bool {
	if !s.MenuOpen || delta == 0 {
		return false
	}
	if _, ok := s.MenuSubmenuItemIndexAt(x, y); ok {
		if s.MenuSubHoverIndex < 0 {
			s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
		}
		s.MoveMenuSubmenuSelection(delta)
		return true
	}
	if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
		if s.MenuSubHoverIndex < 0 {
			s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
		}
		s.MoveMenuSubmenuSelection(delta)
		return true
	}
	if _, ok := s.MenuItemIndexAt(x, y); ok {
		if s.MenuHoverIndex < 0 {
			s.MenuHoverIndex = s.firstSelectableMenuIndex()
		}
		s.MoveMenuSelection(delta)
		return true
	}
	if mx, my, mw, mh := s.MenuPanelRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		if s.MenuHoverIndex < 0 {
			s.MenuHoverIndex = s.firstSelectableMenuIndex()
		}
		s.MoveMenuSelection(delta)
		return true
	}
	return false
}
