package app

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionTogglePanel
	ActionMoveUp
	ActionMoveDown
	ActionToggleOne
	ActionStageAll
	ActionUnstageAll
)

type OpKind int

const (
	OpStagePath OpKind = iota
	OpUnstagePath
	OpStageAll
	OpUnstageAll
)

type Operation struct {
	Kind OpKind
	Path string
}

type ApplyResult struct {
	Quit           bool
	Operations     []Operation
	RefreshChanges bool
}
