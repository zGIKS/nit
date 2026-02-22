package model

import "nit/internal/nit/ui"

func (m Model) View() string {
	return ui.Render(m.State)
}
