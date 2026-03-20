package config

import (
	"encoding/json"
	"os"
)

type MCPServerConfig struct {
	Name       string            `json:"name"`
	Enabled    bool              `json:"enabled"`
	Command    string            `json:"command"`
	Args       []string          `json:"args"`
	Env        map[string]string `json:"env"`
	WorkingDir string            `json:"workingDir"`
	Transport  string            `json:"transport"` // "stdio", "http"
	Endpoint   string            `json:"endpoint"`
	Headers    map[string]string `json:"headers"`
	TimeoutMS  int               `json:"timeoutMs"`
	AllowTools []string          `json:"allowTools"`
}

type MCPToolsConfig struct {
	Enabled          bool                       `json:"enabled"`
	ClientName       string                     `json:"clientName"`
	ClientVersion    string                     `json:"clientVersion"`
	ConnectTimeoutMS int                        `json:"connectTimeoutMs"`
	CallTimeoutMS    int                        `json:"callTimeoutMs"`
	Headers          map[string]string          `json:"headers"`
	Servers          map[string]MCPServerConfig `json:"mcpServers"`
}

// LoadMCPConfig reads an MCP servers JSON configuration file (Claude desktop style).
func LoadMCPConfig(path string) (*MCPToolsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Provide an empty but valid disabled config if file missing
			return &MCPToolsConfig{Enabled: false}, nil
		}
		return nil, err
	}

	type rawServerConfig struct {
		Name       string            `json:"name"`
		Enabled    *bool             `json:"enabled"`
		Command    string            `json:"command"`
		Args       []string          `json:"args"`
		Env        map[string]string `json:"env"`
		WorkingDir string            `json:"workingDir"`
		Transport  string            `json:"transport"`
		Endpoint   string            `json:"endpoint"`
		Headers    map[string]string `json:"headers"`
		TimeoutMS  int               `json:"timeoutMs"`
		AllowTools []string          `json:"allowTools"`
	}
	type rawToolsConfig struct {
		Enabled          *bool                      `json:"enabled"`
		ClientName       string                     `json:"clientName"`
		ClientVersion    string                     `json:"clientVersion"`
		ConnectTimeoutMS int                        `json:"connectTimeoutMs"`
		CallTimeoutMS    int                        `json:"callTimeoutMs"`
		Headers          map[string]string          `json:"headers"`
		Servers          map[string]rawServerConfig `json:"mcpServers"`
	}

	var raw rawToolsConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &MCPToolsConfig{
		ClientName:       raw.ClientName,
		ClientVersion:    raw.ClientVersion,
		ConnectTimeoutMS: raw.ConnectTimeoutMS,
		CallTimeoutMS:    raw.CallTimeoutMS,
		Headers:          raw.Headers,
		Servers:          make(map[string]MCPServerConfig, len(raw.Servers)),
	}
	if raw.Enabled != nil {
		cfg.Enabled = *raw.Enabled
	}

	enabledServerCount := 0
	for name, rawSrv := range raw.Servers {
		srv := MCPServerConfig{
			Name:       rawSrv.Name,
			Command:    rawSrv.Command,
			Args:       rawSrv.Args,
			Env:        rawSrv.Env,
			WorkingDir: rawSrv.WorkingDir,
			Transport:  rawSrv.Transport,
			Endpoint:   rawSrv.Endpoint,
			Headers:    rawSrv.Headers,
			TimeoutMS:  rawSrv.TimeoutMS,
			AllowTools: rawSrv.AllowTools,
		}
		if srv.Name == "" {
			srv.Name = name
		}
		if srv.Transport == "" && srv.Command != "" {
			srv.Transport = "stdio"
		}
		if rawSrv.Enabled != nil {
			srv.Enabled = *rawSrv.Enabled
		} else {
			srv.Enabled = true
		}
		if srv.Enabled {
			enabledServerCount++
		}
		cfg.Servers[name] = srv
	}
	if raw.Enabled == nil && enabledServerCount > 0 {
		cfg.Enabled = true
	}

	return cfg, nil
}
