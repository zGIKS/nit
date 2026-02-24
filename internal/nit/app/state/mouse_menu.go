package state

import "github.com/zGIKS/nit/internal/nit/app/actions"

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
	return s.MenuActivateIndex(idx)
}

func (s *AppState) MenuActivateIndex(idx int) (actions.Action, bool) {
	items := s.MenuItems()
	if idx < 0 || idx >= len(items) || items[idx].Separator {
		return actions.ActionNone, false
	}
	item := items[idx]
	switch item.Label {
	case "Pull":
		s.CloseMenu()
		return actions.ActionPull, true
	case "Fetch":
		s.CloseMenu()
		return actions.ActionFetch, true
	case "Push":
		s.CloseMenu()
		return actions.ActionPush, true
	case "Commit":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Changes":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Pull, Push":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Branch":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Remote":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Stash":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	case "Tags":
		s.OpenSubmenuForMenuIndex(idx)
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	default:
		// Keep menu open for category rows until their submenus are implemented.
		s.MenuSubmenuKind = ""
		s.MenuSubHoverIndex = -1
		s.MenuHoverIndex = idx
		return actions.ActionNone, false
	}
}

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
			return actions.ActionNone, false, true // placeholder
		}
	case "pull_push":
		switch item.Label {
		case "Sync":
			s.CloseMenu()
			return actions.ActionPull, true, true // placeholder: sync currently maps to pull
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
			return actions.ActionNone, false, true // placeholders
		}
	case "branch":
		switch item.Label {
		case "Create Branch...", "Create Branch From...":
			s.CloseMenu()
			s.OpenBranchCreate()
			return actions.ActionNone, false, true
		case "Merge...", "Rebase Branch...", "Rename Branch...", "Delete Branch...", "Delete Remote Branch...", "Publish Branch...":
			s.CloseMenu()
			return actions.ActionNone, false, true // placeholders
		}
	case "remote":
		switch item.Label {
		case "Add Remote...", "Remove Remote":
			s.CloseMenu()
			return actions.ActionNone, false, true // placeholders
		}
	case "stash":
		switch item.Label {
		case "Stash", "Stash (Include Untracked)", "Stash Staged",
			"Apply Latest Stash", "Apply Stash...",
			"Pop Latest Stash", "Pop Stash...",
			"Drop Stash...", "Drop All Stashes...",
			"View Stash...":
			s.CloseMenu()
			return actions.ActionNone, false, true // placeholders
		}
	case "tags":
		switch item.Label {
		case "Create Tag...", "Delete Tag...", "Delete Remote Tag...", "Push Tags":
			s.CloseMenu()
			return actions.ActionNone, false, true // placeholders
		}
	}
	s.CloseMenu()
	// Placeholder until commit submenu actions are implemented.
	return actions.ActionNone, false, true
}

func (s *AppState) ToggleMenuClick(x, y int) bool {
	if mx, my, mw, mh := s.MenuButtonRect(); x >= mx && x < mx+mw && y >= my && y < my+mh {
		s.CloseBranchCreate()
		s.ToggleMenu()
		return true
	}
	return false
}

func (s *AppState) ToggleBranchCreateClick(x, y int) bool {
	if bx, by, bw, bh := s.BranchButtonRect(); x >= bx && x < bx+bw && y >= by && y < by+bh {
		s.CloseMenu()
		s.ToggleBranchCreate()
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
	if sx, sy, sw, sh := s.MenuSubmenuRect(); sw > 0 && sh > 0 && x >= sx && x < sx+sw && y >= sy && y < sy+sh {
		return
	}
	s.CloseMenu()
}

func (s *AppState) CloseTopMenusOnOutsideClick(x, y int) {
	s.CloseMenuOnOutsideClick(x, y)
	s.CloseBranchCreateOnOutsideClick(x, y)
}
