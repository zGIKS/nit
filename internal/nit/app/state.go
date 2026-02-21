package app

import "nit/internal/nit/git"

type FocusState int

type Section string

const (
	FocusChanges FocusState = iota
	FocusGraph
)

const (
	SectionStaged   Section = "staged"
	SectionUnstaged Section = "unstaged"
)

type ChangeRow struct {
	Text       string
	Selectable bool
	Section    Section
	EntryIndex int
}

type ChangesState struct {
	Entries       []git.ChangeEntry
	Staged        []git.ChangeEntry
	Unstaged      []git.ChangeEntry
	Rows          []ChangeRow
	Cursor        int
	Offset        int
	StickySection Section
}

type GraphState struct {
	Lines  []string
	Cursor int
	Offset int
}

type Viewport struct {
	Width  int
	Height int
}

type AppState struct {
	Focus    FocusState
	Changes  ChangesState
	Graph    GraphState
	Viewport Viewport
	Keys     Keymap
	LastErr  string
}

func New(keys Keymap) AppState {
	return AppState{
		Focus: FocusChanges,
		Changes: ChangesState{
			StickySection: SectionUnstaged,
		},
		Keys: keys,
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
	s.Changes.Staged = s.Changes.Staged[:0]
	s.Changes.Unstaged = s.Changes.Unstaged[:0]
	for _, e := range entries {
		if e.Staged {
			s.Changes.Staged = append(s.Changes.Staged, e)
		}
		if e.Changed || !e.Staged {
			s.Changes.Unstaged = append(s.Changes.Unstaged, e)
		}
	}

	rows := make([]ChangeRow, 0, len(entries)+4)
	if len(s.Changes.Staged) > 0 {
		rows = append(rows, ChangeRow{Text: "Staged Changes"})
		for i, e := range s.Changes.Staged {
			rows = append(rows, ChangeRow{
				Text:       "  " + codeForStaged(e) + "  " + e.Path,
				Selectable: true,
				Section:    SectionStaged,
				EntryIndex: i,
			})
		}
	}
	if len(s.Changes.Unstaged) > 0 {
		rows = append(rows, ChangeRow{Text: "Changes"})
		for i, e := range s.Changes.Unstaged {
			rows = append(rows, ChangeRow{
				Text:       "  " + codeForUnstaged(e) + "  " + e.Path,
				Selectable: true,
				Section:    SectionUnstaged,
				EntryIndex: i,
			})
		}
	}
	if len(rows) == 0 {
		rows = []ChangeRow{{Text: "Working tree clean."}}
	}
	s.Changes.Rows = rows

	if hadPrev && s.moveCursorToPath(prevPath, prevSection) {
		s.Clamp()
		return
	}
	if !s.moveCursorToSection(s.Changes.StickySection) {
		s.moveCursorToFirstSelectable()
	}
	s.Clamp()
}

func (s *AppState) Apply(action Action) ApplyResult {
	res := ApplyResult{}
	switch action {
	case ActionQuit:
		res.Quit = true
	case ActionTogglePanel:
		if s.Focus == FocusChanges {
			s.Focus = FocusGraph
		} else {
			s.Focus = FocusChanges
			s.snapChangesCursor(1)
		}
	case ActionMoveDown:
		s.moveCursor(1)
	case ActionMoveUp:
		s.moveCursor(-1)
	case ActionToggleOne:
		if s.Focus != FocusChanges {
			break
		}
		entry, section, ok := s.selectedChange()
		if !ok {
			break
		}
		if section == SectionStaged {
			s.Changes.StickySection = SectionStaged
			res.Operations = []Operation{{Kind: OpUnstagePath, Path: entry.Path}}
		} else {
			// Keep staging flow on unstaged while there are items left.
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []Operation{{Kind: OpStagePath, Path: entry.Path}}
		}
		res.RefreshChanges = true
	case ActionStageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionStaged
			res.Operations = []Operation{{Kind: OpStageAll}}
			res.RefreshChanges = true
		}
	case ActionUnstageAll:
		if s.Focus == FocusChanges {
			s.Changes.StickySection = SectionUnstaged
			res.Operations = []Operation{{Kind: OpUnstageAll}}
			res.RefreshChanges = true
		}
	}
	s.Clamp()
	return res
}

