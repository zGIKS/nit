package state

import (
	"strings"

	"github.com/mattn/go-runewidth"
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

func (s *AppState) TopBarActionAt(x, y int) (actions.Action, bool) {
	// Top bar is the first 3 rows of the command pane.
	if y < 0 || y >= 3 || x < 0 {
		return actions.ActionNone, false
	}

	totalW := max(40, s.Viewport.Width)
	repoName := s.RepoName
	if repoName == "" {
		repoName = "unknown"
	}
	branchName := s.BranchName
	if branchName == "" {
		branchName = "-"
	}

	repoText := strings.TrimSpace(s.RepoLabel + " " + repoName)
	branchText := strings.TrimSpace(s.BranchLabel + " " + branchName)
	fetchText := strings.TrimSpace(s.FetchLabel)
	menuText := strings.TrimSpace(s.MenuLabel)

	repoW := max(16, runewidth.StringWidth(repoText)+4)
	branchW := max(16, runewidth.StringWidth(branchText)+4)
	fetchW := max(14, runewidth.StringWidth(fetchText)+4)
	menuW := max(8, runewidth.StringWidth(menuText)+4)
	minRepoW, minBranchW, minFetchW, minMenuW := 14, 12, 10, 8

	totalNeeded := repoW + branchW + fetchW + menuW + 3
	overflow := totalNeeded - totalW
	shrink := func(w *int, minW int) {
		if overflow <= 0 {
			return
		}
		can := *w - minW
		if can <= 0 {
			return
		}
		d := min(can, overflow)
		*w -= d
		overflow -= d
	}
	shrink(&repoW, minRepoW)
	shrink(&branchW, minBranchW)
	shrink(&fetchW, minFetchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		repoW = max(minRepoW, repoW-overflow)
	}

	rightTopW := branchW + fetchW + menuW + 2
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}

	branchX := repoW + 1 + gapW + 1
	fetchX := branchX + branchW + 1
	menuX := fetchX + fetchW + 1
	_ = menuX // reserved for future click actions

	if x >= fetchX && x < fetchX+fetchW {
		return actions.ActionFetch, true
	}
	return actions.ActionNone, false
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
