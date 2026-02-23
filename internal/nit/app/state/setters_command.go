package state

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
