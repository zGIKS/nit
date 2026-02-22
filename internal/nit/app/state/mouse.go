package state

func (s *AppState) HandleMouseClick(x, y int) {
	_ = x // boxes span full width for now
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

	if s.clickGraphBox(y, top) {
		s.Clamp()
		return
	}
	top += s.GraphPaneHeight()

	if s.clickCommandLogBox(y, top) {
		s.Clamp()
		return
	}
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

	if s.wheelGraphBox(y, top, delta) {
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
	if idx, ok := boxContentLine(y, top, h); ok && idx < len(s.Changes.Rows) && s.Changes.Rows[idx].Selectable {
		s.Changes.Cursor = idx
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

func (s *AppState) clickGraphBox(y, top int) bool {
	h := s.GraphPaneHeight()
	if y < top || y >= top+h {
		return false
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

func (s *AppState) wheelGraphBox(y, top, delta int) bool {
	h := s.GraphPaneHeight()
	if y < top || y >= top+h {
		return false
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
