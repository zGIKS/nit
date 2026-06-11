package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	ActiveBorderColor   string
	InactiveBorderColor string
	ActiveTitleBg       string
	ActiveTitleFg       string
	InactiveTitleFg     string
	CursorBg            string
	CursorFg            string
	SuccessFg           string
	ErrorFg             string
}

var DefaultTheme = Theme{
	ActiveBorderColor:   "#7D56F4", // Indigo
	InactiveBorderColor: "#444444", // Dark Gray
	ActiveTitleBg:       "#7D56F4",
	ActiveTitleFg:       "#FFFFFF",
	InactiveTitleFg:     "#888888",
	CursorBg:            "#383838",
	CursorFg:            "#FFFFFF",
	SuccessFg:           "#04B575",
	ErrorFg:             "#FF5555",
}

// Border Styles
var (
	ActiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(DefaultTheme.ActiveBorderColor))

	InactiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(DefaultTheme.InactiveBorderColor))

	TitleActiveStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(DefaultTheme.ActiveTitleBg)).
				Foreground(lipgloss.Color(DefaultTheme.ActiveTitleFg)).
				Bold(true).
				Padding(0, 1)

	TitleInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DefaultTheme.InactiveTitleFg)).
				Bold(true).
				Padding(0, 1)

	CursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(DefaultTheme.CursorBg)).
			Foreground(lipgloss.Color(DefaultTheme.CursorFg))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DefaultTheme.SuccessFg)).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DefaultTheme.ErrorFg)).
			Bold(true)
)
