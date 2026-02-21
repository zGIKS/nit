package core

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	g "nit/internal/nit/git"
	"nit/internal/nit/ui"
)

type model struct {
	state                app.AppState
	git                  g.Service
	clipCfg              config.ClipboardConfig
	pasteHintAlreadySeen bool
}

type pollMsg struct{}

type changesLoadedMsg struct {
	entries []g.ChangeEntry
	err     error
}

type graphLoadedMsg struct {
	lines []string
	err   error
}

type opDoneMsg struct {
	err            error
	refreshChanges bool
	refreshGraph   bool
	command        string
}

const pollInterval = 700 * time.Millisecond

func schedulePoll() tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg { return pollMsg{} })
}

func loadChangesCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		entries, err := svc.LoadChanges()
		return changesLoadedMsg{entries: entries, err: err}
	}
}

func loadGraphCmd(svc g.Service) tea.Cmd {
	return func() tea.Msg {
		lines, err := svc.LoadGraph()
		return graphLoadedMsg{lines: lines, err: err}
	}
}

func execOpCmd(svc g.Service, op app.Operation, refreshChanges, refreshGraph bool) tea.Cmd {
	return func() tea.Msg {
		cmd, err := execOperation(svc, op)
		if err != nil {
			return opDoneMsg{err: err, command: cmd}
		}
		return opDoneMsg{refreshChanges: refreshChanges, refreshGraph: refreshGraph, command: cmd}
	}
}

func execOperation(svc g.Service, op app.Operation) (string, error) {
	switch op.Kind {
	case app.OpStagePath:
		return svc.StagePath(op.Path)
	case app.OpUnstagePath:
		return svc.UnstagePath(op.Path)
	case app.OpStageAll:
		return svc.StageAll()
	case app.OpUnstageAll:
		return svc.UnstageAll()
	case app.OpCommit:
		return svc.Commit(op.Message)
	case app.OpPush:
		return svc.Push()
	default:
		return "", nil
	}
}

