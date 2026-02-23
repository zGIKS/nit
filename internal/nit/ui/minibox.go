package ui

import "strings"

func MiniBoxView(text string, width int) string {
	return miniBoxViewStyled(text, width, false)
}

func MiniBoxViewUnderline(text string, width int) string {
	return miniBoxViewStyled(text, width, true)
}

func miniBoxViewStyled(text string, width int, underline bool) string {
	w := max(8, width)
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	top := "┌" + strings.Repeat("─", innerW) + "┐"
	midText := fitText(" "+strings.TrimSpace(text)+" ", innerW, ' ')
	if underline {
		midText = ansiUnderline(midText)
	}
	mid := "│" + midText + "│"
	bot := "└" + strings.Repeat("─", innerW) + "┘"
	return top + "\n" + mid + "\n" + bot
}

func TopBarView(width int, left, right string) string {
	w := max(8, width)
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	leftLen := displayWidth(left)
	rightLen := displayWidth(right)

	if rightLen >= w {
		return fitText(right, w, ' ')
	}

	space := w - leftLen - rightLen
	if space < 1 {
		maxLeft := w - rightLen - 1
		if maxLeft < 0 {
			maxLeft = 0
		}
		left = fitText(left, maxLeft, ' ')
		left = strings.TrimRight(left, " ")
		space = w - displayWidth(left) - rightLen
		if space < 1 {
			space = 1
		}
	}

	return left + strings.Repeat(" ", space) + right
}
