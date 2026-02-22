package core

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/core/model"
)

func Run() error {
	p := tea.NewProgram(model.New(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
