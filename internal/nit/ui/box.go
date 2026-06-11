package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func BoxView(title string, width, boxHeight int, lines []string, cursor, offset int, active bool, footer string) string {
	return boxViewWithTitles(title, "", width, boxHeight, lines, cursor, offset, active, footer)
}

func BoxViewTitleRight(title, titleRight string, width, boxHeight int, lines []string, cursor, offset int, active bool, footer string) string {
	return boxViewWithTitles(title, titleRight, width, boxHeight, lines, cursor, offset, active, footer)
}

func BoxViewPinnedTop(title string, width, boxHeight int, pinned []string, lines []string, cursor, offset int, active bool, footer string) string {
	w := max(8, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	var borderStyle lipgloss.Style
	if active {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.ActiveBorderColor))
	} else {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.InactiveBorderColor))
	}

	top := renderTopBorder(title, "", innerW, active)

	var b strings.Builder
	b.WriteString(top + "\n")

	pinnedRows := min(len(pinned), contentHeight)
	scrollRows := contentHeight - pinnedRows
	if scrollRows < 0 {
		scrollRows = 0
	}

	leftB := borderStyle.Render("│")
	rightB := borderStyle.Render("│")

	for i := 0; i < pinnedRows; i++ {
		text := fitText("  "+pinned[i], innerW-2, ' ')
		b.WriteString(leftB + " " + text + " " + rightB + "\n")
	}

	maxOffset := max(0, len(lines)-scrollRows)
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	end := min(len(lines), offset+scrollRows)
	for i := 0; i < scrollRows; i++ {
		idx := offset + i
		text := ""
		if idx < end {
			prefix := "  "
			if idx == cursor {
				prefix = "▌ "
				rawText := prefix + lines[idx]
				paddedText := fitText(rawText, innerW-2, ' ')
				text = CursorStyle.Render(paddedText)
			} else {
				rawText := prefix + lines[idx]
				text = fitText(rawText, innerW-2, ' ')
			}
		} else {
			text = strings.Repeat(" ", innerW-2)
		}
		b.WriteString(leftB + " " + text + " " + rightB + "\n")
	}

	bottom := renderBottomBorder(footer, innerW, active)
	b.WriteString(bottom)
	return b.String()
}

func boxViewWithTitles(title, titleRight string, width, boxHeight int, lines []string, cursor, offset int, active bool, footer string) string {
	w := max(8, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	contentHeight := boxHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	var borderStyle lipgloss.Style
	if active {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.ActiveBorderColor))
	} else {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.InactiveBorderColor))
	}

	top := renderTopBorder(title, titleRight, innerW, active)

	var b strings.Builder
	b.WriteString(top + "\n")

	maxOffset := max(0, len(lines)-contentHeight)
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	end := min(len(lines), offset+contentHeight)

	leftB := borderStyle.Render("│")
	rightB := borderStyle.Render("│")

	for i := 0; i < contentHeight; i++ {
		idx := offset + i
		text := ""
		if idx < end {
			prefix := "  "
			if idx == cursor {
				prefix = "▌ "
				rawText := prefix + lines[idx]
				paddedText := fitText(rawText, innerW-2, ' ')
				text = CursorStyle.Render(paddedText)
			} else {
				rawText := prefix + lines[idx]
				text = fitText(rawText, innerW-2, ' ')
			}
		} else {
			text = strings.Repeat(" ", innerW-2)
		}
		b.WriteString(leftB + " " + text + " " + rightB + "\n")
	}

	bottom := renderBottomBorder(footer, innerW, active)
	b.WriteString(bottom)
	return b.String()
}

func renderTopBorder(title, titleRight string, innerW int, active bool) string {
	var borderStyle lipgloss.Style
	var titleStyle lipgloss.Style
	if active {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.ActiveBorderColor))
		titleStyle = TitleActiveStyle
	} else {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.InactiveBorderColor))
		titleStyle = TitleInactiveStyle
	}

	head := title
	if active {
		head = "● " + head
	}
	headerText := titleStyle.Render(" " + head + " ")

	if r := strings.TrimSpace(titleRight); r != "" {
		rightText := TitleInactiveStyle.Render(" " + r + " ")
		spaces := innerW - displayWidth(headerText) - displayWidth(rightText)
		if spaces >= 1 {
			headerText = headerText + strings.Repeat(" ", spaces) + rightText
		}
	}

	textW := displayWidth(headerText)
	if textW > innerW {
		return borderStyle.Render("┌") + truncateDisplayWidth(headerText, innerW) + borderStyle.Render("┐")
	}

	filler := borderStyle.Render("─")
	leftLen := 1
	leftBorder := strings.Repeat(filler, leftLen)
	rightBorder := strings.Repeat(filler, innerW-textW-leftLen)

	return borderStyle.Render("┌") + leftBorder + headerText + rightBorder + borderStyle.Render("┐")
}

func renderBottomBorder(footer string, innerW int, active bool) string {
	var borderStyle lipgloss.Style
	if active {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.ActiveBorderColor))
	} else {
		borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.InactiveBorderColor))
	}

	if footer == "" {
		return borderStyle.Render("└" + strings.Repeat("─", innerW) + "┘")
	}

	footerText := " " + footer + " "
	textW := displayWidth(footerText)
	if textW > innerW {
		return borderStyle.Render("└") + truncateDisplayWidth(footerText, innerW) + borderStyle.Render("┘")
	}

	filler := borderStyle.Render("─")
	leftBorder := borderStyle.Render("─")
	rightBorder := strings.Repeat(filler, innerW-textW-1)
	return borderStyle.Render("└") + leftBorder + footerText + rightBorder + borderStyle.Render("┘")
}

