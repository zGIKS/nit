package state

func (s *AppState) OpenBranchCreate() {
	s.BranchCreateOpen = true
	s.CloseMenu()
	s.BranchCreateHoverIndex = -1
	s.BranchCreateSelectAll = false
	s.syncBranchCreateSources()
	moveTextInputCursorEnd(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) CloseBranchCreate() {
	s.BranchCreateOpen = false
	s.BranchCreateHoverIndex = -1
}

func (s *AppState) ToggleBranchCreate() {
	if s.BranchCreateOpen {
		s.CloseBranchCreate()
		return
	}
	s.OpenBranchCreate()
}

func (s *AppState) BranchCreateAppendText(text string) {
	appendTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll, text)
}

func (s *AppState) BranchCreateBackspace() {
	backspaceTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateDelete() {
	deleteTextInput(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorLeft() {
	moveTextInputCursorLeft(&s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorRight() {
	moveTextInputCursorRight(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorHome() {
	moveTextInputCursorHome(&s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateCursorEnd() {
	moveTextInputCursorEnd(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s *AppState) BranchCreateSelectAllText() {
	selectAllTextInput(s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}

func (s AppState) SelectedBranchCreateText() string {
	if s.BranchCreateSelectAll {
		return s.BranchCreateName
	}
	return ""
}

func (s *AppState) DeleteBranchCreateSelection() {
	clearSelectedText(&s.BranchCreateName, &s.BranchCreateCursor, &s.BranchCreateSelectAll)
}
