package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMCPConfig_PreservesExplicitDisabledServers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp_servers.json")
	data := []byte(`{
		"enabled": true,
		"mcpServers": {
			"disabled-server": {
				"enabled": false,
				"command": "disabled-bin"
			},
			"implicit-enabled-server": {
				"command": "enabled-bin"
			}
		}
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadMCPConfig(path)
	if err != nil {
		t.Fatalf("LoadMCPConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadMCPConfig() returned nil config")
	}
	if !cfg.Enabled {
		t.Fatal("cfg.Enabled = false, want true")
	}
	if cfg.Servers["disabled-server"].Enabled {
		t.Fatal("disabled-server.Enabled = true, want false")
	}
	if !cfg.Servers["implicit-enabled-server"].Enabled {
		t.Fatal("implicit-enabled-server.Enabled = false, want true")
	}
}

func TestLoadMCPConfig_DoesNotEnableTopLevelWhenAllServersDisabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp_servers.json")
	data := []byte(`{
		"mcpServers": {
			"disabled-server": {
				"enabled": false,
				"command": "disabled-bin"
			}
		}
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadMCPConfig(path)
	if err != nil {
		t.Fatalf("LoadMCPConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadMCPConfig() returned nil config")
	}
	if cfg.Enabled {
		t.Fatal("cfg.Enabled = true, want false")
	}
}
