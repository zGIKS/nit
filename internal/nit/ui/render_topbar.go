package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/zGIKS/nit/internal/nit/app"
)

func buildTopBar(state app.AppState, totalW int) string {
	repoName := state.RepoName
	if repoName == "" {
		repoName = "unknown"
	}
	branchName := state.BranchName
	if branchName == "" {
		branchName = "-"
	}
	sep := strings.TrimSpace(state.RepoBranchSeparator)
	if sep == "" {
		sep = "->"
	}
	repoText := strings.TrimSpace(
		strings.TrimSpace(state.RepoLabel+" "+repoName) +
			" " + sep + " " +
			strings.TrimSpace(state.BranchLabel+" "+branchName),
	)
	createText := strings.TrimSpace(state.BranchesCreateButtonLabel())
	fetchText := strings.TrimSpace(state.FetchLabel)
	menuText := strings.TrimSpace(state.MenuLabel)

	repoW := max(16, runewidth.StringWidth(repoText)+4)
	createW := max(12, runewidth.StringWidth(createText)+4)
	fetchW := max(8, runewidth.StringWidth(fetchText)+4)
	menuW := max(8, runewidth.StringWidth(menuText)+4)
	minRepoW := 14
	minCreateW := 12
	minFetchW := 8
	minMenuW := 8
	totalNeeded := repoW + createW + fetchW + menuW + 3
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
	shrink(&createW, minCreateW)
	shrink(&fetchW, minFetchW)
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		// Last resort: give remaining width to repo box and let text truncate.
		repoW = max(minRepoW, repoW-overflow)
	}

	leftTop := MiniBoxView(repoText, repoW)
	rightTopW := createW + fetchW + menuW + 2
	rightTop := HStackMany(
		[]string{
			func() string {
				if state.HoverBranch {
					return MiniBoxViewUnderline(createText, createW)
				}
				return MiniBoxView(createText, createW)
			}(),
			func() string {
				if state.HoverFetch {
					return MiniBoxViewUnderline(fetchText, fetchW)
				}
				return MiniBoxView(fetchText, fetchW)
			}(),
			func() string {
				if state.HoverMenu {
					return MiniBoxViewUnderline(menuText, menuW)
				}
				return MiniBoxView(menuText, menuW)
			}(),
		},
		[]int{createW, fetchW, menuW},
	)
	gapW := totalW - repoW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}
	spacerLine := strings.Repeat(" ", gapW)
	spacer := spacerLine + "\n" + spacerLine + "\n" + spacerLine
	return HStackMany(
		[]string{leftTop, spacer, rightTop},
		[]int{repoW, gapW, rightTopW},
	)
}
