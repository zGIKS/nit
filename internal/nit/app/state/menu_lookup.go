package state

var menuIndexByLabelMap = func() map[string]int {
	m := make(map[string]int, len(dropdownMenuItems))
	for i, item := range dropdownMenuItems {
		if item.Label != "" {
			m[item.Label] = i
		}
	}
	return m
}()

func (s AppState) MenuItems() []DropdownMenuItem {
	return dropdownMenuItems
}

func (s AppState) MenuSubmenuItems() []DropdownMenuItem {
	return submenuItemsByKind[s.MenuSubmenuKind]
}

func (s AppState) firstSelectableMenuIndex() int {
	return firstSelectableIndex(dropdownMenuItems)
}

func (s AppState) firstSelectableSubmenuIndex() int {
	return firstSelectableIndex(s.MenuSubmenuItems())
}

func (s AppState) menuIndexByLabel(label string) int {
	if idx, ok := menuIndexByLabelMap[label]; ok {
		return idx
	}
	return -1
}

func (s AppState) commitMenuIndex() int {
	return s.menuIndexByLabel("Commit")
}

func (s AppState) changesMenuIndex() int {
	return s.menuIndexByLabel("Changes")
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
	for i, item := range dropdownMenuItems {
		if kind, ok := submenuKindByLabel[item.Label]; ok && kind == s.MenuSubmenuKind {
			return i
		}
	}
	return -1
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
