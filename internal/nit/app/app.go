package app

import (
	actionspkg "nit/internal/nit/app/actions"
	inputpkg "nit/internal/nit/app/input"
	statepkg "nit/internal/nit/app/state"
	"nit/internal/nit/config"
)

type (
	Action      = actionspkg.Action
	Operation   = actionspkg.Operation
	ApplyResult = actionspkg.ApplyResult
	OpKind      = actionspkg.OpKind

	AppState   = statepkg.AppState
	FocusState = statepkg.FocusState
	Section    = statepkg.Section
	Keymap     = inputpkg.Keymap
)

const (
	ActionNone         = actionspkg.ActionNone
	ActionQuit         = actionspkg.ActionQuit
	ActionTogglePanel  = actionspkg.ActionTogglePanel
	ActionFocusCommand = actionspkg.ActionFocusCommand
	ActionMoveUp       = actionspkg.ActionMoveUp
	ActionMoveDown     = actionspkg.ActionMoveDown
	ActionToggleOne    = actionspkg.ActionToggleOne
	ActionStageAll     = actionspkg.ActionStageAll
	ActionUnstageAll   = actionspkg.ActionUnstageAll
	ActionPush         = actionspkg.ActionPush

	OpStagePath   = actionspkg.OpStagePath
	OpUnstagePath = actionspkg.OpUnstagePath
	OpStageAll    = actionspkg.OpStageAll
	OpUnstageAll  = actionspkg.OpUnstageAll
	OpCommit      = actionspkg.OpCommit
	OpPush        = actionspkg.OpPush

	FocusCommand = statepkg.FocusCommand
	FocusChanges = statepkg.FocusChanges
	FocusGraph   = statepkg.FocusGraph
)

func New(keys Keymap) AppState {
	return statepkg.New(keys)
}

func LoadKeymap(cfg config.KeyConfig) (Keymap, string) {
	return inputpkg.LoadKeymap(cfg)
}
