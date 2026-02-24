package state

func (s *AppState) CloseMenu() {
	s.MenuOpen = false
	s.MenuHoverIndex = -1
	s.MenuOffset = 0
	s.MenuSubActive = false
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
	s.MenuSubActive = false
	s.MenuSubmenuKind = ""
	s.OpenSubmenuForMenuIndex(s.MenuHoverIndex)
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
}

func (s *AppState) openSubmenu(kind string) {
	s.MenuSubmenuKind = kind
	if s.MenuSubHoverIndex < 0 {
		s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
	}
	s.ensureMenuScrollVisible()
}

func (s *AppState) OpenCommitSubmenu()   { s.openSubmenu("commit") }
func (s *AppState) OpenChangesSubmenu()  { s.openSubmenu("changes") }
func (s *AppState) OpenPullPushSubmenu() { s.openSubmenu("pull_push") }
func (s *AppState) OpenBranchSubmenu()   { s.openSubmenu("branch") }
func (s *AppState) OpenRemoteSubmenu()   { s.openSubmenu("remote") }
func (s *AppState) OpenStashSubmenu()    { s.openSubmenu("stash") }
func (s *AppState) OpenTagsSubmenu()     { s.openSubmenu("tags") }

func (s *AppState) OpenSubmenuForMenuIndex(idx int) bool {
	if idx < 0 || idx >= len(dropdownMenuItems) {
		return false
	}
	s.MenuHoverIndex = idx
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
	kind, ok := submenuKindByLabel[dropdownMenuItems[idx].Label]
	if !ok {
		s.CloseSubmenu()
		return false
	}
	s.openSubmenu(kind)
	return true
}

func (s *AppState) OpenHoveredSubmenu() bool {
	if s.OpenSubmenuForMenuIndex(s.MenuHoverIndex) {
		s.MenuSubActive = true
		if s.MenuSubHoverIndex < 0 {
			s.MenuSubHoverIndex = s.firstSelectableSubmenuIndex()
		}
		return true
	}
	return false
}

func (s *AppState) CloseSubmenu() {
	s.MenuSubmenuKind = ""
	s.MenuSubHoverIndex = -1
	s.MenuSubOffset = 0
	s.MenuSubActive = false
}

func (s *AppState) MoveMenuSelection(delta int) {
	if !s.MenuOpen || delta == 0 {
		return
	}
	s.MenuHoverIndex = nextSelectableIndex(dropdownMenuItems, s.MenuHoverIndex, delta)
	if s.MenuHoverIndex >= 0 && s.MenuHoverIndex < len(dropdownMenuItems) {
		s.OpenSubmenuForMenuIndex(s.MenuHoverIndex)
		s.MenuSubHoverIndex = -1
		s.MenuSubOffset = 0
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

func (s *AppState) ensureMenuScrollVisible() {
	_, _, _, h := s.MenuPanelRect()
	clampScrollSelection(dropdownMenuItems, &s.MenuHoverIndex, &s.MenuOffset, menuPageSizeForRectHeight(h))
	if s.MenuSubmenuKind != "" {
		_, _, _, sh := s.MenuSubmenuRect()
		items := s.MenuSubmenuItems()
		clampScrollSelection(items, &s.MenuSubHoverIndex, &s.MenuSubOffset, menuPageSizeForRectHeight(sh))
	}
}
