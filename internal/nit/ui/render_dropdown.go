package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/zGIKS/nit/internal/nit/app"
)

type dropdownViewParams struct {
	items      []app.DropdownMenuItem
	offset     int
	hoverIndex int
	panelH     int
	chevron    string
	indicator  string
}

func dropdownView(p dropdownViewParams, width int) string {
	w := max(18, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	bottom := "└" + strings.Repeat("─", innerW) + "┘"
	page := max(1, p.panelH-2)
	start := p.offset
	if start < 0 {
		start = 0
	}
	if start > max(0, len(p.items)-page) {
		start = max(0, len(p.items)-page)
	}
	lines := make([]string, 0, page+2)
	lines = append(lines, top)
	for row := 0; row < page; row++ {
		i := start + row
		if i >= len(p.items) {
			lines = append(lines, "│"+fitText("", innerW, ' ')+"│")
			continue
		}
		item := p.items[i]
		if item.Separator {
			lines = append(lines, "├"+strings.Repeat("─", innerW)+"┤")
			continue
		}
		text := buildItemText(item, p.chevron, innerW, "  ")
		if p.hoverIndex == i {
			indicator := p.indicator
			if len(indicator) > 1 {
				indicator = indicator[:1]
			}
			text = buildItemText(item, p.chevron, innerW, indicator+" ")
		}
		lines = append(lines, "│"+text+"│")
	}
	lines = append(lines, bottom)
	return strings.Join(lines, "\n")
}

func buildItemText(item app.DropdownMenuItem, chevron string, innerW int, prefix string) string {
	if item.HasChevron && innerW >= 2 {
		target := max(0, innerW-2)
		pad := max(0, target-runewidth.StringWidth(prefix+item.Label))
		return fitText(prefix+item.Label+strings.Repeat(" ", pad)+" "+chevron, innerW, ' ')
	}
	return fitText(prefix+item.Label+" ", innerW, ' ')
}

func menuDropdownView(state app.AppState, width int) string {
	_, _, _, panelH := state.MenuPanelRect()
	return dropdownView(dropdownViewParams{
		items:      state.MenuItems(),
		offset:     state.MenuOffset,
		hoverIndex: state.MenuHoverIndex,
		panelH:     panelH,
		chevron:    state.MenuChevron,
		indicator:  state.MenuSelectionIndicator,
	}, width)
}

func menuSubmenuView(state app.AppState, width int) string {
	_, _, _, panelH := state.MenuSubmenuRect()
	return dropdownView(dropdownViewParams{
		items:      state.MenuSubmenuItems(),
		offset:     state.MenuSubOffset,
		hoverIndex: state.MenuSubHoverIndex,
		panelH:     panelH,
		chevron:    state.MenuChevron,
		indicator:  state.MenuSelectionIndicator,
	}, width)
}
