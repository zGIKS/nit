package core

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	g "nit/internal/nit/git"
	"nit/internal/nit/ui"
)

type model struct {
	state app.AppState
	git   g.Service
}

type pollMsg struct{}

const pollInterval = 700 * time.Millisecond

func schedulePoll() tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg { return pollMsg{} })
}

func newModel() model {
	keys, keyErr := app.LoadKeymap()
	state := app.New(keys)
	if keyErr != "" {
		state.SetError(keyErr)
	}

	runner := g.NewRunner(4 * time.Second)
	svc := g.NewService(runner)
	graph, err := svc.LoadGraph()
	if err != nil {
		state.SetError(err.Error())
	}
	state.SetGraph(graph)

	changes, err := svc.LoadChanges()
	if err != nil {
		state.SetError(err.Error())
	}
	state.SetChanges(changes)
	state.Clamp()

	return model{state: state, git: svc}
}

func (m model) Init() tea.Cmd { return schedulePoll() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.state.Viewport.Width == msg.Width && m.state.Viewport.Height == msg.Height {
			return m, nil
		}
		m.state.SetViewport(msg.Width, msg.Height)
		return m, nil
	case pollMsg:
		changes, err := m.git.LoadChanges()
		if err != nil {
			m.state.SetError(err.Error())
		} else {
			m.state.SetError("")
			if !sameChanges(m.state.Changes.Entries, changes) {
				m.state.SetChanges(changes)
			}
		}
		m.state.Clamp()
		return m, schedulePoll()
	case tea.KeyMsg:
		action := m.state.Keys.Match(msg.String())
		result := m.state.Apply(action)
		if result.Quit {
			return m, tea.Quit
		}
		if len(result.Operations) > 0 {
			for _, op := range result.Operations {
				if err := m.execOp(op); err != nil {
					m.state.SetError(err.Error())
					break
				}
				m.state.SetError("")
			}
		}
		if result.RefreshChanges {
			changes, err := m.git.LoadChanges()
			if err != nil {
				m.state.SetError(err.Error())
			} else {
				m.state.SetError("")
				m.state.SetChanges(changes)
			}
		}
		m.state.Clamp()
		return m, nil
	}
	return m, nil
}

func sameChanges(a, b []g.ChangeEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Raw != b[i].Raw {
			return false
		}
	}
	return true
}

func (m *model) execOp(op app.Operation) error {
	switch op.Kind {
	case app.OpStagePath:
		return m.git.StagePath(op.Path)
	case app.OpUnstagePath:
		return m.git.UnstagePath(op.Path)
	case app.OpStageAll:
		return m.git.StageAll()
	case app.OpUnstageAll:
		return m.git.UnstageAll()
	default:
		return nil
	}
}

func (m model) View() string {
	return ui.Render(m.state)
}
