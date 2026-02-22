package state

import (
	"nit/internal/nit/app/input"
	"nit/internal/nit/git"
)

type FocusState int

type Section string

const (
	FocusCommand FocusState = iota
	FocusChanges
	FocusGraph
	FocusCommandLog
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

type CommandState struct {
	Input       string
	Cursor      int
	SelectAll   bool
	Clipboard   string
	ReturnFocus FocusState
}

type CommandLogState struct {
	Cursor int
	Offset int
}

type Viewport struct {
	Width  int
	Height int
}

type AppState struct {
	Focus          FocusState
	Command        CommandState
	Changes        ChangesState
	Graph          GraphState
	CommandLogView CommandLogState
	CommandLog     []string
	Viewport       Viewport
	Keys           input.Keymap
	LastErr        string
	RepoName       string
	BranchName     string
	RepoLabel      string
	BranchLabel    string
	FetchLabel     string
	MenuLabel      string
}

func New(keys input.Keymap) AppState {
	return AppState{
		Focus: FocusChanges,
		Command: CommandState{
			ReturnFocus: FocusChanges,
		},
		Changes: ChangesState{
			StickySection: SectionUnstaged,
		},
		Keys:        keys,
		RepoName:    "loading...",
		BranchName:  "loading...",
		RepoLabel:   "repo",
		BranchLabel: "branch",
		FetchLabel:  "[f] fetch",
		MenuLabel:   "...",
	}
}
