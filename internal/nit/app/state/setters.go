package state

import (
	"strings"

	"github.com/zGIKS/nit/internal/nit/git"
)

func setIfNotBlank(dst *string, src string) {
	if v := strings.TrimSpace(src); v != "" {
		*dst = v
	}
}

func (s *AppState) SetViewport(width, height int) {
	s.Viewport.Width = width
	s.Viewport.Height = height
	s.Clamp()
}

func (s *AppState) SetError(errMsg string) {
	s.LastErr = errMsg
}

func (s *AppState) SetRepoSummary(repo, branch string) {
	setIfNotBlank(&s.RepoName, repo)
	setIfNotBlank(&s.BranchName, branch)
}

func (s *AppState) SetTopBarLabels(repo, branch, fetch, menu string) {
	setIfNotBlank(&s.RepoLabel, repo)
	setIfNotBlank(&s.BranchLabel, branch)
	setIfNotBlank(&s.FetchLabel, fetch)
	setIfNotBlank(&s.MenuLabel, menu)
}

func (s *AppState) SetRepoBranchSeparator(label string) {
	setIfNotBlank(&s.RepoBranchSeparator, label)
}

func (s *AppState) SetUISymbols(branchSourceSelectedMark, menuChevron, menuSelectionIndicator string) {
	setIfNotBlank(&s.BranchSourceSelectedMark, branchSourceSelectedMark)
	setIfNotBlank(&s.MenuChevron, menuChevron)
	setIfNotBlank(&s.MenuSelectionIndicator, menuSelectionIndicator)
}

func (s *AppState) SetUIText(branchCreateTitle, branchCreateEnterHint, branchCreatePushHint, branchCreateNameLabel, branchCreateSourceLabel string) {
	setIfNotBlank(&s.BranchCreateTitle, branchCreateTitle)
	setIfNotBlank(&s.BranchCreateEnterHint, branchCreateEnterHint)
	setIfNotBlank(&s.BranchCreatePushHint, branchCreatePushHint)
	setIfNotBlank(&s.BranchCreateNameLabel, branchCreateNameLabel)
	setIfNotBlank(&s.BranchCreateSourceLabel, branchCreateSourceLabel)
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


