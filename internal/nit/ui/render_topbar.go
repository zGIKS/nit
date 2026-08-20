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

	repoW := max(12, runewidth.StringWidth(repoText)+4)
	branchW := max(12, runewidth.StringWidth(branchText)+4)
	menuW := max(8, runewidth.StringWidth(menuText)+4)

	minRepoW := 10
	minBranchW := 10
	minMenuW := 8

	totalNeeded := repoW + branchW + 1 + menuW + 2
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
		repoW = max(minRepoW, repoW-overflow)
	}

	leftTop := HStack(MiniBoxView(repoText, repoW), repoW, MiniBoxView(branchText, branchW), branchW)
	leftTopW := repoW + branchW + 1

	rightTopW := menuW
	rightTop := MiniBoxView(menuText, menuW)
	if state.HoverMenu {
		rightTop = MiniBoxViewUnderline(menuText, menuW)
	}
	gapW := totalW - leftTopW - rightTopW - 2
	if gapW < 1 {
		gapW = 1
	}
	spacerLine := strings.Repeat(" ", gapW)
	spacer := spacerLine + "\n" + spacerLine + "\n" + spacerLine
	return HStackMany(
		[]string{leftTop, spacer, rightTop},
		[]int{leftTopW, gapW, rightTopW},
	)
}
