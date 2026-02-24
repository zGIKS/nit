package state

import (
	"github.com/zGIKS/nit/internal/nit/app/actions"
)

func (s *AppState) HandleMouseClick(x, y int) {
	if y < 0 {
		return
	}

	top := 0

	if s.clickCommandBox(y, top) {
		s.Clamp()
		return
	}
	top += s.CommandPaneHeight()

	if s.clickChangesBox(y, top) {
		s.Clamp()
		return
	}
	top += s.ChangesPaneHeight()

	if s.clickGraphBox(x, y, top) {
		s.Clamp()
		return
	}
	top += s.GraphPaneHeight()

	if s.clickCommandLogBox(y, top) {
		s.Clamp()
		return
	}
}

func (s *AppState) TopBarActionAt(x, y int) (actions.Action, bool) {
	// Top bar is the first 3 rows of the command pane.
	if y < 0 || y >= 3 || x < 0 {
		return actions.ActionNone, false
	}
	return actions.ActionNone, false
}

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
			switch item.Label {
			case "Commit":
				s.MenuSubmenuKind = "commit"
			case "Changes":
				s.MenuSubmenuKind = "changes"
			case "Pull, Push":
				s.MenuSubmenuKind = "pull_push"
			case "Branch":
				s.MenuSubmenuKind = "branch"
			case "Remote":
				s.MenuSubmenuKind = "remote"
			case "Stash":
				s.MenuSubmenuKind = "stash"
			case "Tags":
				s.MenuSubmenuKind = "tags"
			}
		}
	}
	if s.MenuOpen && s.MenuSubmenuKind == "" && prevSubmenuKind != "" {
		s.MenuSubmenuKind = prevSubmenuKind
		if idx, ok := s.MenuSubmenuItemIndexAt(x, y); ok {
			s.MenuSubHoverIndex = idx
		} else if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
			// Keep submenu open while hovering over separators/empty rows.
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

func (s *AppState) HandleMouseWheel(x, y, delta int) {
	if y < 0 || delta == 0 {
		return
	}
	if s.MenuWheelAt(x, y, delta) {
		return
	}

	top := 0

	if s.wheelCommandBox(y, top, delta) {
		s.Clamp()
		return
	}
	top += s.CommandPaneHeight()

	if s.wheelChangesBox(y, top, delta) {
		s.Clamp()
		return
	}
	top += s.ChangesPaneHeight()

	if s.wheelGraphBox(x, y, top, delta) {
		s.Clamp()
		return
	}
	top += s.GraphPaneHeight()

	if s.wheelCommandLogBox(y, top, delta) {
		s.Clamp()
		return
	}
}
