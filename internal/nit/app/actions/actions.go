package actions

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionTogglePanel
	ActionFocusCommand
	ActionMoveUp
	ActionMoveDown
	ActionToggleOne
	ActionStageAll
	ActionUnstageAll
	ActionDiscardAll
	ActionPull
	ActionFetch
	ActionPush
	ActionMenuRight
	ActionMenuLeft
	ActionUndoLastCommit
	ActionAbortRebase
)

type OpKind int

const (
	OpStagePath OpKind = iota
	OpUnstagePath
	OpStageAll
	OpUnstageAll
	OpDiscardAll
	OpCommit
	OpPull
	OpFetch
	OpPush
	OpUndoLastCommit
	OpAbortRebase
)

type Operation struct {
	Kind          OpKind
	Path          string
	Message       string
	CommitAll     bool
	CommitAmend   bool
	CommitSignoff bool
}

type ApplyResult struct {
	Quit           bool
	Operations     []Operation
	RefreshChanges bool
	RefreshGraph   bool
}
