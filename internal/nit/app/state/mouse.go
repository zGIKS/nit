package state

import (
	"nit/internal/nit/app/actions"
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
	if y < 0 || x < 0 {
		s.MenuHoverIndex = -1
		return
	}

	if fx, fy, fw, fh := s.FetchButtonRect(); x >= fx && x < fx+fw && y >= fy && y < fy+fh {
		s.HoverFetch = true
	}
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		s.HoverMenu = true
	}
	if idx, ok := s.MenuItemIndexAt(x, y); ok {
		s.MenuHoverIndex = idx
		return
	}
	s.MenuHoverIndex = -1
}

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
		s.ToggleMenu()
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

func (s *AppState) clickCommandBox(y, top int) bool {
	h := s.CommandPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusCommand)
	return true
}

func (s *AppState) wheelCommandBox(y, top, delta int) bool {
	h := s.CommandPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusCommand)
	return true
}

func (s *AppState) clickChangesBox(y, top int) bool {
	h := s.ChangesPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusChanges)
	if idx, ok := boxContentLine(y, top, h); ok {
		row := s.Changes.Offset + idx
		if row >= 0 && row < len(s.Changes.Rows) && s.Changes.Rows[row].Selectable {
			s.Changes.Cursor = row
		}
	}
	return true
}

func (s *AppState) wheelChangesBox(y, top, delta int) bool {
	h := s.ChangesPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusChanges)
	s.Changes.Cursor += delta
	if delta >= 0 {
		s.snapChangesCursor(1)
	} else {
		s.snapChangesCursor(-1)
	}
	return true
}

func (s *AppState) clickGraphBox(x, y, top int) bool {
	h := s.GraphPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	graphW, _ := s.GraphBranchesPaneWidths()
	if x > graphW {
		s.focusByMouse(FocusBranches)
		if idx, ok := boxContentLine(y, top, h); ok {
			line := s.Branches.Offset + idx
			if line >= 0 && line < len(s.Branches.Lines) {
				s.Branches.Cursor = line
			}
		}
		return true
	}
	s.focusByMouse(FocusGraph)
	if idx, ok := boxContentLine(y, top, h); ok {
		line := s.Graph.Offset + idx
		if line >= 0 && line < len(s.Graph.Lines) {
			s.Graph.Cursor = line
		}
	}
	return true
}

func (s *AppState) wheelGraphBox(x, y, top, delta int) bool {
	h := s.GraphPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	graphW, _ := s.GraphBranchesPaneWidths()
	if x > graphW {
		s.focusByMouse(FocusBranches)
		s.Branches.Cursor += delta
		return true
	}
	s.focusByMouse(FocusGraph)
	s.Graph.Cursor += delta
	return true
}

func (s *AppState) clickCommandLogBox(y, top int) bool {
	h := s.CommandLogPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusCommandLog)
	if idx, ok := boxContentLine(y, top, h); ok {
		line := s.CommandLogView.Offset + idx
		if line >= 0 && line < len(s.CommandLog) {
			s.CommandLogView.Cursor = line
		}
	}
	return true
}

func (s *AppState) wheelCommandLogBox(y, top, delta int) bool {
	h := s.CommandLogPaneHeight()
	if y < top || y >= top+h {
		return false
	}
	s.focusByMouse(FocusCommandLog)
	s.CommandLogView.Cursor += delta
	return true
}

func (s *AppState) focusByMouse(target FocusState) {
	if target == FocusCommand {
		if s.Focus != FocusCommand {
			s.Command.ReturnFocus = s.Focus
		}
		s.Focus = FocusCommand
		s.Command.SelectAll = false
		s.MoveCommandCursorToEnd()
		return
	}

	s.Focus = target
	s.Command.SelectAll = false
	if target == FocusChanges {
		s.snapChangesCursor(1)
	}
}

func boxContentLine(y, top, boxHeight int) (int, bool) {
	contentTop := top + 1
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}
	if y < contentTop || y >= contentTop+contentHeight {
		return 0, false
	}
	return y - contentTop, true
}
