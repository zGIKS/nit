package state

import (
	"strings"

	"nit/internal/nit/git"
)

func (s *AppState) SetViewport(width, height int) {
	s.Viewport.Width = width
	s.Viewport.Height = height
	s.Clamp()
}

func (s *AppState) SetError(errMsg string) {
	s.LastErr = errMsg
}

func (s *AppState) SetRepoSummary(repo, branch string) {
	if strings.TrimSpace(repo) != "" {
		s.RepoName = strings.TrimSpace(repo)
	}
	if strings.TrimSpace(branch) != "" {
		s.BranchName = strings.TrimSpace(branch)
	}
}

func (s *AppState) SetTopBarLabels(repo, branch, fetch, menu string) {
	if strings.TrimSpace(repo) != "" {
		s.RepoLabel = strings.TrimSpace(repo)
	}
	if strings.TrimSpace(branch) != "" {
		s.BranchLabel = strings.TrimSpace(branch)
	}
	if strings.TrimSpace(fetch) != "" {
		s.FetchLabel = strings.TrimSpace(fetch)
	}
	if strings.TrimSpace(menu) != "" {
		s.MenuLabel = strings.TrimSpace(menu)
	}
}

func (s *AppState) SetUISymbols(branchSourceSelectedMark string) {
	if strings.TrimSpace(branchSourceSelectedMark) != "" {
		s.BranchSourceSelectedMark = strings.TrimSpace(branchSourceSelectedMark)
	}
}

func (s *AppState) SetUIText(branchCreateTitle, branchCreateEnterHint, branchCreatePushHint, branchCreateNameLabel, branchCreateSourceLabel string) {
	if strings.TrimSpace(branchCreateTitle) != "" {
		s.BranchCreateTitle = strings.TrimSpace(branchCreateTitle)
	}
	if strings.TrimSpace(branchCreateEnterHint) != "" {
		s.BranchCreateEnterHint = strings.TrimSpace(branchCreateEnterHint)
	}
	if strings.TrimSpace(branchCreatePushHint) != "" {
		s.BranchCreatePushHint = strings.TrimSpace(branchCreatePushHint)
	}
	if strings.TrimSpace(branchCreateNameLabel) != "" {
		s.BranchCreateNameLabel = strings.TrimSpace(branchCreateNameLabel)
	}
	if strings.TrimSpace(branchCreateSourceLabel) != "" {
		s.BranchCreateSourceLabel = strings.TrimSpace(branchCreateSourceLabel)
	}
}

func (s *AppState) SetGraph(lines []string) {
	if len(lines) == 0 {
		lines = []string{"No commits to display."}
	}
	s.Graph.Lines = lines
	if s.Graph.Cursor >= len(s.Graph.Lines) {
		s.Graph.Cursor = max(0, len(s.Graph.Lines)-1)
	}
	s.Clamp()
}

func (s *AppState) SetBranches(lines []string) {
	if len(lines) == 0 {
		lines = []string{"No local branches."}
	}
	s.Branches.Lines = lines
	if s.Branches.Cursor < 0 {
		s.Branches.Cursor = 0
	}
	if s.Branches.Cursor >= len(s.Branches.Lines) {
		s.Branches.Cursor = max(0, len(s.Branches.Lines)-1)
	}
	page := s.branchesPageSize()
	if s.Branches.Cursor < s.Branches.Offset {
		s.Branches.Offset = s.Branches.Cursor
	}
	if s.Branches.Cursor >= s.Branches.Offset+page {
		s.Branches.Offset = s.Branches.Cursor - page + 1
	}
	maxOffset := max(0, len(s.Branches.Lines)-page)
	if s.Branches.Offset > maxOffset {
		s.Branches.Offset = maxOffset
	}
	if s.Branches.Offset < 0 {
		s.Branches.Offset = 0
	}
	s.syncBranchCreateSources()
}

func (s *AppState) SetChanges(entries []git.ChangeEntry) {
	prevPath, prevSection, hadPrev := s.selectedPath()
	if s.Changes.StickySection == "" {
		s.Changes.StickySection = SectionUnstaged
	}

	s.Changes.Entries = entries
	s.rebuildChangesSlices()
	s.rebuildChangesRows()

	if hadPrev && s.moveCursorToPath(prevPath, prevSection) {
		s.Clamp()
		return
	}
	if !s.moveCursorToSection(s.Changes.StickySection) {
		s.moveCursorToFirstSelectable()
	}
	s.Clamp()
}

func (s *AppState) AppendCommandText(text string) {
	appendTextInput(&s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll, text)
}

func (s *AppState) BackspaceCommandText() {
	backspaceTextInput(&s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) DeleteCommandText() {
	deleteTextInput(&s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) MoveCommandCursorLeft() {
	moveTextInputCursorLeft(&s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) MoveCommandCursorRight() {
	moveTextInputCursorRight(s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) MoveCommandCursorToStart() {
	moveTextInputCursorHome(&s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) MoveCommandCursorToEnd() {
	moveTextInputCursorEnd(s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}

func (s *AppState) ExitCommandFocus() {
	target := s.Command.ReturnFocus
	if target == FocusCommand {
		target = FocusChanges
	}
	s.Focus = target
	s.Command.SelectAll = false
	if s.Focus == FocusChanges {
		s.snapChangesCursor(1)
	}
	s.Clamp()
}

func (s *AppState) SelectAllCommandText() {
	selectAllTextInput(s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}

func (s AppState) SelectedCommandText() string {
	if s.Command.SelectAll {
		return s.Command.Input
	}
	return ""
}

func (s *AppState) SetCommandClipboard(text string) {
	s.Command.Clipboard = text
}

func (s AppState) CommandClipboard() string {
	return s.Command.Clipboard
}

func (s *AppState) AddCommandLog(cmd string) {
	s.CommandLog = append(s.CommandLog, cmd)
	if len(s.CommandLog) > 100 {
		s.CommandLog = s.CommandLog[len(s.CommandLog)-100:]
	}
	if len(s.CommandLog) == 0 {
		s.CommandLogView.Cursor = 0
		s.CommandLogView.Offset = 0
		return
	}
	s.CommandLogView.Cursor = len(s.CommandLog) - 1
	page := s.commandLogPageSize()
	s.CommandLogView.Offset = max(0, len(s.CommandLog)-page)
}

func (s *AppState) DeleteCommandSelection() {
	clearSelectedText(&s.Command.Input, &s.Command.Cursor, &s.Command.SelectAll)
}
