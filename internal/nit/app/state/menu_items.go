package state

import "github.com/mattn/go-runewidth"

// DropdownMenuItem represents a single item in a dropdown menu.
type DropdownMenuItem struct {
	Label      string
	HasChevron bool
	Separator  bool
}

var dropdownMenuItems = []DropdownMenuItem{
	{Label: "Pull"},
	{Label: "Fetch"},
	{Separator: true},
	{Label: "Commit", HasChevron: true},
	{Label: "Changes", HasChevron: true},
	{Label: "Pull, Push", HasChevron: true},
	{Label: "Branch", HasChevron: true},
	{Label: "Remote", HasChevron: true},
	{Label: "Stash", HasChevron: true},
	{Label: "Tags", HasChevron: true},
}

var commitDropdownMenuItems = []DropdownMenuItem{
	{Label: "Commit"},
	{Label: "Commit Staged"},
	{Label: "Commit All"},
	{Label: "Undo Last Commit"},
	{Label: "Abort Rebase"},
	{Separator: true},
	{Label: "Commit (Amend)"},
	{Label: "Commit Staged (Amend)"},
	{Label: "Commit All (Amend)"},
	{Separator: true},
	{Label: "Commit (Signed Off)"},
	{Label: "Commit Staged (Signed Off)"},
	{Label: "Commit All (Signed Off)"},
}

var changesDropdownMenuItems = []DropdownMenuItem{
	{Label: "Stage All Changes"},
	{Label: "Unstage All Changes"},
	{Label: "Discard All Changes"},
}

var pullPushDropdownMenuItems = []DropdownMenuItem{
	{Label: "Sync"},
	{Separator: true},
	{Label: "Pull"},
	{Label: "Pull (Rebase)"},
	{Label: "Pull from..."},
	{Separator: true},
	{Label: "Push"},
	{Label: "Push to..."},
	{Separator: true},
	{Label: "Fetch"},
	{Label: "Fetch (Prune)"},
	{Label: "Fetch From All Remotes"},
}

var branchDropdownMenuItems = []DropdownMenuItem{
	{Label: "Merge..."},
	{Label: "Rebase Branch..."},
	{Separator: true},
	{Label: "Create Branch..."},
	{Label: "Create Branch From..."},
	{Separator: true},
	{Label: "Rename Branch..."},
	{Label: "Delete Branch..."},
	{Label: "Delete Remote Branch..."},
	{Separator: true},
	{Label: "Publish Branch..."},
}

var remoteDropdownMenuItems = []DropdownMenuItem{
	{Label: "Add Remote..."},
	{Label: "Remove Remote"},
}

var stashDropdownMenuItems = []DropdownMenuItem{
	{Label: "Stash"},
	{Label: "Stash (Include Untracked)"},
	{Label: "Stash Staged"},
	{Separator: true},
	{Label: "Apply Latest Stash"},
	{Label: "Apply Stash..."},
	{Separator: true},
	{Label: "Pop Latest Stash"},
	{Label: "Pop Stash..."},
	{Separator: true},
	{Label: "Drop Stash..."},
	{Label: "Drop All Stashes..."},
	{Separator: true},
	{Label: "View Stash..."},
}

var tagsDropdownMenuItems = []DropdownMenuItem{
	{Label: "Create Tag..."},
	{Label: "Delete Tag..."},
	{Label: "Delete Remote Tag..."},
	{Separator: true},
	{Label: "Push Tags"},
}

// submenuItemsByKind maps a submenu kind to its menu items.
var submenuItemsByKind = map[string][]DropdownMenuItem{
	"commit":    commitDropdownMenuItems,
	"changes":   changesDropdownMenuItems,
	"pull_push": pullPushDropdownMenuItems,
	"branch":    branchDropdownMenuItems,
	"remote":    remoteDropdownMenuItems,
	"stash":     stashDropdownMenuItems,
	"tags":      tagsDropdownMenuItems,
}

// submenuKindByLabel maps a main menu item label to its submenu kind.
var submenuKindByLabel = map[string]string{
	"Commit":     "commit",
	"Changes":    "changes",
	"Pull, Push": "pull_push",
	"Branch":     "branch",
	"Remote":     "remote",
	"Stash":      "stash",
	"Tags":       "tags",
}

func dropdownItemsMaxWidth(items []DropdownMenuItem) int {
	maxItemW := 0
	for _, item := range items {
		if item.Separator {
			continue
		}
		itemW := runewidth.StringWidth(item.Label) + 2
		if item.HasChevron {
			itemW += 2
		}
		if itemW > maxItemW {
			maxItemW = itemW
		}
	}
	return maxItemW
}
