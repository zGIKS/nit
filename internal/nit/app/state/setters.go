package state

import (
	"strings"
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
