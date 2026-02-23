package model

import "github.com/zGIKS/nit/internal/nit/ui"

func (m Model) View() string {
	return ui.Render(m.State)
}
