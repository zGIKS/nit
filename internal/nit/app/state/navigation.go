package state

import (
	"nit/internal/nit/app/actions"
	"nit/internal/nit/git"
)

func (s *AppState) Apply(action actions.Action) actions.ApplyResult {
	res := actions.ApplyResult{}
	switch action {
	case actions.ActionQuit:
		res.Quit = true
	case actions.ActionFocusCommand:
		s.Focus = FocusCommand
		s.MoveCommandCursorToEnd()
	case actions.ActionTogglePanel:
		switch s.Focus {
		case FocusCommand:
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
		case FocusChanges:
			s.Focus = FocusGraph
		default:
			s.Focus = FocusCommand
		}
	case actions.ActionMoveDown:
		s.moveCursor(1)
	case actions.ActionMoveUp:
		s.moveCursor(-1)
	case actions.ActionToggleOne:
		if s.Focus == FocusCommand {
			msg := s.Command.Input
			if msg == "" {
				break
			}
			res.Operations = []actions.Operation{{Kind: actions.OpCommit, Message: msg}}
			res.RefreshChanges = true
			s.Command.Input = ""
			s.Command.Cursor = 0
			break
		}
		if s.Focus != FocusChanges {
			break
		}
		entry, section, ok := s.selectedChange()
		if !ok {
			break
		}
		if section == SectionStaged {
			s.Changes.StickySection = SectionStaged
			res.Operations = []actions.Operation{{Kind: actions.OpUnstagePath, Path: entry.Path}}
		} else {
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []actions.Operation{{Kind: actions.OpStagePath, Path: entry.Path}}
		}
		res.RefreshChanges = true
	case actions.ActionStageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionStaged
			res.Operations = []actions.Operation{{Kind: actions.OpStageAll}}
			res.RefreshChanges = true
		}
	case actions.ActionUnstageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []actions.Operation{{Kind: actions.OpUnstageAll}}
			res.RefreshChanges = true
		}
	case actions.ActionPush:
		res.Operations = []actions.Operation{{Kind: actions.OpPush}}
	}
	s.Clamp()
	return res
}

func (s *AppState) moveCursor(delta int) {
	if s.Focus == FocusGraph {
		s.Graph.Cursor += delta
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

func (s *AppState) selectedChange() (git.ChangeEntry, Section, bool) {
	if len(s.Changes.Rows) == 0 || s.Changes.Cursor < 0 || s.Changes.Cursor >= len(s.Changes.Rows) {
		return git.ChangeEntry{}, "", false
	}
	row := s.Changes.Rows[s.Changes.Cursor]
	if !row.Selectable {
		return git.ChangeEntry{}, "", false
	}
	if row.Section == SectionStaged {
		if row.EntryIndex < 0 || row.EntryIndex >= len(s.Changes.Staged) {
			return git.ChangeEntry{}, "", false
		}
		return s.Changes.Staged[row.EntryIndex], row.Section, true
	}
	if row.EntryIndex < 0 || row.EntryIndex >= len(s.Changes.Unstaged) {
		return git.ChangeEntry{}, "", false
	}
	return s.Changes.Unstaged[row.EntryIndex], row.Section, true
}

func (s *AppState) selectedPath() (string, Section, bool) {
	e, sec, ok := s.selectedChange()
	if !ok {
		return "", "", false
	}
	return e.Path, sec, true
}

func (s *AppState) moveCursorToPath(path string, section Section) bool {
	for i, row := range s.Changes.Rows {
		if !row.Selectable || row.Section != section {
			continue
		}
		if section == SectionStaged {
			if row.EntryIndex >= 0 && row.EntryIndex < len(s.Changes.Staged) && s.Changes.Staged[row.EntryIndex].Path == path {
				s.Changes.Cursor = i
				return true
			}
			continue
		}
		if row.EntryIndex >= 0 && row.EntryIndex < len(s.Changes.Unstaged) && s.Changes.Unstaged[row.EntryIndex].Path == path {
			s.Changes.Cursor = i
			return true
		}
	}
	return false
}

func (s *AppState) moveCursorToSection(section Section) bool {
	for i, row := range s.Changes.Rows {
		if row.Selectable && row.Section == section {
			s.Changes.Cursor = i
			s.Changes.Offset = 0
			return true
		}
	}
	return false
}

func (s *AppState) moveCursorToFirstSelectable() bool {
	for i, row := range s.Changes.Rows {
		if row.Selectable {
			s.Changes.Cursor = i
			s.Changes.Offset = 0
			return true
		}
	}
	s.Changes.Cursor = 0
	s.Changes.Offset = 0
	return false
}

func (s *AppState) snapChangesCursor(dir int) {
	if len(s.Changes.Rows) == 0 {
		s.Changes.Cursor = 0
		return
	}
	if s.Changes.Cursor < 0 {
		s.Changes.Cursor = 0
	}
	if s.Changes.Cursor >= len(s.Changes.Rows) {
		s.Changes.Cursor = len(s.Changes.Rows) - 1
	}
	if s.Changes.Rows[s.Changes.Cursor].Selectable {
		return
	}

	i := s.Changes.Cursor
	for i >= 0 && i < len(s.Changes.Rows) {
		if s.Changes.Rows[i].Selectable {
			s.Changes.Cursor = i
			return
		}
		i += dir
	}
	for i = 0; i < len(s.Changes.Rows); i++ {
		if s.Changes.Rows[i].Selectable {
			s.Changes.Cursor = i
			return
		}
	}
	s.Changes.Cursor = 0
}
