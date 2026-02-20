package nit

import tea "github.com/charmbracelet/bubbletea"

func Run() error {
	keys, keyErr := loadKeyConfig()
	m := initialModel(keys)
	if keyErr != "" {
		if m.err == "" {
			m.err = keyErr
		} else {
			m.err = m.err + " | " + keyErr
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