func (s *AppState) moveCursor(delta int) {
	if s.Focus == FocusGraph {
		s.Graph.Cursor += delta
		s.Clamp()
		return
	}
	s.Changes.Cursor += delta
	if delta >= 0 {
		s.snapChangesCursor(1)
	} else {
		s.snapChangesCursor(-1)
	}
	s.Clamp()
}

func (s *AppState) selectedChange() (git.ChangeEntry, Section, bool) {
	if len(s.Changes.Rows) == 0 || s.Changes.Cursor < 0 || s.Changes.Cursor >= len(s.Changes.Rows) {
		return git.ChangeEntry{}, "", false
	}
	row := s.Changes.Rows[s.Changes.Cursor]
	if !row.Selectable {
		return git.ChangeEntry{}, "", false
	}
	if row.Section == SectionStaged {
		if row.EntryIndex < 0 || row.EntryIndex >= len(s.Changes.Staged) {
			return git.ChangeEntry{}, "", false
		}
		return s.Changes.Staged[row.EntryIndex], row.Section, true
	}
	if row.EntryIndex < 0 || row.EntryIndex >= len(s.Changes.Unstaged) {
		return git.ChangeEntry{}, "", false
	}
	return s.Changes.Unstaged[row.EntryIndex], row.Section, true
}

func (s *AppState) selectedPath() (string, Section, bool) {
	e, sec, ok := s.selectedChange()
	if !ok {
		return "", "", false
	}
	return e.Path, sec, true
}

func (s *AppState) moveCursorToPath(path string, section Section) bool {
	for i, row := range s.Changes.Rows {
		if !row.Selectable || row.Section != section {
			continue
		}
		if section == SectionStaged {
			if row.EntryIndex >= 0 && row.EntryIndex < len(s.Changes.Staged) && s.Changes.Staged[row.EntryIndex].Path == path {
				s.Changes.Cursor = i
				return true
			}
			continue
		}
		if row.EntryIndex >= 0 && row.EntryIndex < len(s.Changes.Unstaged) && s.Changes.Unstaged[row.EntryIndex].Path == path {
			s.Changes.Cursor = i
			return true
		}
	}
	return false
}

func (s *AppState) moveCursorToSection(section Section) bool {
	for i, row := range s.Changes.Rows {
		if row.Selectable && row.Section == section {
			s.Changes.Cursor = i
			s.Changes.Offset = 0
			return true
		}
	}
	return false
}

func (s *AppState) moveCursorToFirstSelectable() bool {
	for i, row := range s.Changes.Rows {
		if row.Selectable {
			s.Changes.Cursor = i
			s.Changes.Offset = 0
			return true
		}
	}
	s.Changes.Cursor = 0
	s.Changes.Offset = 0
	return false
}

func (s *AppState) snapChangesCursor(dir int) {
	if len(s.Changes.Rows) == 0 {
		s.Changes.Cursor = 0
		return
	}
	if s.Changes.Cursor < 0 {
		s.Changes.Cursor = 0
	}
	if s.Changes.Cursor >= len(s.Changes.Rows) {
		s.Changes.Cursor = len(s.Changes.Rows) - 1
	}
	if s.Changes.Rows[s.Changes.Cursor].Selectable {
		return
	}

	i := s.Changes.Cursor
	for i >= 0 && i < len(s.Changes.Rows) {
		if s.Changes.Rows[i].Selectable {
			s.Changes.Cursor = i
			return
		}
		i += dir
	}
	for i = 0; i < len(s.Changes.Rows); i++ {
		if s.Changes.Rows[i].Selectable {
			s.Changes.Cursor = i
			return
		}
	}
	s.Changes.Cursor = 0
}

func codeForStaged(e git.ChangeEntry) string {
	if e.X == '?' {
		return "U"
	}
	if e.X != ' ' {
		return string(e.X)
	}
	if e.Staged {
		return "M"
	}
	return "-"
}

func codeForUnstaged(e git.ChangeEntry) string {
	if e.X == '?' {
		return "U"
	}
	if e.Y != ' ' {
		return string(e.Y)
	}
	if !e.Staged {
		return "M"
	}
	return "-"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
