package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/observability"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var invalidToolNameChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// ── Types ─────────────────────────────────────────────────────────────────────

type ToolSpec struct {
	RegistryName string
	ServerName   string
	RemoteName   string
	Description  string
	Parameters   map[string]interface{}
}

type CallResult struct {
	Content string
	IsError bool
}

type Caller interface {
	CallTool(ctx context.Context, serverName, remoteToolName string, args map[string]interface{}) (*CallResult, error)
}

type serverSession struct {
	name        string
	session     *mcpsdk.ClientSession
	callTimeout time.Duration
}

type serverResult struct {
	name    string
	session *serverSession
	specs   []ToolSpec
	err     error
}

type closeResult struct {
	name string
	err  error
}

type Manager struct {
	mu      sync.RWMutex
	servers map[string]*serverSession
	tools   []ToolSpec
}

// ── Bootstrap ─────────────────────────────────────────────────────────────────

func NewManager(cfg config.MCPToolsConfig, workspace string) (*Manager, error) {
	if !cfg.Enabled || len(cfg.Servers) == 0 {
		return nil, nil
	}
	connectTimeout := timeoutFromMS(cfg.ConnectTimeoutMS, 10*time.Second)
	defaultCallTimeout := timeoutFromMS(cfg.CallTimeoutMS, 60*time.Second)
	enabledCfgs := enabledMCPServers(cfg.Servers)
	if len(enabledCfgs) == 0 {
		return nil, nil
	}
	manager := &Manager{servers: make(map[string]*serverSession)}
	manager.connectEnabledServers(cfg, enabledCfgs, workspace, connectTimeout, defaultCallTimeout)
	if len(manager.servers) == 0 {
		_ = manager.Close()
		return nil, fmt.Errorf("no MCP servers connected successfully")
	}
	return manager, nil
}

func (m *Manager) connectEnabledServers(
	cfg config.MCPToolsConfig,
	servers []config.MCPServerConfig,
	workspace string,
	connectTimeout, defaultCallTimeout time.Duration,
) {
	logger := observability.Logger("mcp.bootstrap")
	for result := range connectServers(cfg, servers, workspace, connectTimeout, defaultCallTimeout) {
		if result.err != nil {
			logger.Warn("failed to connect MCP server", slog.String("server", result.name), slog.Any("err", result.err))
			continue
		}
		m.servers[result.name] = result.session
		m.tools = append(m.tools, result.specs...)
		logger.Info("connected MCP server", slog.String("server", result.name), slog.Int("tool_count", len(result.specs)))
	}
}

func enabledMCPServers(servers map[string]config.MCPServerConfig) []config.MCPServerConfig {
	enabled := make([]config.MCPServerConfig, 0, len(servers))
	for _, serverCfg := range servers {
		if serverCfg.Enabled {
			enabled = append(enabled, serverCfg)
		}
	}
	return enabled
}

// ── Close ─────────────────────────────────────────────────────────────────────

func (m *Manager) Close() error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	var errs []string
	for result := range closeSessions(m.servers) {
		if result.err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", result.name, result.err))
		}
	}
	m.servers = map[string]*serverSession{}
	m.tools = nil
	if len(errs) > 0 {
		return fmt.Errorf("failed to close MCP sessions: %s", strings.Join(errs, "; "))
	}
	return nil
}

func closeSessions(servers map[string]*serverSession) <-chan closeResult {
	results := make(chan closeResult, len(servers))
	var wg sync.WaitGroup
	for name, server := range servers {
		if server == nil || server.session == nil {
			continue
		}
		wg.Add(1)
		go func(n string, sess *mcpsdk.ClientSession) {
			defer wg.Done()
			done := make(chan error, 1)
			go func() { done <- sess.Close() }()
			select {
			case err := <-done:
				results <- closeResult{name: n, err: err}
			case <-time.After(5 * time.Second):
				results <- closeResult{name: n, err: fmt.Errorf("close timed out after 5s")}
			}
		}(name, server.session)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

// ── Call ──────────────────────────────────────────────────────────────────────

func (m *Manager) ToolSpecs() []ToolSpec {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ToolSpec, len(m.tools))
	copy(out, m.tools)
	return out
}

func (m *Manager) CallTool(ctx context.Context, serverName, remoteToolName string, args map[string]interface{}) (*CallResult, error) {
	server, err := m.lookupServer(serverName)
	if err != nil {
		return nil, err
	}
	if args == nil {
		args = map[string]interface{}{}
	}
	callCtx, cancel := newCallContext(ctx, server.callTimeout)
	defer cancel()
	result, err := server.session.CallTool(callCtx, &mcpsdk.CallToolParams{
		Name:      remoteToolName,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("MCP server returned nil tool result")
	}
	return &CallResult{Content: formatCallToolResult(result), IsError: result.IsError}, nil
}

func (m *Manager) lookupServer(serverName string) (*serverSession, error) {
	if m == nil {
		return nil, fmt.Errorf("MCP manager is not initialized")
	}
	m.mu.RLock()
	server, ok := m.servers[serverName]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("MCP server %q is not connected", serverName)
	}
	return server, nil
}

func newCallContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

// ── Format ────────────────────────────────────────────────────────────────────

func formatCallToolResult(result *mcpsdk.CallToolResult) string {
	parts := make([]string, 0, len(result.Content)+1)
	for _, block := range result.Content {
		switch content := block.(type) {
		case *mcpsdk.TextContent:
			if strings.TrimSpace(content.Text) != "" {
				parts = append(parts, content.Text)
			}
		case *mcpsdk.ImageContent:
			parts = append(parts, fmt.Sprintf("[image content mime=%s bytes=%d]", content.MIMEType, len(content.Data)))
		case *mcpsdk.EmbeddedResource:
			if formatted := formatEmbeddedResource(content); formatted != "" {
				parts = append(parts, formatted)
			}
		default:
			if raw, err := json.Marshal(block); err == nil && len(raw) > 0 {
				parts = append(parts, string(raw))
			}
		}
	}
	if len(parts) == 0 {
		if result.IsError {
			return "MCP tool failed with empty response"
		}
		return "MCP tool completed with empty response"
	}
	return strings.Join(parts, "\n")
}

func formatEmbeddedResource(content *mcpsdk.EmbeddedResource) string {
	if content == nil || content.Resource == nil {
		return ""
	}
	switch {
	case strings.TrimSpace(content.Resource.Text) != "":
		return content.Resource.Text
	case len(content.Resource.Blob) > 0:
		return fmt.Sprintf("[embedded resource %s blob bytes=%d]", content.Resource.URI, len(content.Resource.Blob))
	default:
		return fmt.Sprintf("[embedded resource %s]", content.Resource.URI)
	}
}
