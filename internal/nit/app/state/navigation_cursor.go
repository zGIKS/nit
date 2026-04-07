package state

func (s *AppState) canPush() bool {
	if len(s.Graph.Lines) == 0 {
		return true
	}
	if len(s.Graph.Lines) == 1 {
		line := s.Graph.Lines[0]
		if line == "No commits to display." || line == "Not a git repo or no commits yet." {
			return false
		}
	}
	return true
}

func (s *AppState) moveCursor(delta int) {
	if s.Focus == FocusGraph {
		s.Graph.Cursor += delta
		s.Clamp()
		return
	}
	if s.Focus == FocusBranches {
		s.Branches.Cursor += delta
		s.Clamp()
		return
	}
	if s.Focus == FocusCommandLog {
		s.CommandLogView.Cursor += delta
		s.Clamp()
		return
	}
	if s.Focus != FocusChanges {
		return
	}
	s.Changes.Cursor += delta
	if delta >= 0 {
		s.snapChangesCursor(1)
	} else {
		s.snapChangesCursor(-1)
	}
	s.Clamp()
}
