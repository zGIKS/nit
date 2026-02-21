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
	ActionPush
)

type OpKind int

const (
	OpStagePath OpKind = iota
	OpUnstagePath
	OpStageAll
	OpUnstageAll
	OpCommit
	OpPush
)

type Operation struct {
	Kind    OpKind
	Path    string
	Message string
}

type ApplyResult struct {
	Quit           bool
	Operations     []Operation
	RefreshChanges bool
}
