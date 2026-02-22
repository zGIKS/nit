package common

import g "nit/internal/nit/git"

type PollMsg struct{}
type GraphPollMsg struct{}
type WatchTickMsg struct{}

type WatchReadyMsg struct {
	Watcher *g.FSWatcher
	Err     error
}

type ChangesLoadedMsg struct {
	Entries []g.ChangeEntry
	Err     error
}

type GraphLoadedMsg struct {
	Lines []string
	Err   error
}

type RepoSummaryLoadedMsg struct {
	Repo   string
	Branch string
	Err    error
}

type OpDoneMsg struct {
	Err            error
	RefreshChanges bool
	RefreshGraph   bool
	Command        string
}
