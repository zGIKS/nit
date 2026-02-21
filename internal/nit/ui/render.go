package ui

import (
	"fmt"

	"nit/internal/nit/app"
)

func Render(state app.AppState) string {
	changesActive := state.Focus == app.FocusChanges
	graphActive := state.Focus == app.FocusGraph
	changeSel, changeTotal := state.ChangesPosition()
	graphSel, graphTotal := state.GraphPosition()

	changeLines := make([]string, 0, len(state.Changes.Rows))
	for _, r := range state.Changes.Rows {
		changeLines = append(changeLines, r.Text)
	}
	if len(changeLines) == 0 {
		changeLines = []string{"Working tree clean."}
	}

	changes := BoxView(
		"Changes",
		state.Viewport.Width,
		state.ChangesPaneHeight(),
		changeLines,
		state.Changes.Cursor,
		state.Changes.Offset,
		changesActive,
		fmt.Sprintf("%d of %d", changeSel, changeTotal),
	)
	graph := BoxView(
		"Commits - Reflog",
		state.Viewport.Width,
		state.GraphPaneHeight(),
		state.Graph.Lines,
		state.Graph.Cursor,
		state.Graph.Offset,
		graphActive,
		fmt.Sprintf("%d of %d", graphSel, graphTotal),
	)

	if state.LastErr != "" {
		err := BoxView(
			"Error",
			state.Viewport.Width,
			3,
			[]string{state.LastErr},
			0,
			0,
			true,
			"diagnostics",
		)
		return changes + "\n" + graph + "\n" + err
	}

	return changes + "\n" + graph
}
