package agent

import (
	"strings"
	"unicode"
)

// ToolCatalogEntry enriquece uma definição de Tool com exemplos de casos de uso
// e palavras-chave para matching semântico sem dependência de embeddings externos.
type ToolCatalogEntry struct {
	Tool
	// Keywords são termos que associam esta tool a intenções específicas.
	// São derivados da descrição + schema automaticamente e podem ser enriquecidos
	// manualmente via WithKeywords().
	Keywords []string
	// Score é o escore de relevância calculado para a tarefa atual (ephemeral).
	Score int
}

// ToolCatalog gerencia o mapa semântico de tools disponíveis no registry.
// Em vez de enviar todas as 77+ tools para a LLM em todo prompt,
// o catalog filtra as top-K mais relevantes para a tarefa.
type ToolCatalog struct {
	entries []ToolCatalogEntry
}

// NewToolCatalog constrói o catálogo a partir de um ToolRegistry existente.
// Extrai automaticamente keywords de Name e Description de cada tool.
func NewToolCatalog(registry *ToolRegistry) *ToolCatalog {
	defs := registry.GetDefinitions()
	entries := make([]ToolCatalogEntry, 0, len(defs))
	for _, def := range defs {
		entry := ToolCatalogEntry{
			Tool:     def,
			Keywords: extractKeywords(def.Name, def.Description),
		}
		entries = append(entries, entry)
	}
	return &ToolCatalog{entries: entries}
}

// MatchForTask retorna até k tools mais relevantes para o prompt dado.
// Usa scoring léxico: match por tokens do prompt contra keywords da tool.
// Ferramentas core de raciocínio (think, read, write) recebem boost.
func (c *ToolCatalog) MatchForTask(prompt string, k int) []Tool {
	if k <= 0 || len(c.entries) == 0 {
		return nil
	}

	// Tokeniza o prompt
	promptTokens := tokenize(prompt)
	if len(promptTokens) == 0 {
		// prompt vazio → retornar as k primeiras (tools core)
		return topK(c.entries, k)
	}

	// Realiza scoring para cada entry
	scored := make([]ToolCatalogEntry, len(c.entries))
	copy(scored, c.entries)

	for i := range scored {
		scored[i].Score = scoreEntry(&scored[i], promptTokens)
	}

	// Ordena por score decrescente (insertion sort para slices pequenos)
	for i := 1; i < len(scored); i++ {
		for j := i; j > 0 && scored[j].Score > scored[j-1].Score; j-- {
			scored[j], scored[j-1] = scored[j-1], scored[j]
		}
	}

	// Coleta as top-k com score > 0; se não houver suficientes, inclui as core
	result := make([]Tool, 0, k)
	for _, e := range scored {
		if len(result) >= k {
			break
		}
		if e.Score > 0 {
			result = append(result, e.Tool)
		}
	}

	// Preenche com ferramentas core caso não haja k matches
	if len(result) < k {
		coreNames := coreToolNames()
		for _, e := range scored {
			if len(result) >= k {
				break
			}
			if coreNames[e.Name] && !containsTool(result, e.Name) {
				result = append(result, e.Tool)
			}
		}
	}

	return result
}

// AllTools retorna todas as tools sem filtro (equivalente ao comportamento atual).
func (c *ToolCatalog) AllTools() []Tool {
	all := make([]Tool, len(c.entries))
	for i, e := range c.entries {
		all[i] = e.Tool
	}
	return all
}

// Len retorna o total de tools no catálogo.
func (c *ToolCatalog) Len() int {
	return len(c.entries)
}

// --- helpers internos ---

func scoreEntry(e *ToolCatalogEntry, promptTokens []string) int {
	score := 0
	keywordIndex := make(map[string]bool, len(e.Keywords))
	for _, kw := range e.Keywords {
		keywordIndex[kw] = true
	}
	for _, pt := range promptTokens {
		if keywordIndex[pt] {
			score++
		}
		// Match parcial: token do prompt contido no nome da tool
		if strings.Contains(strings.ToLower(e.Name), pt) {
			score += 2 // peso maior para match no nome
		}
	}
	// Boost para tools core de raciocínio
	if coreToolNames()[e.Name] {
		score++
	}
	return score
}

// tokenize converte um texto em tokens normalizados (lowercase, sem pontuação).
func tokenize(text string) []string {
	text = strings.ToLower(text)
	// Split por pontuação e espaço
	f := func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r)
	}
	parts := strings.FieldsFunc(text, f)

	// Remove stopwords PT-BR + EN
	stopwords := stopwordSet()
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) > 2 && !stopwords[p] {
			result = append(result, p)
		}
	}
	return result
}

// extractKeywords extrai keywords do nome e descrição de uma tool.
func extractKeywords(name, description string) []string {
	combined := name + " " + description
	return tokenize(combined)
}

func topK(entries []ToolCatalogEntry, k int) []Tool {
	result := make([]Tool, 0, k)
	for i, e := range entries {
		if i >= k {
			break
		}
		result = append(result, e.Tool)
	}
	return result
}

func containsTool(tools []Tool, name string) bool {
	for _, t := range tools {
		if t.Name == name {
			return true
		}
	}
	return false
}

// coreToolNames são ferramentas sempre incluídas pois fazem parte do raciocínio básico.
func coreToolNames() map[string]bool {
	return map[string]bool{
		"run_command":  true,
		"read_file":    true,
		"write_file":   true,
		"list_dir":     true,
		"search_files": true,
		"grep":         true,
	}
}

// stopwordSet retorna um conjunto de stopwords ignoradas no matching.
func stopwordSet() map[string]bool {
	words := []string{
		// PT-BR
		"que", "com", "para", "por", "uma", "um", "como", "não", "mas",
		"vai", "vou", "tem", "ser", "está", "ele", "ela", "nos", "nas",
		"dos", "das", "aos", "esse", "essa", "isso", "isto", "seu", "sua",
		// EN
		"the", "and", "for", "that", "this", "with", "from", "are", "was",
		"not", "but", "have", "you", "can", "will", "what", "how", "its",
	}
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
