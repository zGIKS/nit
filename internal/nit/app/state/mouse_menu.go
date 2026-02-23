package state

import "nit/internal/nit/app/actions"

func (s *AppState) handleMenuClick(x, y int) bool {
	if idx, ok := s.MenuItemIndexAt(x, y); ok {
		s.CloseMenu()
		_ = idx
		return true
	}
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		s.ToggleMenu()
		if s.MenuOpen {
			s.MenuHoverIndex = -1
		}
		return true
	}
	if s.MenuOpen {
		s.CloseMenu()
		return false
	}
	return false
}

func (s *AppState) MenuClickActionAt(x, y int) (actions.Action, bool) {
	idx, ok := s.MenuItemIndexAt(x, y)
	if !ok {
		return actions.ActionNone, false
	}
	item := s.MenuItems()[idx]
	s.CloseMenu()
	switch item {
	case "Pull":
		return actions.ActionPull, true
	case "Fetch":
		return actions.ActionFetch, true
	case "Push":
		return actions.ActionPush, true
	default:
		s.SetError("menu action not implemented yet")
		return actions.ActionNone, false
	}
}

func (s *AppState) ToggleMenuClick(x, y int) bool {
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		s.CloseBranchCreate()
		s.ToggleMenu()
		return true
	}
	return false
}

func (s *AppState) ToggleBranchCreateClick(x, y int) bool {
	if bx, by, bw, bh := s.BranchButtonRect(); x >= bx && x < bx+bw && y >= by && y < by+bh {
		s.CloseMenu()
		s.ToggleBranchCreate()
		return true
	}
	return false
}

func (s *AppState) CloseMenuOnOutsideClick(x, y int) {
	if !s.MenuOpen {
		return
	}
	if _, ok := s.MenuItemIndexAt(x, y); ok {
		return
	}
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		return
	}
	px, py, pw, ph := s.MenuPanelRect()
	if x >= px && x < px+pw && y >= py && y < py+ph {
		return
	}
	s.CloseMenu()
}

func (s *AppState) CloseTopMenusOnOutsideClick(x, y int) {
	s.CloseMenuOnOutsideClick(x, y)
	s.CloseBranchCreateOnOutsideClick(x, y)
}
