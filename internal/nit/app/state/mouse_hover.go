package state

func (s *AppState) HandleMouseMove(x, y int) {
	s.HoverFetch = false
	s.HoverMenu = false
	s.HoverBranch = false
	if y < 0 || x < 0 {
		s.MenuHoverIndex = -1
		s.MenuSubmenuKind = ""
		s.MenuSubHoverIndex = -1
		s.BranchCreateHoverIndex = -1
		return
	}

	if fx, fy, fw, fh := s.FetchButtonRect(); x >= fx && x < fx+fw && y >= fy && y < fy+fh {
		s.HoverFetch = true
	}
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		s.HoverMenu = true
	}
	if bx, by, bw, bh := s.BranchButtonRect(); x >= bx && x < bx+bw && y >= by && y < by+bh {
		s.HoverBranch = true
	}
	menuIdx, menuOK := s.MenuItemIndexAt(x, y)
	if menuOK {
		s.MenuHoverIndex = menuIdx
	} else {
		s.MenuHoverIndex = -1
	}
	prevSubmenuKind := s.MenuSubmenuKind
	s.MenuSubmenuKind = ""
	s.MenuSubHoverIndex = -1
	if s.MenuOpen && menuOK {
		if item := s.MenuItems()[menuIdx]; item.HasChevron {
			if kind, ok := submenuKindByLabel[item.Label]; ok {
				s.MenuSubmenuKind = kind
			}
		}
	}
	if s.MenuOpen && s.MenuSubmenuKind == "" && prevSubmenuKind != "" {
		s.MenuSubmenuKind = prevSubmenuKind
		if idx, ok := s.MenuSubmenuItemIndexAt(x, y); ok {
			s.MenuSubHoverIndex = idx
		} else if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
			if anchorIdx := s.submenuAnchorIndex(); anchorIdx >= 0 {
				s.MenuHoverIndex = anchorIdx
			}
		} else {
			s.MenuSubmenuKind = ""
		}
	} else if s.MenuOpen && s.MenuSubmenuKind != "" {
		if idx, ok := s.MenuSubmenuItemIndexAt(x, y); ok {
			s.MenuSubHoverIndex = idx
		} else if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
			if anchorIdx := s.submenuAnchorIndex(); anchorIdx >= 0 {
				s.MenuHoverIndex = anchorIdx
			}
		}
	}
	s.BranchCreateHoverAt(x, y)
}
