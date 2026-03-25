package config

import (
	"encoding/json"
	"os"

	"github.com/kocar/aurelia/internal/runtime"
)

// SaveBots persists the current bot list to the config file,
// preserving all other config fields unchanged.
func SaveBots(r *runtime.PathResolver, bots []BotConfig) error {
	cfg := defaultFileConfig(r)
	if data, err := os.ReadFile(r.AppConfig()); err == nil && len(data) != 0 {
		_ = json.Unmarshal(data, &cfg)
	}
	cfg.Bots = bots
	return writeConfigFile(r.AppConfig(), cfg)
}
