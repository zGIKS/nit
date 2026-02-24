package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/zGIKS/nit/internal/nit/app"
)

func menuDropdownView(state app.AppState, width int) string {
	items := state.MenuItems()
	w := max(18, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	bottom := "└" + strings.Repeat("─", innerW) + "┘"
	_, _, _, panelH := state.MenuPanelRect()
	page := max(1, panelH-2)
	start := state.MenuOffset
	if start < 0 {
		start = 0
	}
	if start > max(0, len(items)-page) {
		start = max(0, len(items)-page)
	}
	lines := make([]string, 0, page+2)
	lines = append(lines, top)
	for row := 0; row < page; row++ {
		i := start + row
		if i >= len(items) {
			lines = append(lines, "│"+fitText("", innerW, ' ')+"│")
			continue
		}
		item := items[i]
		if item.Separator {
			lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
			continue
		}
		text := fitText("  "+item.Label+" ", innerW, ' ')
		if item.HasChevron && innerW >= 2 {
			target := max(0, innerW-2)
			pad := max(0, target-runewidth.StringWidth("  "+item.Label))
			text = fitText("  "+item.Label+strings.Repeat(" ", pad)+" "+state.MenuChevron, innerW, ' ')
		}
		if state.MenuHoverIndex == i {
			indicator := state.MenuSelectionIndicator
			if len(indicator) > 1 {
				indicator = indicator[:1]
			}
			text = fitText(indicator+" "+item.Label+" ", innerW, ' ')
			if item.HasChevron && innerW >= 2 {
				target := max(0, innerW-2)
				pad := max(0, target-runewidth.StringWidth(indicator+" "+item.Label))
				text = fitText(indicator+" "+item.Label+strings.Repeat(" ", pad)+" "+state.MenuChevron, innerW, ' ')
			}
		}
		lines = append(lines, "│"+text+"│")
	}
	lines = append(lines, bottom)
	return strings.Join(lines, "\n")
}

func menuSubmenuView(state app.AppState, width int) string {
	items := state.MenuSubmenuItems()
	w := max(18, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	bottom := "└" + strings.Repeat("─", innerW) + "┘"
	_, _, _, panelH := state.MenuSubmenuRect()
	page := max(1, panelH-2)
	start := state.MenuSubOffset
	if start < 0 {
		start = 0
	}
	if start > max(0, len(items)-page) {
		start = max(0, len(items)-page)
	}
	lines := make([]string, 0, page+2)
	lines = append(lines, top)
	for row := 0; row < page; row++ {
		i := start + row
		if i >= len(items) {
			lines = append(lines, "│"+fitText("", innerW, ' ')+"│")
			continue
		}
		item := items[i]
		if item.Separator {
			lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
			continue
		}
		text := fitText("  "+item.Label+" ", innerW, ' ')
		if item.HasChevron && innerW >= 2 {
			target := max(0, innerW-2)
			pad := max(0, target-runewidth.StringWidth("  "+item.Label))
			text = fitText("  "+item.Label+strings.Repeat(" ", pad)+" "+state.MenuChevron, innerW, ' ')
		}
		if state.MenuSubHoverIndex == i {
			indicator := state.MenuSelectionIndicator
			if len(indicator) > 1 {
				indicator = indicator[:1]
			}
			text = fitText(indicator+" "+item.Label+" ", innerW, ' ')
			if item.HasChevron && innerW >= 2 {
				target := max(0, innerW-2)
				pad := max(0, target-runewidth.StringWidth(indicator+" "+item.Label))
				text = fitText(indicator+" "+item.Label+strings.Repeat(" ", pad)+" "+state.MenuChevron, innerW, ' ')
			}
		}
		lines = append(lines, "│"+text+"│")
	}
	lines = append(lines, bottom)
	return strings.Join(lines, "\n")
}

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
	if pushHint == "" {
		pushHint = "Ctrl+b: create and push to origin"
	}
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

func overlayBlock(base, overlay string, x, y, width int) string {
	if base == "" || overlay == "" || x < 0 || y < 0 || width <= 0 {
		return base
	}
	baseLines := strings.Split(base, "\n")
	overLines := strings.Split(overlay, "\n")
	for i, ol := range overLines {
		row := y + i
		if row < 0 || row >= len(baseLines) {
			continue
		}
		bl := []rune(baseLines[row])
		if len(bl) < x {
			bl = append(bl, []rune(strings.Repeat(" ", x-len(bl)))...)
		}
		end := x + width
		if len(bl) < end {
			bl = append(bl, []rune(strings.Repeat(" ", end-len(bl)))...)
		}
		left := string(bl[:x])
		right := string(bl[end:])
		baseLines[row] = left + ol + right
	}
	return strings.Join(baseLines, "\n")
}
