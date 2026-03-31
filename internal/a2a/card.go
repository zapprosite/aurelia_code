// Package a2a implements the Agent-to-Agent (A2A) protocol for inter-agent
// communication. Aurélia acts as an A2A agent that can be discovered and called
// by other agents on the network.
//
// Specification: https://github.com/a2aproject/A2A
package a2a

import (
	"time"
)

// AgentCard describes this agent to the A2A network.
// It is served at /.well-known/agent.json per the A2A spec.
type AgentCard struct {
	// Name of the agent.
	Name string `json:"name"`
	// Human-readable description of the agent's capabilities.
	Description string `json:"description"`
	// URL where the A2A server listens.
	URL string `json:"url"`
	// Semantic version of this agent.
	Version string `json:"version"`
	// Supported A2A protocol version(s).
	ProtocolVersion string `json:"protocolVersion"`
	// Capabilities this agent supports.
	Capabilities Capabilities `json:"capabilities"`
	// Skills this agent can perform.
	Skills []Skill `json:"skills"`
	// MIME types this agent accepts as input.
	InputModes []string `json:"inputModes"`
	// MIME types this agent produces as output.
	OutputModes []string `json:"outputModes"`
	// Authentication requirements for this agent.
	Authentication Authentication `json:"authentication,omitempty"`
	// Provider that hosts this agent.
	Provider *Provider `json:"provider,omitempty"`
}

// Capabilities lists the optional features this agent supports.
type Capabilities struct {
	// Streaming responses via SSE.
	Streaming bool `json:"streaming"`
	// Push notifications to external endpoints.
	PushNotifications bool `json:"pushNotifications"`
	// Long-running tasks with status tracking.
	TaskTracking bool `json:"taskTracking"`
	// Multiple agents can collaborate on a task.
	MultiAgent bool `json:"multiAgent"`
}

// Skill describes a discrete capability of this agent.
type Skill struct {
	// Unique identifier for this skill.
	ID string `json:"id"`
	// Human-readable name.
	Name string `json:"name"`
	// Description of what this skill does.
	Description string `json:"description"`
	 // MIME types this skill accepts.
	InputModes []string `json:"inputModes,omitempty"`
	// MIME types this skill produces.
	OutputModes []string `json:"outputModes,omitempty"`
	// Tags for discovery.
	Tags []string `json:"tags,omitempty"`
}

// Authentication describes how clients authenticate with this agent.
type Authentication struct {
	// Schemes supported: "none", "bearer", "api_key"
	Schemes []string `json:"schemes"`
	// Credentials endpoint for bearer token.
	CredentialsURL string `json:"credentialsUrl,omitempty"`
}

// Provider describes the infrastructure provider.
type Provider struct {
	Organization string `json:"organization,omitempty"`
	URL          string `json:"url,omitempty"`
}

// DefaultAgentCard returns the canonical AgentCard for Aurélia.
func DefaultAgentCard(baseURL string) *AgentCard {
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	return &AgentCard{
		Name:             "aurelia",
		Description:      "Assistente IA soberana com voz PT-BR, memória semântica e ferramentas MCP. Opera em português brasileiro.",
		URL:             baseURL,
		Version:          "1.0.0",
		ProtocolVersion:  "2025-03-01",
		Capabilities: Capabilities{
			Streaming:       true,
			TaskTracking:    true,
			PushNotifications: false,
			MultiAgent:      false,
		},
		Skills: []Skill{
			{
				ID:          "voice-synthesis",
				Name:        "Synthese de Voz",
				Description: "TTS em PT-BR com Kokoro/Kodoro — streaming de áudio",
				InputModes:  []string{"text/plain"},
				OutputModes: []string{"audio/mp3", "audio/opus"},
				Tags:        []string{"tts", "voice", "audio", "kokoro"},
			},
			{
				ID:          "mcp-tools",
				Name:        "Ferramentas MCP",
				Description: "Executa chamadas de ferramentas via Model Context Protocol",
				InputModes:  []string{"application/json"},
				OutputModes: []string{"application/json", "text/plain"},
				Tags:        []string{"mcp", "tools", "function-calling"},
			},
			{
				ID:          "semantic-memory",
				Name:        "Memoria Semantica",
				Description: "Busca e armazenamento de memória semântica via Qdrant",
				InputModes:  []string{"text/plain", "application/json"},
				OutputModes: []string{"application/json"},
				Tags:        []string{"memory", "qdrant", "vector", "rag"},
			},
			{
				ID:          "telegram-bot",
				Name:        "Bot Telegram",
				Description: "Interface de chat via Telegram com personas configuraveis",
				InputModes:  []string{"text/plain"},
				OutputModes: []string{"text/plain", "text/markdown", "audio/mp3"},
				Tags:        []string{"telegram", "bot", "chat"},
			},
		},
		InputModes:  []string{"text/plain", "application/json"},
		OutputModes: []string{"text/plain", "text/markdown", "application/json", "audio/mp3"},
		Authentication: Authentication{
			Schemes: []string{"bearer"},
		},
		Provider: &Provider{
			Organization: "Aurelia Sovereign 2026",
			URL:          "https://github.com/kocar/aurelia",
		},
	}
}

// DefaultAgentCardWithTimestamp is like DefaultAgentCard but includes a cache hint.
func DefaultAgentCardWithTimestamp(baseURL string) AgentCard {
	card := *DefaultAgentCard(baseURL)
	// A2A spec recommends exposing agent card at /.well-known/agent.json
	// with Cache-Control for caching. The handler sets this header.
	_ = time.Now // can be used by the server handler for ETag generation
	return card
}
