package state

import (
	"nit/internal/nit/app/actions"
)

func (s *AppState) Apply(action actions.Action) actions.ApplyResult {
	res := actions.ApplyResult{}
	switch action {
	case actions.ActionQuit:
		res.Quit = true
	case actions.ActionFocusCommand:
		if s.Focus != FocusCommand {
			s.Command.ReturnFocus = s.Focus
		}
		s.Focus = FocusCommand
		s.Command.SelectAll = false
		s.MoveCommandCursorToEnd()
	case actions.ActionTogglePanel:
		switch s.Focus {
		case FocusCommand:
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
		case FocusChanges:
			s.Focus = FocusGraph
		case FocusGraph:
			s.Focus = FocusBranches
		case FocusBranches:
			s.Focus = FocusCommandLog
		default:
			s.Command.ReturnFocus = FocusBranches
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
			if len(s.Changes.Staged) == 0 {
				s.SetError("nothing staged to commit")
				break
			}
			res.Operations = []actions.Operation{{Kind: actions.OpCommit, Message: msg}}
			res.RefreshChanges = true
			res.RefreshGraph = true
			s.Command.Input = ""
			s.Command.Cursor = 0
			s.Command.SelectAll = false
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
	case actions.ActionPull:
		res.Operations = []actions.Operation{{Kind: actions.OpPull}}
		res.RefreshChanges = true
		res.RefreshGraph = true
	case actions.ActionFetch:
		res.Operations = []actions.Operation{{Kind: actions.OpFetch}}
		res.RefreshGraph = true
	case actions.ActionPush:
		if !s.canPush() {
			s.SetError("nothing to push")
			break
		}
		res.Operations = []actions.Operation{{Kind: actions.OpPush}}
		res.RefreshGraph = true
	}
	s.Clamp()
	return res
}

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
