package nit

import (
	"encoding/json"
	"os"
)

func defaultKeyConfig() keyConfig {
	return keyConfig{
		Quit:            keyBinding{Keys: []string{"ctrl+c", "q"}},
		OpenMenu:        keyBinding{Keys: []string{"m"}},
		TogglePanel:     keyBinding{Keys: []string{"tab"}},
		ShowOutput:      keyBinding{Keys: []string{"o"}},
		Reload:          keyBinding{Keys: []string{"r"}},
		Down:            keyBinding{Keys: []string{"down", "j"}},
		Up:              keyBinding{Keys: []string{"up", "k"}},
		PageDown:        keyBinding{Keys: []string{"pgdown", "f"}},
		PageUp:          keyBinding{Keys: []string{"pgup", "b"}},
		Home:            keyBinding{Keys: []string{"home", "g"}},
		End:             keyBinding{Keys: []string{"end", "G"}},
		StageSelected:   keyBinding{Keys: []string{"+"}},
		UnstageSelected: keyBinding{Keys: []string{"-"}},
		StageAll:        keyBinding{Keys: []string{"S"}},
		UnstageAll:      keyBinding{Keys: []string{"U"}},
		MenuDown:        keyBinding{Keys: []string{"down", "j"}},
		MenuUp:          keyBinding{Keys: []string{"up", "k"}},
		MenuSelect:      keyBinding{Keys: []string{"enter"}},
		MenuClose:       keyBinding{Keys: []string{"esc", "q"}},
		PromptSubmit:    keyBinding{Keys: []string{"enter"}},
		PromptCancel:    keyBinding{Keys: []string{"esc"}},
		PromptBackspace: keyBinding{Keys: []string{"backspace", "delete"}},
	}
}

func loadKeyConfig() (keyConfig, string) {
	cfg := defaultKeyConfig()
	path := os.Getenv("NIT_KEYS_FILE")
	if path == "" {
		path = ".nit.keys.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, ""
	}

	var user keyConfig
	if err := json.Unmarshal(data, &user); err != nil {
		return cfg, "invalid key config: " + err.Error()
	}

	mergeBindings(&cfg.Quit, user.Quit)
	mergeBindings(&cfg.OpenMenu, user.OpenMenu)
	mergeBindings(&cfg.TogglePanel, user.TogglePanel)
	mergeBindings(&cfg.ShowOutput, user.ShowOutput)
	mergeBindings(&cfg.Reload, user.Reload)
	mergeBindings(&cfg.Down, user.Down)
	mergeBindings(&cfg.Up, user.Up)
	mergeBindings(&cfg.PageDown, user.PageDown)
	mergeBindings(&cfg.PageUp, user.PageUp)
	mergeBindings(&cfg.Home, user.Home)
	mergeBindings(&cfg.End, user.End)
	mergeBindings(&cfg.StageSelected, user.StageSelected)
	mergeBindings(&cfg.UnstageSelected, user.UnstageSelected)
	mergeBindings(&cfg.StageAll, user.StageAll)
	mergeBindings(&cfg.UnstageAll, user.UnstageAll)
	mergeBindings(&cfg.MenuDown, user.MenuDown)
	mergeBindings(&cfg.MenuUp, user.MenuUp)
	mergeBindings(&cfg.MenuSelect, user.MenuSelect)
	mergeBindings(&cfg.MenuClose, user.MenuClose)
	mergeBindings(&cfg.PromptSubmit, user.PromptSubmit)
	mergeBindings(&cfg.PromptCancel, user.PromptCancel)
	mergeBindings(&cfg.PromptBackspace, user.PromptBackspace)

	return cfg, ""
}

func mergeBindings(base *keyBinding, override keyBinding) {
	if len(override.Keys) > 0 {
		base.Keys = override.Keys
	}
}

func hasKey(binding keyBinding, key string) bool {
	for _, k := range binding.Keys {
		if k == key {
			return true
		}
	}
	return false
}

func primaryKey(binding keyBinding) string {
	if len(binding.Keys) == 0 {
		return "-"
	}
	return binding.Keys[0]
}
