package app

func (s *AppState) Apply(action Action) ApplyResult {
	res := ApplyResult{}
	switch action {
	case ActionQuit:
		res.Quit = true
	case ActionTogglePanel:
		if s.Focus == FocusChanges {
			s.Focus = FocusGraph
		} else {
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
		}
	case ActionMoveDown:
		s.moveCursor(1)
	case ActionMoveUp:
		s.moveCursor(-1)
	case ActionToggleOne:
		if s.Focus != FocusChanges {
			break
		}
		entry, section, ok := s.selectedChange()
		if !ok {
			break
		}
		if section == SectionStaged {
			s.Changes.StickySection = SectionStaged
			res.Operations = []Operation{{Kind: OpUnstagePath, Path: entry.Path}}
		} else {
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []Operation{{Kind: OpStagePath, Path: entry.Path}}
		}
		res.RefreshChanges = true
	case ActionStageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionStaged
			res.Operations = []Operation{{Kind: OpStageAll}}
			res.RefreshChanges = true
		}
	case ActionUnstageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []Operation{{Kind: OpUnstageAll}}
			res.RefreshChanges = true
		}
	}
	s.Clamp()
	return res
}
