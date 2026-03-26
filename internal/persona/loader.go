package persona

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
	"gopkg.in/yaml.v3"
)

// Config holds the persona frontmatter configuration.
type Config struct {
	Name             string   `yaml:"name"`
	Role             string   `yaml:"role"`
	MemoryWindowSize int      `yaml:"memory_window_size"`
	Tools            []string `yaml:"tools"`
}

// Persona is the parsed identity package used by the prompt builder.
type Persona struct {
	Config            Config
	SystemPrompt      string
	PromptBody        string
	CanonicalIdentity CanonicalIdentity
}

var placeholderUserPattern = regexp.MustCompile(`(?i)^usuario\s+\d+$`)

// CanonicalIdentity holds resolved identity values before prompt assembly.
type CanonicalIdentity struct {
	AgentName string
	AgentRole string
	UserName  string
}

// LoadPersona reads IDENTITY.md, SOUL.md and USER.md into a single persona.
func LoadPersona(identityPath, soulPath, userPath string) (*Persona, error) {
	identityBytes, err := readPersonaFile(identityPath, "identity")
	if err != nil {
		return nil, err
	}
	soulBytes, err := readPersonaFile(soulPath, "soul")
	if err != nil {
		return nil, err
	}
	userBytes, err := readPersonaFile(userPath, "user")
	if err != nil {
		return nil, err
	}

	config, identityBody, err := parseIdentityFrontmatter(identityBytes)
	if err != nil {
		return nil, err
	}

	promptBody := buildPromptBody(identityBody, soulBytes, userBytes)

	canonicalIdentity := CanonicalIdentity{
		AgentName: canonicalValue(config.Name),
		AgentRole: canonicalValue(config.Role),
		UserName:  canonicalUserName(userBytes),
	}

	sysPrompt := buildSystemPrompt(canonicalIdentity, promptBody)

	return &Persona{
		Config:            config,
		SystemPrompt:      sysPrompt,
		PromptBody:        promptBody,
		CanonicalIdentity: canonicalIdentity,
	}, nil
}

func buildCanonicalIdentityBlock(identity CanonicalIdentity) string {
	lines := []string{
		"# CANONICAL IDENTITY",
		fmt.Sprintf("Nome canonico do agente: %s", canonicalValue(identity.AgentName)),
		fmt.Sprintf("Papel canonico do agente: %s", canonicalValue(identity.AgentRole)),
		fmt.Sprintf("Nome canonico do usuario: %s", canonicalValue(identity.UserName)),
		"Esses fatos canonicos tem prioridade sobre historico conversacional e placeholders.",
	}

	return strings.Join(lines, "\n")
}

func canonicalUserName(userBytes []byte) string {
	name := extractCanonicalUserName(string(userBytes))
	if name == "" {
		return "nao definido"
	}
	return name
}

func extractCanonicalUserName(userContent string) string {
	for _, line := range strings.Split(userContent, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(trimmed), "nome:") {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(trimmed, "Nome:"))
		name = strings.TrimSpace(strings.Trim(name, "*"))
		if name == "" {
			return ""
		}
		if placeholderUserPattern.MatchString(name) {
			return ""
		}

		return name
	}

	return ""
}

func canonicalValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "nao definido"
	}
	return value
}

func buildLongTermMemoryBlock(facts []memory.Fact, notes []memory.Note) string {
	if len(facts) == 0 && len(notes) == 0 {
		return ""
	}

	lines := []string{"# LONG-TERM MEMORY"}
	if len(facts) > 0 {
		lines = append(lines, "Facts:")
		for _, fact := range facts {
			key := strings.TrimSpace(fact.Key)
			value := strings.TrimSpace(fact.Value)
			if key == "" || value == "" {
				continue
			}
			lines = append(lines, fmt.Sprintf("- %s: %s", key, value))
		}
	}
	if len(notes) > 0 {
		lines = append(lines, "Relevant Notes:")
	}
	for _, note := range notes {
		topic := canonicalValue(note.Topic)
		kind := canonicalValue(note.Kind)
		summary := strings.TrimSpace(note.Summary)
		if summary == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("- [%s/%s] %s", topic, kind, summary))
	}

	if len(lines) <= 1 {
		return ""
	}

	return strings.Join(lines, "\n")
}

// ── File helpers ──────────────────────────────────────────────────────────────

func readOptionalFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}

func readPersonaFile(path, kind string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", kind, err)
	}
	return content, nil
}

func parseIdentityFrontmatter(identityBytes []byte) (Config, string, error) {
	var config Config
	parts := bytes.SplitN(identityBytes, []byte("---"), 3)
	if len(parts) != 3 {
		return config, string(bytes.TrimSpace(identityBytes)), nil
	}
	if err := yaml.Unmarshal(parts[1], &config); err != nil {
		return Config{}, "", fmt.Errorf("failed to parse yaml frontmatter: %w", err)
	}
	return config, string(bytes.TrimSpace(parts[2])), nil
}

func buildPromptBody(identityBody string, soulBytes, userBytes []byte) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		identityBody,
		string(bytesTrimSpace(soulBytes)),
		string(bytesTrimSpace(userBytes)),
	)
}

func buildSystemPrompt(identity CanonicalIdentity, promptBody string) string {
	return fmt.Sprintf("%s\n\n%s", buildCanonicalIdentityBlock(identity), promptBody)
}

func bytesTrimSpace(content []byte) []byte {
	return []byte(strings.TrimSpace(string(content)))
}

// RenderSystemPrompt assembles the final prompt with canonical identity and long-term memory.
func (p *Persona) RenderSystemPrompt(identity CanonicalIdentity, facts []memory.Fact, notes []memory.Note) string {
	if p == nil {
		return ""
	}
	sections := []string{buildCanonicalIdentityBlock(identity)}
	if memoryBlock := buildLongTermMemoryBlock(facts, notes); memoryBlock != "" {
		sections = append(sections, memoryBlock)
	}
	sections = append(sections, strings.TrimSpace(p.PromptBody))
	return strings.Join(sections, "\n\n")
}
