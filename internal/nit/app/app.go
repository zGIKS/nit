package app

import (
	actionspkg "github.com/zGIKS/nit/internal/nit/app/actions"
	inputpkg "github.com/zGIKS/nit/internal/nit/app/input"
	statepkg "github.com/zGIKS/nit/internal/nit/app/state"
	"github.com/zGIKS/nit/internal/nit/config"
)

type (
	Action      = actionspkg.Action
	Operation   = actionspkg.Operation
	ApplyResult = actionspkg.ApplyResult
	OpKind      = actionspkg.OpKind

	AppState         = statepkg.AppState
	FocusState       = statepkg.FocusState
	Section          = statepkg.Section
	DropdownMenuItem = statepkg.DropdownMenuItem
	Keymap           = inputpkg.Keymap
)

const (
	ActionNone           = actionspkg.ActionNone
	ActionQuit           = actionspkg.ActionQuit
	ActionTogglePanel    = actionspkg.ActionTogglePanel
	ActionFocusCommand   = actionspkg.ActionFocusCommand
	ActionMoveUp         = actionspkg.ActionMoveUp
	ActionMoveDown       = actionspkg.ActionMoveDown
	ActionToggleOne      = actionspkg.ActionToggleOne
	ActionStageAll       = actionspkg.ActionStageAll
	ActionUnstageAll     = actionspkg.ActionUnstageAll
	ActionDiscardAll     = actionspkg.ActionDiscardAll
	ActionPull           = actionspkg.ActionPull
	ActionFetch          = actionspkg.ActionFetch
	ActionPush           = actionspkg.ActionPush
	ActionMenuRight      = actionspkg.ActionMenuRight
	ActionMenuLeft       = actionspkg.ActionMenuLeft
	ActionUndoLastCommit = actionspkg.ActionUndoLastCommit
	ActionAbortRebase    = actionspkg.ActionAbortRebase

	OpStagePath      = actionspkg.OpStagePath
	OpUnstagePath    = actionspkg.OpUnstagePath
	OpStageAll       = actionspkg.OpStageAll
	OpUnstageAll     = actionspkg.OpUnstageAll
	OpDiscardAll     = actionspkg.OpDiscardAll
	OpCommit         = actionspkg.OpCommit
	OpPull           = actionspkg.OpPull
	OpFetch          = actionspkg.OpFetch
	OpPush           = actionspkg.OpPush
	OpUndoLastCommit = actionspkg.OpUndoLastCommit
	OpAbortRebase    = actionspkg.OpAbortRebase

	FocusCommand    = statepkg.FocusCommand
	FocusChanges    = statepkg.FocusChanges
	FocusGraph      = statepkg.FocusGraph
	FocusBranches   = statepkg.FocusBranches
	FocusCommandLog = statepkg.FocusCommandLog
)

func New(keys Keymap) AppState {
	return statepkg.New(keys)
}

func LoadKeymap(cfg config.KeyConfig) (Keymap, string) {
	return inputpkg.LoadKeymap(cfg)
}
