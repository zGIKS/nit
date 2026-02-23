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
	if idx, ok := s.MenuItemIndexAt(x, y); ok {
		s.MenuHoverIndex = idx
		s.BranchCreateHoverAt(x, y)
		return
	}
	s.MenuHoverIndex = -1
	s.BranchCreateHoverAt(x, y)
}

func (s *AppState) HandleMouseWheel(x, y, delta int) {
	_ = x // boxes span full width for now
	if y < 0 || delta == 0 {
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
