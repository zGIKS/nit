package ui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"nit/internal/nit/app"
)

func Render(state app.AppState) string {
	commandActive := state.Focus == app.FocusCommand
	changesActive := state.Focus == app.FocusChanges
	graphActive := state.Focus == app.FocusGraph
	commandLogActive := state.Focus == app.FocusCommandLog
	changeSel, changeTotal := state.ChangesPosition()
	graphSel, graphTotal := state.GraphPosition()
	commandText := state.Command.Input
	if commandActive {
		commandText = commandLineViewport(state, max(1, commitContentWidth(state.Viewport.Width)))
	} else if commandText == "" {
		commandText = "Message (c focus, Enter commit)"
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

	repoName := state.RepoName
	if repoName == "" {
		repoName = "unknown"
	}
	branchName := state.BranchName
	if branchName == "" {
		branchName = "-"
	}
	repoText := strings.TrimSpace(state.RepoLabel + " " + repoName)
	branchText := strings.TrimSpace(state.BranchLabel + " " + branchName)
	fetchText := strings.TrimSpace(state.FetchLabel)
	menuText := strings.TrimSpace(state.MenuLabel)

	repoW := max(16, runewidth.StringWidth(repoText)+4)
	branchW := max(16, runewidth.StringWidth(branchText)+4)
	fetchW := max(14, runewidth.StringWidth(fetchText)+4)
	menuW := max(8, runewidth.StringWidth(menuText)+4)
	minRepoW := 14
	minBranchW := 12
	minFetchW := 10
	minMenuW := 8
	totalNeeded := repoW + branchW + fetchW + menuW + 3
	overflow := totalNeeded - totalW
	shrink := func(w *int, minW int) {
		if overflow <= 0 {
			return
		}
		can := *w - minW
		if can <= 0 {
			return
		}
		d := min(can, overflow)
		*w -= d
		overflow -= d
	}
	shrink(&repoW, minRepoW)
	shrink(&branchW, minBranchW)
	shrink(&fetchW, minFetchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		// Last resort: give remaining width to repo box and let text truncate.
		repoW = max(minRepoW, repoW-overflow)
	}

	leftTop := MiniBoxView(repoText, repoW)
	rightTopW := branchW + fetchW + menuW + 2
	rightTop := HStackMany(
		[]string{
			MiniBoxView(branchText, branchW),
			MiniBoxView(fetchText, fetchW),
			MiniBoxView(menuText, menuW),
		},
		[]int{branchW, fetchW, menuW},
	)
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}
	spacerLine := strings.Repeat(" ", gapW)
	spacer := spacerLine + "\n" + spacerLine + "\n" + spacerLine
	topBar := HStackMany(
		[]string{leftTop, spacer, rightTop},
		[]int{repoW, gapW, rightTopW},
	)

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
				return "Ctrl+P"
			}
			return "p"
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
	graph := BoxView(
		"Commits - Reflog",
		totalW,
		state.GraphPaneHeight(),
		state.Graph.Lines,
		state.Graph.Cursor,
		state.Graph.Offset,
		graphActive,
		fmt.Sprintf("%d of %d", graphSel, graphTotal),
	)
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

	return command + "\n" + changes + "\n" + graph + "\n" + commandLog
}

func commitContentWidth(totalWidth int) int {
	totalW := max(40, totalWidth)
	pushW := max(18, totalW/4)
	commitW := totalW - pushW - 1
	if commitW < 20 {
		commitW = 20
	}
	// BoxView visible width for content line is (w-4), but it also prepends
	// a 2-char cursor prefix ("▌ " or "  "), so the user text gets (w-6).
	return commitW - 6
}

func commandLineViewport(state app.AppState, width int) string {
	if width < 4 {
		return state.CommandLineWithCaret()
	}
	full := state.CommandLineWithCaret()
	r := []rune(full)
	if len(r) <= width {
		return full
	}

	caret := 0
	for i, ch := range r {
		if ch == '|' {
			caret = i
			break
		}
	}

	start := caret - width/2
	if start < 0 {
		start = 0
	}
	if start+width > len(r) {
		start = len(r) - width
	}
	if start < 0 {
		start = 0
	}
	end := min(len(r), start+width)
	return string(r[start:end])
}
