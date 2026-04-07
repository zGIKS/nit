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
}

var commitDropdownMenuItems = []DropdownMenuItem{
	{Label: "Commit"},
	{Label: "Commit Staged"},
	{Label: "Commit All"},
	{Label: "Undo Last Commit"},
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

// submenuItemsByKind maps a submenu kind to its menu items.
var submenuItemsByKind = map[string][]DropdownMenuItem{
	"commit":  commitDropdownMenuItems,
	"changes": changesDropdownMenuItems,
}

// submenuKindByLabel maps a main menu item label to its submenu kind.
var submenuKindByLabel = map[string]string{
	"Commit":  "commit",
	"Changes": "changes",
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
