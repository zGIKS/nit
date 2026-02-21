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
