package nit

import (
	"fmt"
	"strings"
)

func defaultMenuItems() []action {
	return []action{
		{label: "View Graph + Changes", kind: "view_graph"},
		{label: "Show Git Output", kind: "show_output"},
		{label: "Stage Selected", kind: "stage_selected"},
		{label: "Unstage Selected", kind: "unstage_selected"},
		{label: "Stage All", kind: "stage_all"},
		{label: "Unstage All", kind: "unstage_all"},
		{label: "Pull", kind: "pull"},
		{label: "Push", kind: "push"},
		{label: "Clone", kind: "clone"},
		{label: "Checkout to...", kind: "checkout"},
		{label: "Fetch", kind: "fetch"},
		{label: "Commit", kind: "commit"},
		{label: "Changes (status)", kind: "changes"},
		{label: "Pull, Push", kind: "pull_push"},
		{label: "Branch", kind: "branch"},
		{label: "Remote", kind: "remote"},
		{label: "Stash", kind: "stash"},
		{label: "Tags", kind: "tags"},
	}
}

func initialModel(keys keyConfig) model {
	graph, graphErr := loadGraphLines()
	changes, changeLines, changesErr := loadChanges()

	errs := make([]string, 0, 2)
	if graphErr != nil {
		errs = append(errs, fmt.Sprintf("graph error: %v", graphErr))
	}
	if changesErr != nil {
		errs = append(errs, fmt.Sprintf("changes error: %v", changesErr))
	}

	m := model{
		ui:            uiBrowse,
		panel:         panelGraph,
		focus:         focusChanges,
		graphLines:    normalizeLines(graph),
		changeLines:   normalizeLines(changeLines),
		changeEntries: changes,
		outputLines:   []string{"Run an action from the menu to see command output."},
		height:        24,
		err:           strings.Join(errs, " | "),
		status:        "Ready",
		menuItems:     defaultMenuItems(),
		keys:          keys,
	}
	m.setActiveLines()
	m.clamp()
	return m
}
