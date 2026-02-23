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
	repoText := strings.TrimSpace(state.RepoLabel + " " + repoName)
	branchText := strings.TrimSpace(state.BranchLabel + " " + branchName)
	menuText := strings.TrimSpace(state.MenuLabel)

	repoW := max(16, runewidth.StringWidth(repoText)+4)
	branchW := max(16, runewidth.StringWidth(branchText)+4)
	menuW := max(8, runewidth.StringWidth(menuText)+4)
	minRepoW := 14
	minBranchW := 12
	minMenuW := 8
	totalNeeded := repoW + branchW + menuW + 2
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
	shrink(&menuW, minMenuW)
	if overflow > 0 {
		// Last resort: give remaining width to repo box and let text truncate.
		repoW = max(minRepoW, repoW-overflow)
	}

	leftTop := MiniBoxView(repoText, repoW)
	rightTopW := branchW + menuW + 1
	rightTop := HStackMany(
		[]string{
			func() string {
				if state.HoverBranch {
					return MiniBoxViewUnderline(branchText, branchW)
				}
				return MiniBoxView(branchText, branchW)
			}(),
			func() string {
				if state.HoverMenu {
					return MiniBoxViewUnderline(menuText, menuW)
				}
				return MiniBoxView(menuText, menuW)
			}(),
		},
		[]int{branchW, menuW},
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
