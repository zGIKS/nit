package state

import "nit/internal/nit/git"

func (s *AppState) SetViewport(width, height int) {
	s.Viewport.Width = width
	s.Viewport.Height = height
	s.Clamp()
}

func (s *AppState) SetError(errMsg string) {
	s.LastErr = errMsg
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
	if text == "" {
		return
	}
	if s.Command.SelectAll {
		s.Command.Input = ""
		s.Command.Cursor = 0
		s.Command.SelectAll = false
	}
	r := []rune(s.Command.Input)
	if s.Command.Cursor < 0 {
		s.Command.Cursor = 0
	}
	if s.Command.Cursor > len(r) {
		s.Command.Cursor = len(r)
	}
	insert := []rune(text)
	out := make([]rune, 0, len(r)+len(insert))
	out = append(out, r[:s.Command.Cursor]...)
	out = append(out, insert...)
	out = append(out, r[s.Command.Cursor:]...)
	s.Command.Input = string(out)
	s.Command.Cursor += len(insert)
}

func (s *AppState) BackspaceCommandText() {
	if s.Command.SelectAll {
		s.Command.Input = ""
		s.Command.Cursor = 0
		s.Command.SelectAll = false
		return
	}
	r := []rune(s.Command.Input)
	if len(r) == 0 || s.Command.Cursor <= 0 {
		return
	}
	if s.Command.Cursor > len(r) {
		s.Command.Cursor = len(r)
	}
	out := make([]rune, 0, len(r)-1)
	out = append(out, r[:s.Command.Cursor-1]...)
	out = append(out, r[s.Command.Cursor:]...)
	s.Command.Input = string(out)
	s.Command.Cursor--
}

func (s *AppState) DeleteCommandText() {
	if s.Command.SelectAll {
		s.Command.Input = ""
		s.Command.Cursor = 0
		s.Command.SelectAll = false
		return
	}
	r := []rune(s.Command.Input)
	if len(r) == 0 {
		return
	}
	if s.Command.Cursor < 0 {
		s.Command.Cursor = 0
	}
	if s.Command.Cursor >= len(r) {
		return
	}
	out := make([]rune, 0, len(r)-1)
	out = append(out, r[:s.Command.Cursor]...)
	out = append(out, r[s.Command.Cursor+1:]...)
	s.Command.Input = string(out)
}

func (s *AppState) MoveCommandCursorLeft() {
	s.Command.SelectAll = false
	if s.Command.Cursor > 0 {
		s.Command.Cursor--
	}
}

func (s *AppState) MoveCommandCursorRight() {
	s.Command.SelectAll = false
	r := []rune(s.Command.Input)
	if s.Command.Cursor < len(r) {
		s.Command.Cursor++
	}
}

func (s *AppState) MoveCommandCursorToStart() {
	s.Command.SelectAll = false
	s.Command.Cursor = 0
}

func (s *AppState) MoveCommandCursorToEnd() {
	s.Command.SelectAll = false
	s.Command.Cursor = len([]rune(s.Command.Input))
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
	if s.Command.Input == "" {
		s.Command.SelectAll = false
		s.Command.Cursor = 0
		return
	}
	s.Command.SelectAll = true
	s.Command.Cursor = len([]rune(s.Command.Input))
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
}

func (s *AppState) DeleteCommandSelection() {
	if !s.Command.SelectAll {
		return
	}
	s.Command.Input = ""
	s.Command.Cursor = 0
	s.Command.SelectAll = false
}

func (s AppState) CommandLineWithCaret() string {
	if s.Command.SelectAll && s.Command.Input != "" {
		return "[" + s.Command.Input + "]"
	}
	r := []rune(s.Command.Input)
	cursor := s.Command.Cursor
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(r) {
		cursor = len(r)
	}
	out := make([]rune, 0, len(r)+1)
	out = append(out, r[:cursor]...)
	out = append(out, '|')
	out = append(out, r[cursor:]...)
	return string(out)
}
