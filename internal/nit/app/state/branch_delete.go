package state

import "strings"

func (s *AppState) OpenBranchDeleteConfirm() bool {
	branch, ok := s.SelectedBranchName()
	if !ok || strings.TrimSpace(branch) == strings.TrimSpace(s.BranchName) {
		return false
	}
	s.BranchDeleteConfirmOpen = true
	s.BranchDeleteBranch = branch
	s.BranchDeleteChoice = 1
	s.CloseMenu()
	s.CloseBranchCreate()
	return true
}

func (s *AppState) CloseBranchDeleteConfirm() {
	s.BranchDeleteConfirmOpen = false
	s.BranchDeleteBranch = ""
	s.BranchDeleteChoice = 1
}

func (s *AppState) MoveBranchDeleteChoice(delta int) {
	if delta == 0 {
		return
	}
	s.BranchDeleteChoice = (s.BranchDeleteChoice + delta) % 2
	if s.BranchDeleteChoice < 0 {
		s.BranchDeleteChoice += 2
	}
}

func (s AppState) BranchDeleteConfirmed() bool {
	return s.BranchDeleteChoice == 0
}
