package ui

import (
	"fmt"
	"strings"

	"github.com/zGIKS/nit/internal/nit/app"
)

func Render(state app.AppState) string {
	commandActive := state.Focus == app.FocusCommand
	changesActive := state.Focus == app.FocusChanges
	graphActive := state.Focus == app.FocusGraph
	branchesActive := state.Focus == app.FocusBranches
	commandLogActive := state.Focus == app.FocusCommandLog
	changeSel, changeTotal := state.ChangesPosition()
	graphSel, graphTotal := state.GraphPosition()
	branchSel, branchTotal := state.BranchesPosition()
	commandText := state.Command.Input
	focusCommandKey := state.Keys.DisplayBinding(app.ActionFocusCommand)
	if focusCommandKey == "" {
		focusCommandKey = "c"
	}
	pushKeyNormal := state.Keys.DisplayBindingMatching(app.ActionPush, func(k string) bool { return !strings.HasPrefix(k, "ctrl+") })
	if pushKeyNormal == "" {
		pushKeyNormal = state.Keys.DisplayBinding(app.ActionPush)
	}
	if pushKeyNormal == "" {
		pushKeyNormal = "p"
	}
	pushKeyInCommand := state.Keys.DisplayBindingMatching(app.ActionPush, func(k string) bool { return strings.HasPrefix(k, "ctrl+") })
	if pushKeyInCommand == "" {
		pushKeyInCommand = pushKeyNormal
	}
	if commandActive {
		commandText = commandLineViewport(state, max(1, commitContentWidth(state.Viewport.Width)))
	} else if commandText == "" {
		commandText = fmt.Sprintf("Message (%s focus, Enter commit)", focusCommandKey)
	}

	changeLines := make([]string, 0, len(state.Changes.Rows))
	for _, r := range state.Changes.Rows {
		changeLines = append(changeLines, r.Text)
	}
	if len(changeLines) == 0 {
		changeLines = []string{"Working tree clean."}
	}

	totalW := max(40, state.Viewport.Width)
	pushW := max(18, totalW/4)
	commitW := totalW - pushW - 1
	if commitW < 20 {
		commitW = 20
		pushW = max(8, totalW-commitW-1)
	}

	topBar := buildTopBar(state, totalW)
	commandBox := BoxView(
		"Commit",
		commitW,
		3,
		[]string{commandText},
		0,
		0,
		commandActive,
		"",
	)
	pushBox := BoxView(
		"Push",
		pushW,
		3,
		[]string{func() string {
			if commandActive {
				return pushKeyInCommand
			}
			return pushKeyNormal
		}()},
		0,
		0,
		false,
		"",
	)
	commandRow := HStack(commandBox, commitW, pushBox, pushW)
	command := topBar + "\n" + commandRow
	changes := BoxView(
		"Changes",
		totalW,
		state.ChangesPaneHeight(),
		changeLines,
		state.Changes.Cursor,
		state.Changes.Offset,
		changesActive,
		fmt.Sprintf("%d of %d", changeSel, changeTotal),
	)
	graphPaneW, branchPaneW := state.GraphBranchesPaneWidths()
	graphBox := BoxView(
		"Commits - Reflog",
		graphPaneW,
		state.GraphPaneHeight(),
		state.Graph.Lines,
		state.Graph.Cursor,
		state.Graph.Offset,
		graphActive,
		fmt.Sprintf("%d of %d", graphSel, graphTotal),
	)
	branchesBox := BoxView(
		"Branches",
		branchPaneW,
		state.GraphPaneHeight(),
		state.Branches.Lines,
		state.Branches.Cursor,
		state.Branches.Offset,
		branchesActive,
		fmt.Sprintf("%d of %d", branchSel, branchTotal),
	)
	graph := HStack(graphBox, graphPaneW, branchesBox, branchPaneW)
	commandLogFooter := ""
	if state.LastErr != "" {
		commandLogFooter = "error: " + state.LastErr
	}
	commandLog := BoxView(
		"Command Log",
		totalW,
		state.CommandLogPaneHeight(),
		state.CommandLog,
		func() int {
			if commandLogActive {
				return state.CommandLogView.Cursor
			}
			return len(state.CommandLog) - 1
		}(),
		func() int {
			if commandLogActive {
				return state.CommandLogView.Offset
			}
			return max(0, len(state.CommandLog)-(state.CommandLogPaneHeight()-2))
		}(),
		commandLogActive,
		commandLogFooter,
	)

	out := command + "\n" + changes + "\n" + graph + "\n" + commandLog
	if state.MenuOpen {
		menuPanelX, menuPanelY, menuPanelW, _ := state.MenuPanelRect()
		out = overlayBlock(out, menuDropdownView(state, menuPanelW), menuPanelX, menuPanelY, menuPanelW)
	}
	if state.BranchCreateOpen {
		panelX, panelY, panelW, panelH := state.BranchCreatePanelRect()
		out = overlayBlock(out, branchCreateModalView(state, panelW, panelH), panelX, panelY, panelW)
	}
	return out
}
