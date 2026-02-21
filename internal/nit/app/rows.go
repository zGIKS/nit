package app

import "nit/internal/nit/git"

func (s *AppState) rebuildChangesSlices() {
	s.Changes.Staged = s.Changes.Staged[:0]
	s.Changes.Unstaged = s.Changes.Unstaged[:0]
	for _, e := range s.Changes.Entries {
		if e.Staged {
			s.Changes.Staged = append(s.Changes.Staged, e)
		}
		if e.Changed || !e.Staged {
			s.Changes.Unstaged = append(s.Changes.Unstaged, e)
		}
	}
}

func (s *AppState) rebuildChangesRows() {
	rows := make([]ChangeRow, 0, len(s.Changes.Entries)+4)
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
