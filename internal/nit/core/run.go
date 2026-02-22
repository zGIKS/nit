package core

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/core/model"
)

func Run() error {
	opts := []tea.ProgramOption{tea.WithAltScreen()}
	switch strings.ToLower(strings.TrimSpace(os.Getenv("NIT_MOUSE_MODE"))) {
	case "", "cell":
		opts = append(opts, tea.WithMouseCellMotion())
	case "all":
		opts = append(opts, tea.WithMouseAllMotion())
	case "off":
		// Mouse disabled for terminals with incompatible mouse tracking support.
	default:
		opts = append(opts, tea.WithMouseCellMotion())
	}
	p := tea.NewProgram(model.New(), opts...)
	_, err := p.Run()
	return err
}
