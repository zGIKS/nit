package state

import (
	"strings"
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

func (s *AppState) SetUISymbols(branchSourceSelectedMark, menuChevron, menuSelectionIndicator string) {
	if strings.TrimSpace(branchSourceSelectedMark) != "" {
		s.BranchSourceSelectedMark = strings.TrimSpace(branchSourceSelectedMark)
	}
	if strings.TrimSpace(menuChevron) != "" {
		s.MenuChevron = strings.TrimSpace(menuChevron)
	}
	if strings.TrimSpace(menuSelectionIndicator) != "" {
		s.MenuSelectionIndicator = strings.TrimSpace(menuSelectionIndicator)
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
