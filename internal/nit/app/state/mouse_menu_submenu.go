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
		case "Commit", "Commit Staged":
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
		case "Abort Rebase":
			s.CloseMenu()
			return actions.ActionAbortRebase, true, true
		case "Commit (Amend)", "Commit Staged (Amend)":
			s.CloseMenu()
			s.PrepareCommandCommit(false, true, false)
			return actions.ActionNone, false, true
		case "Commit All (Amend)":
			s.CloseMenu()
			s.PrepareCommandCommit(true, true, false)
			return actions.ActionNone, false, true
		case "Commit (Signed Off)", "Commit Staged (Signed Off)":
			s.CloseMenu()
			s.PrepareCommandCommit(false, false, true)
			return actions.ActionNone, false, true
		case "Commit All (Signed Off)":
			s.CloseMenu()
			s.PrepareCommandCommit(true, false, true)
			return actions.ActionNone, false, true
		}
	case "changes":
		switch item.Label {
		case "Stage All Changes":
			s.CloseMenu()
			return actions.ActionStageAll, true, true
		case "Unstage All Changes":
			s.CloseMenu()
			return actions.ActionUnstageAll, true, true
		case "Discard All Changes":
			s.CloseMenu()
			return actions.ActionDiscardAll, true, true
		}
	case "pull_push":
		switch item.Label {
		case "Sync":
			s.CloseMenu()
			return actions.ActionPull, true, true
		case "Pull":
			s.CloseMenu()
			return actions.ActionPull, true, true
		case "Push":
			s.CloseMenu()
			return actions.ActionPush, true, true
		case "Fetch":
			s.CloseMenu()
			return actions.ActionFetch, true, true
		case "Pull (Rebase)", "Pull from...", "Push to...", "Fetch (Prune)", "Fetch From All Remotes":
			s.CloseMenu()
			return actions.ActionNone, false, true
		}
	case "branch":
		switch item.Label {
		case "Create Branch...", "Create Branch From...":
			s.CloseMenu()
			s.OpenBranchCreate()
			return actions.ActionNone, false, true
		case "Merge...", "Rebase Branch...", "Rename Branch...", "Delete Branch...", "Delete Remote Branch...", "Publish Branch...":
			s.CloseMenu()
			return actions.ActionNone, false, true
		}
	case "remote":
		switch item.Label {
		case "Add Remote...", "Remove Remote":
			s.CloseMenu()
			return actions.ActionNone, false, true
		}
	case "stash":
		switch item.Label {
		case "Stash", "Stash (Include Untracked)", "Stash Staged",
			"Apply Latest Stash", "Apply Stash...",
			"Pop Latest Stash", "Pop Stash...",
			"Drop Stash...", "Drop All Stashes...",
			"View Stash...":
			s.CloseMenu()
			return actions.ActionNone, false, true
		}
	case "tags":
		switch item.Label {
		case "Create Tag...", "Delete Tag...", "Delete Remote Tag...", "Push Tags":
			s.CloseMenu()
			return actions.ActionNone, false, true
		}
	}
	s.CloseMenu()
	return actions.ActionNone, false, true
}
