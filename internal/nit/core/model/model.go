package model

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"nit/internal/nit/app"
	"nit/internal/nit/config"
	"nit/internal/nit/core/model/cmds"
	g "nit/internal/nit/git"
)

type Model struct {
	State                app.AppState
	Git                  g.Service
	Watcher              *g.FSWatcher
	ClipCfg              config.ClipboardConfig
	PasteHintAlreadySeen bool
}

func New() Model {
	cfg, cfgWarn := config.Load()
	keys, keyErr := app.LoadKeymap(cfg.Keys)
	state := app.New(keys)
	state.SetTopBarLabels(cfg.UI.RepoLabel, cfg.UI.BranchLabel, cfg.UI.FetchLabel, cfg.UI.MenuLabel)
	state.SetUISymbols(cfg.UI.BranchSourceSelectedMark)
	state.SetUIText(
		cfg.UI.BranchCreateTitle,
		cfg.UI.BranchCreateEnterHint,
		cfg.UI.BranchCreatePushHint,
		cfg.UI.BranchCreateNameLabel,
		cfg.UI.BranchCreateSourceLabel,
	)
	state.SetGraph([]string{"Loading graph..."})
	state.SetBranches([]string{"Loading branches..."})
	state.SetChanges(nil)
	if keyErr != "" {
		state.SetError(keyErr)
	} else if cfgWarn != "" {
		state.SetError(cfgWarn)
	}

	runner := g.NewRunner(4 * time.Second)
	svc := g.NewService(runner)

	return Model{
		State:   state,
		Git:     svc,
		ClipCfg: cfg.Clipboard,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		cmds.ScheduleChangesPoll(),
		cmds.ScheduleGraphPoll(),
		cmds.LoadChangesCmd(m.Git),
		cmds.LoadGraphCmd(m.Git),
		cmds.LoadBranchesCmd(m.Git),
		cmds.LoadRepoSummaryCmd(m.Git),
		cmds.InitWatchCmd(m.Git),
	)
}
