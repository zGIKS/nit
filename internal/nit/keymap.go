package nit

import (
	"encoding/json"
	"os"
)

func defaultKeyConfig() keyConfig {
	return keyConfig{
		Quit:        keyBinding{Keys: []string{"ctrl+c", "q"}},
		TogglePanel: keyBinding{Keys: []string{"tab"}},
		Down:        keyBinding{Keys: []string{"down", "j"}},
		Up:          keyBinding{Keys: []string{"up", "k"}},
		ToggleOne:   keyBinding{Keys: []string{"enter"}},
		StageAll:    keyBinding{Keys: []string{"s"}},
		UnstageAll:  keyBinding{Keys: []string{"u"}},
	}
}

func loadKeyConfig() keyConfig {
	cfg := defaultKeyConfig()
	path := os.Getenv("NIT_KEYS_FILE")
	if path == "" {
		path = ".nit.keys.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var user keyConfig
	if err := json.Unmarshal(data, &user); err != nil {
		return cfg
	}

	mergeBindings(&cfg.Quit, user.Quit)
	mergeBindings(&cfg.TogglePanel, user.TogglePanel)
	mergeBindings(&cfg.Down, user.Down)
	mergeBindings(&cfg.Up, user.Up)
	mergeBindings(&cfg.ToggleOne, user.ToggleOne)
	mergeBindings(&cfg.StageAll, user.StageAll)
	mergeBindings(&cfg.UnstageAll, user.UnstageAll)
	return cfg
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
