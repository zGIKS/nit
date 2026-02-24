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
	if y < 0 || y >= 3 || x < 0 {
		return actions.ActionNone, false
	}
	if fx, fy, fw, fh := s.FetchButtonRect(); fw > 0 && x >= fx && x < fx+fw && y >= fy && y < fy+fh {
		return actions.ActionFetch, true
	}
	return actions.ActionNone, false
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
