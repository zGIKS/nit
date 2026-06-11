package ui

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

type cell struct {
	char  rune
	style string
}

func parseAnsiLine(s string) []cell {
	var cells []cell
	var currentStyle strings.Builder
	for i := 0; i < len(s); {
		if end, ok := ansiSeqEnd(s, i); ok {
			seq := s[i:end]
			if seq == "\x1b[0m" {
				currentStyle.Reset()
			} else {
				currentStyle.WriteString(seq)
			}
			i = end
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		rw := runewidth.RuneWidth(r)
		if rw > 0 {
			cells = append(cells, cell{char: r, style: currentStyle.String()})
			for w := 1; w < rw; w++ {
				cells = append(cells, cell{char: 0, style: currentStyle.String()})
			}
		}
		i += size
	}
	return cells
}

func cellsToString(cells []cell) string {
	var sb strings.Builder
	activeStyle := ""
	for _, c := range cells {
		if c.char == 0 {
			continue
		}
		if c.style != activeStyle {
			if activeStyle != "" {
				sb.WriteString("\x1b[0m")
			}
			sb.WriteString(c.style)
			activeStyle = c.style
		}
		sb.WriteRune(c.char)
	}
	if activeStyle != "" {
		sb.WriteString("\x1b[0m")
	}
	return sb.String()
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
		baseCells := parseAnsiLine(baseLines[row])
		overCells := parseAnsiLine(ol)
		targetLen := x + width
		for len(baseCells) < targetLen {
			baseCells = append(baseCells, cell{char: ' ', style: ""})
		}
		for j := 0; j < width; j++ {
			if j < len(overCells) {
				baseCells[x+j] = overCells[j]
			} else {
				baseCells[x+j] = cell{char: ' ', style: ""}
			}
		}
		baseLines[row] = cellsToString(baseCells)
	}
	return strings.Join(baseLines, "\n")
}
