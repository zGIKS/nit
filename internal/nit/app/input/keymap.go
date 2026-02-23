package input

import (
	"fmt"
	"strings"

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
		actions.ActionFetch:        {"f"},
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
	merge(actions.ActionFetch, cfg.Fetch)
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

func (k Keymap) FirstBinding(action actions.Action) string {
	keys := k.bindings[action]
	if len(keys) == 0 {
		return ""
	}
	return keys[0]
}

func (k Keymap) FirstBindingMatching(action actions.Action, match func(string) bool) string {
	keys := k.bindings[action]
	for _, key := range keys {
		if match == nil || match(key) {
			return key
		}
	}
	return ""
}

func (k Keymap) DisplayBinding(action actions.Action) string {
	return displayKey(k.FirstBinding(action))
}

func (k Keymap) DisplayBindingMatching(action actions.Action, match func(string) bool) string {
	key := k.FirstBindingMatching(action, match)
	if key == "" {
		return ""
	}
	return displayKey(key)
}

func displayKey(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	switch s {
	case "ctrl+c":
		return "Ctrl+C"
	case "ctrl+p":
		return "Ctrl+P"
	case "ctrl+b":
		return "Ctrl+B"
	case "ctrl+e":
		return "Ctrl+E"
	case "ctrl+a":
		return "Ctrl+A"
	case "tab":
		return "Tab"
	case "enter":
		return "Enter"
	case "space":
		return "Space"
	case "up":
		return "Up"
	case "down":
		return "Down"
	case "left":
		return "Left"
	case "right":
		return "Right"
	}
	if strings.HasPrefix(s, "ctrl+") && len(s) > len("ctrl+") {
		return "Ctrl+" + strings.ToUpper(s[len("ctrl+"):])
	}
	if len(s) == 1 {
		return strings.ToLower(s)
	}
	return s
}
