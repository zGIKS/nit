package input

import (
	"encoding/json"
	"fmt"
	"os"

	"nit/internal/nit/app/actions"
)

type Keymap struct {
	bindings map[actions.Action][]string
}

type keyBinding struct {
	Keys []string `json:"keys"`
}

type keyConfig struct {
	Quit         keyBinding `json:"quit"`
	TogglePanel  keyBinding `json:"toggle_panel"`
	FocusCommand keyBinding `json:"focus_command"`
	Down         keyBinding `json:"down"`
	Up           keyBinding `json:"up"`
	ToggleOne    keyBinding `json:"toggle_one"`
	StageAll     keyBinding `json:"stage_all"`
	UnstageAll   keyBinding `json:"unstage_all"`
	Push         keyBinding `json:"push"`
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
		actions.ActionPush:         {"p"},
	}}
}

func LoadKeymap() (Keymap, string) {
	km := DefaultKeymap()
	path := os.Getenv("NIT_KEYS_FILE")
	if path == "" {
		path = ".nit.keys.json"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return km, ""
	}

	var cfg keyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return km, "invalid key config: " + err.Error()
	}

	merge := func(a actions.Action, b keyBinding) {
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
