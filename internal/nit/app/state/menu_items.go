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
	{Label: "Commit Staged"},
	{Label: "Commit All"},
	{Label: "Undo Last Commit"},
}

var changesDropdownMenuItems = []DropdownMenuItem{
	{Label: "Stage All Changes"},
	{Label: "Unstage All Changes"},
	{Label: "Discard All Changes"},
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
