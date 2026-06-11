package state

import "github.com/zGIKS/nit/internal/nit/app/actions"

func (s *AppState) MenuSubmenuClickActionAt(x, y int) (actions.Action, bool, bool) {
	idx, ok := s.MenuSubmenuItemIndexAt(x, y)
	if !ok {
		return actions.ActionNone, false, false
	}
	return s.MenuSubmenuActivateIndex(idx)
}

func (s *AppState) MenuSubmenuActivateIndex(idx int) (actions.Action, bool, bool) {
	items := s.MenuSubmenuItems()
	if idx < 0 || idx >= len(items) || items[idx].Separator {
		return actions.ActionNone, false, false
	}
	item := items[idx]
	switch s.MenuSubmenuKind {
	case "commit":
		switch item.Label {
		case "Commit Staged":
			s.CloseMenu()
			s.PrepareCommandCommit(false, false, false)
			return actions.ActionNone, false, true
		case "Commit All":
			s.CloseMenu()
			s.PrepareCommandCommit(true, false, false)
			return actions.ActionNone, false, true
		case "Undo Last Commit":
			s.CloseMenu()
			return actions.ActionUndoLastCommit, true, true
		}
	case "changes":
		switch item.Label {
		case "Stage All Changes":
			s.CloseMenu()
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
			return actions.ActionStageAll, true, true
		case "Unstage All Changes":
			s.CloseMenu()
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
			return actions.ActionUnstageAll, true, true
		case "Discard All Changes":
			s.CloseMenu()
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
			return actions.ActionDiscardAll, true, true
		}
	}
	s.CloseMenu()
	return actions.ActionNone, false, true
}
