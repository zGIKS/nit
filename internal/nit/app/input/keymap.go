package input

import (
	"fmt"

	"nit/internal/nit/app/actions"
	"nit/internal/nit/config"
)

type Keymap struct {
	bindings map[actions.Action][]string
}

func DefaultKeymap() Keymap {
	return Keymap{bindings: map[actions.Action][]string{
		actions.ActionQuit:         {"ctrl+c", "q"},
		actions.ActionTogglePanel:  {"tab"},
		actions.ActionFocusCommand: {"c"},
		actions.ActionMoveDown:     {"down", "j"},
		actions.ActionMoveUp:       {"up", "k"},
		actions.ActionToggleOne:    {"enter"},
		actions.ActionStageAll:     {"s"},
		actions.ActionUnstageAll:   {"u"},
		actions.ActionPush:         {"p", "ctrl+p"},
	}}
}

func LoadKeymap(cfg config.KeyConfig) (Keymap, string) {
	km := DefaultKeymap()
	merge := func(a actions.Action, b config.KeyBinding) {
		if len(b.Keys) > 0 {
			km.bindings[a] = b.Keys
		}
	}
	merge(actions.ActionQuit, cfg.Quit)
	merge(actions.ActionTogglePanel, cfg.TogglePanel)
	merge(actions.ActionFocusCommand, cfg.FocusCommand)
	merge(actions.ActionMoveDown, cfg.Down)
	merge(actions.ActionMoveUp, cfg.Up)
	merge(actions.ActionToggleOne, cfg.ToggleOne)
	merge(actions.ActionStageAll, cfg.StageAll)
	merge(actions.ActionUnstageAll, cfg.UnstageAll)
	merge(actions.ActionPush, cfg.Push)

	if err := validateKeyConflicts(km); err != nil {
		return DefaultKeymap(), "invalid key config: " + err.Error()
	}
	return km, ""
}

func validateKeyConflicts(km Keymap) error {
	seen := map[string]actions.Action{}
	for action, keys := range km.bindings {
		for _, k := range keys {
			if prev, ok := seen[k]; ok {
				return fmt.Errorf("duplicate key %q for actions %d and %d", k, prev, action)
			}
			seen[k] = action
		}
	}
	return nil
}

func (k Keymap) Match(key string) actions.Action {
	for action, keys := range k.bindings {
		for _, cand := range keys {
			if cand == key {
				return action
			}
		}
	}
	return actions.ActionNone
}
