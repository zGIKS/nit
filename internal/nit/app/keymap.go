package app

import (
	"encoding/json"
	"fmt"
	"os"
)

type Keymap struct {
	bindings map[Action][]string
}

type keyBinding struct {
	Keys []string `json:"keys"`
}

type keyConfig struct {
	Quit        keyBinding `json:"quit"`
	TogglePanel keyBinding `json:"toggle_panel"`
	Down        keyBinding `json:"down"`
	Up          keyBinding `json:"up"`
	ToggleOne   keyBinding `json:"toggle_one"`
	StageAll    keyBinding `json:"stage_all"`
	UnstageAll  keyBinding `json:"unstage_all"`
}

func DefaultKeymap() Keymap {
	return Keymap{bindings: map[Action][]string{
		ActionQuit:        {"ctrl+c", "q"},
		ActionTogglePanel: {"tab"},
		ActionMoveDown:    {"down", "j"},
		ActionMoveUp:      {"up", "k"},
		ActionToggleOne:   {"enter"},
		ActionStageAll:    {"s"},
		ActionUnstageAll:  {"u"},
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

	merge := func(a Action, b keyBinding) {
		if len(b.Keys) > 0 {
			km.bindings[a] = b.Keys
		}
	}
	merge(ActionQuit, cfg.Quit)
	merge(ActionTogglePanel, cfg.TogglePanel)
	merge(ActionMoveDown, cfg.Down)
	merge(ActionMoveUp, cfg.Up)
	merge(ActionToggleOne, cfg.ToggleOne)
	merge(ActionStageAll, cfg.StageAll)
	merge(ActionUnstageAll, cfg.UnstageAll)

	if err := validateKeyConflicts(km); err != nil {
		return DefaultKeymap(), "invalid key config: " + err.Error()
	}
	return km, ""
}

func validateKeyConflicts(km Keymap) error {
	seen := map[string]Action{}
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

func (k Keymap) Match(key string) Action {
	for action, keys := range k.bindings {
		for _, cand := range keys {
			if cand == key {
				return action
			}
		}
	}
	return ActionNone
}
