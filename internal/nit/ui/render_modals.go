package ui

import (
	"fmt"
	"strings"

	"github.com/zGIKS/nit/internal/nit/app"
)

func branchCreateModalView(state app.AppState, width, height int) string {
	w := max(36, width)
	if width > 0 {
		w = width
	}
	if w < 4 {
		w = 4
	}
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	lines := make([]string, 0, max(3, height))
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	bottom := "└" + strings.Repeat("─", innerW) + "┘"
	titleText := strings.TrimSpace(state.BranchCreateTitle)
	if titleText == "" {
		titleText = "Create a branch"
	}
	title := fitText(" "+titleText+" ", innerW, ' ')
	lines = append(lines, top)
	lines = append(lines, "│"+title+"│")
	lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
	enterHint := strings.TrimSpace(state.BranchCreateEnterHint)
	if enterHint == "" {
		enterHint = "Enter: create branch"
	}
	lines = append(lines, "│"+fitText(" "+enterHint, innerW, ' ')+"│")
	pushHint := strings.TrimSpace(state.BranchCreatePushHint)
	lines = append(lines, "│"+fitText(" "+pushHint, innerW, ' ')+"│")
	lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
	nameLabel := strings.TrimSpace(state.BranchCreateNameLabel)
	if nameLabel == "" {
		nameLabel = "New branch name"
	}
	lines = append(lines, "│"+fitText(" "+nameLabel, innerW, ' ')+"│")
	inputViewportW := max(1, innerW-1)
	lines = append(lines, "│"+fitText(" "+textInputViewport(state.BranchCreateName, state.BranchCreateCursor, state.BranchCreateSelectAll, inputViewportW), innerW, ' ')+"│")
	lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
	sourceLabel := strings.TrimSpace(state.BranchCreateSourceLabel)
	if sourceLabel == "" {
		sourceLabel = "Source"
	}
	lines = append(lines, "│"+fitText(" "+sourceLabel+": "+state.BranchCreateSource, innerW, ' ')+"│")

	_, _, _, remaining := state.BranchCreateSourceListRect()
	start := state.BranchCreateSourceOffset
	for i := 0; i < remaining; i++ {
		var row string
		idx := start + i
		if idx < len(state.BranchCreateSourceList) {
			name := state.BranchCreateSourceList[idx]
			prefix := "  "
			if name == state.BranchCreateSource {
				mark := state.BranchSourceSelectedMark
				if strings.TrimSpace(mark) == "" {
					mark = "✓"
				}
				prefix = mark + " "
			}
			label := prefix + name
			row = fitText(label, innerW, ' ')
		} else {
			row = fitText("", innerW, ' ')
		}
		lines = append(lines, "│"+row+"│")
	}
	lines = append(lines, bottom)

	if len(lines) > height {
		lines = lines[:height]
		if len(lines) > 0 {
			lines[len(lines)-1] = bottom
		}
	}
	for len(lines) < height {
		lines = append(lines, "│"+fitText("", innerW, ' ')+"│")
	}
	return strings.Join(lines, "\n")
}

func branchDeleteConfirmView(state app.AppState, width int) string {
	w := max(42, width)
	innerW := w - 2
	choice := "[Yes]    No"
	if state.BranchDeleteChoice == 1 {
		choice = "Yes    [No]"
	}
	lines := []string{
		"┌" + strings.Repeat("─", innerW) + "┐",
		"│" + fitText(" Delete branch?", innerW, ' ') + "│",
		"│" + fitText(fmt.Sprintf(" Are you sure you want to delete branch %q?", state.BranchDeleteBranch), innerW, ' ') + "│",
		"│" + fitText(strings.Repeat(" ", max(0, (innerW-displayWidth(choice))/2))+choice, innerW, ' ') + "│",
		"└" + strings.Repeat("─", innerW) + "┘",
	}
	return strings.Join(lines, "\n")
}
