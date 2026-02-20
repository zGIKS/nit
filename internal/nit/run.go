package nit

import tea "github.com/charmbracelet/bubbletea"

func Run() error {
	keys := loadKeyConfig()
	m := initialModel(keys)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
