package memory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ContextAssembler struct {
	qdrantURL        string
	qdrantAPIKey     string
	qdrantCollection string
	ollamaURL        string
	embedURL         string
	embeddingModel   string
	mem              *MemoryManager
	client           *http.Client
}

func NewContextAssembler(qdrantURL, apiKey, collection, embeddingModel, ollamaURL string, mem *MemoryManager) *ContextAssembler {
	return &ContextAssembler{
		qdrantURL:        strings.TrimRight(strings.TrimSpace(qdrantURL), "/"),
		qdrantAPIKey:     strings.TrimSpace(apiKey),
		qdrantCollection: strings.TrimSpace(collection),
		ollamaURL:        strings.TrimRight(strings.TrimSpace(ollamaURL), "/"),
		embedURL:         strings.TrimRight(strings.TrimSpace(ollamaURL), "/") + "/api/embed",
		embeddingModel:   strings.TrimSpace(embeddingModel),
		mem:              mem,
		client:           NewSemanticHTTPClient(10 * time.Second),
	}
}

// AssembleContext returns a formatted markdown string joining semantic matches and recent facts/notes.
func (a *ContextAssembler) AssembleContext(ctx context.Context, query string) string {
	return a.AssembleContextForBot(ctx, "", query)
}

// AssembleContextForBot scopes semantic results to a specific bot when provided.
func (a *ContextAssembler) AssembleContextForBot(ctx context.Context, botID, query string) string {
	var sb strings.Builder

	// 1. Fetch from Qdrant
	if qdrantRes, mode, err := a.searchQdrant(ctx, botID, query); qdrantRes != "" {
		title := "Arquivos Históricos (Qdrant Semantic Search):\n"
		if mode == "lexical-fallback" {
			title = "Arquivos Históricos (Qdrant Fallback Lexical):\n"
		}
		sb.WriteString(title)
		sb.WriteString(qdrantRes)
		sb.WriteString("\n\n")
		if err != nil && mode == "lexical-fallback" {
			sb.WriteString("[memory/degraded] Busca semântica indisponível; usando fallback lexical.\n\n")
		}
	} else if err != nil {
		sb.WriteString(fmt.Sprintf("[memory/degraded] Busca semântica indisponível: %v\n\n", err))
	}

	// 2. Fetch from SQLite Notes
	if a.mem != nil {
		if recentNotes, err := a.mem.GetGlobalTopics(ctx, 10); err == nil && len(recentNotes) > 0 {
			sb.WriteString("Tópicos e Notas Locais (SQLite):\n")
			for _, note := range recentNotes {
				sb.WriteString(fmt.Sprintf("- [%s] %s\n", note.Kind, note.Topic))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

func (a *ContextAssembler) searchQdrant(ctx context.Context, botID, text string) (string, string, error) {
	if a.qdrantURL == "" || text == "" {
		return "", "", nil
	}

	vector, semanticErr := EmbedText(ctx, a.client, a.ollamaURL, a.embeddingModel, text)
	if semanticErr == nil {
		points, err := SearchSemantic(ctx, a.client, a.qdrantURL, a.qdrantCollection, a.qdrantAPIKey, vector, 5)
		if err == nil {
			points = FilterPointsByCanonicalBotID(points, botID)
			if rendered := formatSemanticPoints(points, 0.40); rendered != "" {
				return rendered, "semantic", nil
			}
		} else {
			semanticErr = err
		}
	}

	points, scrollErr := ScrollPoints(ctx, a.client, a.qdrantURL, a.qdrantCollection, a.qdrantAPIKey, 50)
	if scrollErr != nil {
		if semanticErr != nil {
			return "", "degraded", fmt.Errorf("%v; fallback lexical failed: %w", semanticErr, scrollErr)
		}
		return "", "degraded", scrollErr
	}

	points = FilterPointsByCanonicalBotID(points, botID)
	filtered := FilterPointsLexical(points, text)
	if len(filtered) == 0 {
		return "", "degraded", semanticErr
	}
	return formatSemanticPoints(filtered, 0), "lexical-fallback", semanticErr
}

func formatSemanticPoints(points []SemanticPoint, minScore float64) string {
	var out strings.Builder
	for _, hit := range points {
		if hit.Score > 0 && hit.Score < minScore {
			continue
		}

		payload := NormalizeSemanticPayload(hit.Payload)
		text := ExtractSearchableText(payload)
		if text == "" {
			continue
		}

		var labels []string
		if botID := firstString(payload, "canonical_bot_id"); botID != "" {
			labels = append(labels, botID)
		}
		if domain := firstString(payload, "domain"); domain != "" {
			labels = append(labels, domain)
		}
		prefix := ""
		if len(labels) > 0 {
			prefix = "[" + strings.Join(labels, "/") + "] "
		}

		if hit.Score > 0 {
			out.WriteString(fmt.Sprintf("> (Score: %.2f) %s%s\n", hit.Score, prefix, text))
			continue
		}
		out.WriteString(fmt.Sprintf("> %s%s\n", prefix, text))
	}

	return out.String()
}

func (a *ContextAssembler) embed(ctx context.Context, text string) ([]float32, error) {
	return EmbedText(ctx, a.client, a.ollamaURL, a.embeddingModel, text)
}