func newModel() model {
	cfg, cfgWarn := config.Load()
	keys, keyErr := app.LoadKeymap(cfg.Keys)
	state := app.New(keys)
	state.SetGraph([]string{"Loading graph..."})
	state.SetChanges(nil)
	if keyErr != "" {
		state.SetError(keyErr)
	} else if cfgWarn != "" {
		state.SetError(cfgWarn)
	}

	runner := g.NewRunner(4 * time.Second)
	svc := g.NewService(runner)

	return model{
		state:   state,
		git:     svc,
		clipCfg: cfg.Clipboard,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(schedulePoll(), loadChangesCmd(m.git), loadGraphCmd(m.git))
}

func (m model) handleResult(result app.ApplyResult) tea.Cmd {
	if result.Quit {
		return tea.Quit
	}

	cmds := make([]tea.Cmd, 0, len(result.Operations)+2)
	if len(result.Operations) > 0 {
		for _, op := range result.Operations {
			cmds = append(cmds, execOpCmd(m.git, op, result.RefreshChanges, result.RefreshGraph))
		}
	} else {
		if result.RefreshChanges {
			cmds = append(cmds, loadChangesCmd(m.git))
		}
		if result.RefreshGraph {
			cmds = append(cmds, loadGraphCmd(m.git))
		}
	}
	if len(cmds) == 0 {
		return nil
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.state.Viewport.Width == msg.Width && m.state.Viewport.Height == msg.Height {
			return m, nil
		}
		m.state.SetViewport(msg.Width, msg.Height)
		return m, nil

	case pollMsg:
		return m, tea.Batch(schedulePoll(), loadChangesCmd(m.git))

	case changesLoadedMsg:
		if msg.err != nil {
			m.state.SetError(msg.err.Error())
		} else {
			m.state.SetError("")
			if !sameChanges(m.state.Changes.Entries, msg.entries) {
				m.state.SetChanges(msg.entries)
			}
		}
		m.state.Clamp()
		return m, nil

	case graphLoadedMsg:
		if msg.err != nil {
			m.state.SetError(msg.err.Error())
		} else {
			m.state.SetError("")
			m.state.SetGraph(msg.lines)
		}
		m.state.Clamp()
		return m, nil

	case opDoneMsg:
		if msg.command != "" {
			m.state.AddCommandLog(msg.command)
		}
		if msg.err != nil {
			m.state.SetError(msg.err.Error())
			m.state.Clamp()
			return m, nil
		}
		m.state.SetError("")
		cmds := make([]tea.Cmd, 0, 2)
		if msg.refreshChanges {
			cmds = append(cmds, loadChangesCmd(m.git))
		}
		if msg.refreshGraph {
			cmds = append(cmds, loadGraphCmd(m.git))
		}
		m.state.Clamp()
		if len(cmds) == 0 {
			return m, nil
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if m.state.Focus == app.FocusCommand {
			switch msg.Type {
			case tea.KeyEnter:
				result := m.state.Apply(app.ActionToggleOne)
				m.state.Clamp()
				return m, m.handleResult(result)
			case tea.KeyEsc:
				m.state.ExitCommandFocus()
				m.state.Clamp()
				return m, nil
			case tea.KeyCtrlC:
				selected := m.state.SelectedCommandText()
				if selected == "" {
					m.state.Clamp()
					return m, nil
				}
				m.state.SetCommandClipboard(selected)
				_ = copyWithMode(m.clipCfg, selected)
				m.state.SetError("")
				m.state.Clamp()
				return m, nil
			case tea.KeyBackspace:
				m.state.BackspaceCommandText()
				m.state.Clamp()
				return m, nil
			case tea.KeyDelete:
				m.state.DeleteCommandText()
				m.state.Clamp()
				return m, nil
			case tea.KeyLeft:
				m.state.MoveCommandCursorLeft()
				m.state.Clamp()
				return m, nil
			case tea.KeyRight:
				m.state.MoveCommandCursorRight()
				m.state.Clamp()
				return m, nil
			case tea.KeyHome:
				m.state.MoveCommandCursorToStart()
				m.state.Clamp()
				return m, nil
			case tea.KeyEnd, tea.KeyCtrlE:
				m.state.MoveCommandCursorToEnd()
				m.state.Clamp()
				return m, nil
			case tea.KeyCtrlA:
				m.state.SelectAllCommandText()
				m.state.Clamp()
				return m, nil
			case tea.KeyCtrlX:
				selected := m.state.SelectedCommandText()
				if selected == "" {
					m.state.Clamp()
					return m, nil
				}
				m.state.SetCommandClipboard(selected)
				_ = copyWithMode(m.clipCfg, selected)
				m.state.DeleteCommandSelection()
				m.state.SetError("")
				m.state.Clamp()
				return m, nil
			case tea.KeyCtrlV:
				pasted, err := pasteWithMode(m.clipCfg)
				if err != nil || pasted == "" {
					pasted = m.state.CommandClipboard()
				}
				if pasted == "" {
					if !m.pasteHintAlreadySeen && m.clipCfg.Mode == config.ClipboardOnlyCopy {
						m.state.SetError("paste from OS disabled in only_copy mode")
						m.pasteHintAlreadySeen = true
					}
					m.state.Clamp()
					return m, nil
				}
				m.state.AppendCommandText(pasted)
				m.state.SetError("")
				m.state.Clamp()
				return m, nil
			case tea.KeySpace:
				m.state.AppendCommandText(" ")
				m.state.Clamp()
				return m, nil
			case tea.KeyRunes:
				m.state.AppendCommandText(string(msg.Runes))
				m.state.Clamp()
				return m, nil
			}

			action := m.state.Keys.Match(msg.String())
			if action == app.ActionTogglePanel {
				result := m.state.Apply(action)
				m.state.Clamp()
				return m, m.handleResult(result)
			}
			// Ignore quit in commit input to avoid accidental exit while typing.
			m.state.Clamp()
			return m, nil
		}

		action := m.state.Keys.Match(msg.String())
		result := m.state.Apply(action)
		m.state.Clamp()
		return m, m.handleResult(result)
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

func (m model) View() string {
	return ui.Render(m.state)
}
