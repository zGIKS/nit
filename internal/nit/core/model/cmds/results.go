package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	g "nit/internal/nit/git"
)

func HandleResult(git g.Service, result app.ApplyResult) tea.Cmd {
	if result.Quit {
		return tea.Quit
	}

	cmds := make([]tea.Cmd, 0, len(result.Operations)+2)
	if len(result.Operations) > 0 {
		for _, op := range result.Operations {
			cmds = append(cmds, ExecOpCmd(git, op, result.RefreshChanges, result.RefreshGraph))
		}
	} else {
		if result.RefreshChanges {
			cmds = append(cmds, LoadChangesCmd(git))
		}
		if result.RefreshGraph {
			cmds = append(cmds, LoadGraphCmd(git))
		}
	}
	if len(cmds) == 0 {
		return nil
	}
	return tea.Batch(cmds...)
}
