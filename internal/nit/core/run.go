package core

import tea "github.com/charmbracelet/bubbletea"

func Run() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
