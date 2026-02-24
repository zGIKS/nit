package state

import "github.com/zGIKS/nit/internal/nit/app/input"

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
		MenuOffset:               0,
		MenuSubHoverIndex:        -1,
		MenuSubOffset:            0,
		BranchCreateHoverIndex:   -1,
		RepoName:                 "loading...",
		BranchName:               "loading...",
		RepoLabel:                "repo",
		BranchLabel:              "branch",
		RepoBranchSeparator:      "->",
		FetchLabel:               "[f] fetch",
		MenuLabel:                "...",
		MenuChevron:              "›",
		MenuSelectionIndicator:   ">",
		BranchSourceSelectedMark: "✓",
		BranchCreateTitle:        "Create a branch",
		BranchCreateEnterHint:    "Enter: create branch",
		BranchCreatePushHint:     "Ctrl+b: create and push to origin",
		BranchCreateNameLabel:    "New branch name",
		BranchCreateSourceLabel:  "Source",
	}
}
