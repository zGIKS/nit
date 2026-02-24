package state

import "github.com/zGIKS/nit/internal/nit/app/actions"

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
	return s.MenuActivateIndex(idx)
}

func (s *AppState) MenuActivateIndex(idx int) (actions.Action, bool) {
	items := s.MenuItems()
	if idx < 0 || idx >= len(items) || items[idx].Separator {
		return actions.ActionNone, false
	}
	item := items[idx]
	if item.HasChevron {
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	}
	switch item.Label {
	case "Pull":
		s.CloseMenu()
		return actions.ActionPull, true
	case "Fetch":
		s.CloseMenu()
		return actions.ActionFetch, true
	default:
		s.MenuSubmenuKind = ""
		s.MenuSubHoverIndex = -1
		s.MenuHoverIndex = idx
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
	if bx, by, bw, bh := s.BranchButtonRect(); bw > 0 && x >= bx && x < bx+bw && y >= by && y < by+bh {
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
	if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
		return
	}
	s.CloseMenu()
}

func (s *AppState) CloseTopMenusOnOutsideClick(x, y int) {
	s.CloseMenuOnOutsideClick(x, y)
	s.CloseBranchCreateOnOutsideClick(x, y)
}
