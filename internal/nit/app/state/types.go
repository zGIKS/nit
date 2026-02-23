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
	FocusBranches
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

type BranchesState struct {
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
	Focus                    FocusState
	Command                  CommandState
	Changes                  ChangesState
	Graph                    GraphState
	Branches                 BranchesState
	CommandLogView           CommandLogState
	CommandLog               []string
	Viewport                 Viewport
	Keys                     input.Keymap
	LastErr                  string
	MenuOpen                 bool
	MenuHoverIndex           int
	HoverFetch               bool
	HoverMenu                bool
	HoverBranch              bool
	RepoName                 string
	BranchName               string
	RepoLabel                string
	BranchLabel              string
	FetchLabel               string
	MenuLabel                string
	BranchSourceSelectedMark string
	BranchCreateOpen         bool
	BranchCreateName         string
	BranchCreateCursor       int
	BranchCreateSelectAll    bool
	BranchCreateSource       string
	BranchCreateSourceList   []string
	BranchCreateSourceOffset int
	BranchCreateHoverIndex   int
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
		Branches: BranchesState{
			Lines: []string{"Loading branches..."},
		},
		Keys:                     keys,
		MenuHoverIndex:           -1,
		BranchCreateHoverIndex:   -1,
		RepoName:                 "loading...",
		BranchName:               "loading...",
		RepoLabel:                "repo",
		BranchLabel:              "branch",
		FetchLabel:               "[f] fetch",
		MenuLabel:                "...",
		BranchSourceSelectedMark: "✓",
	}
}
