package memory

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ContextAssembler struct {
	qdrantURL          string
	qdrantAPIKey       string
	qdrantCollection   string
	markdownCollection string
	ollamaURL          string
	embedURL           string
	embeddingModel     string
	mem                *MemoryManager
	client             *http.Client
}

func NewContextAssembler(qdrantURL, apiKey, collection, markdownCollection, embeddingModel, ollamaURL string, mem *MemoryManager) *ContextAssembler {
	return &ContextAssembler{
		qdrantURL:          strings.TrimRight(strings.TrimSpace(qdrantURL), "/"),
		qdrantAPIKey:       strings.TrimSpace(apiKey),
		qdrantCollection:   strings.TrimSpace(collection),
		markdownCollection: strings.TrimSpace(markdownCollection),
		ollamaURL:          strings.TrimRight(strings.TrimSpace(ollamaURL), "/"),
		embedURL:           strings.TrimRight(strings.TrimSpace(ollamaURL), "/") + "/api/embed",
		embeddingModel:     strings.TrimSpace(embeddingModel),
		mem:                mem,
		client:             NewSemanticHTTPClient(10 * time.Second),
	}
}

// AssembleContext returns a formatted markdown string joining semantic matches and recent facts/notes.
func (a *ContextAssembler) AssembleContext(ctx context.Context, query string) string {
	return a.AssembleContextForBot(ctx, "", query)
}

// AssembleContextForBot scopes semantic results to a specific bot when provided.
func (a *ContextAssembler) AssembleContextForBot(ctx context.Context, botID, query string) string {
	var sb strings.Builder

	a.appendSemanticSection(ctx, &sb, botID, query, a.qdrantCollection, "Arquivos Históricos", "[memory/degraded]", semanticFormatOptions{
		minScore: 0.40,
		maxItems: 5,
	})
	if a.shouldIncludeMarkdownBrain(botID) {
		a.appendSemanticSection(ctx, &sb, botID, query, a.markdownCollection, "Markdown Brain", "[markdown_brain/degraded]", semanticFormatOptions{
			minScore:       0.38,
			maxItems:       5,
			uniqueBySource: true,
		})
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

type semanticFormatOptions struct {
	minScore       float64
	maxItems       int
	uniqueBySource bool
}

func (a *ContextAssembler) appendSemanticSection(
	ctx context.Context,
	sb *strings.Builder,
	botID, query, collection, titleBase, degradedPrefix string,
	format semanticFormatOptions,
) {
	points, mode, err := a.searchQdrant(ctx, botID, collection, query, max(format.maxItems*2, 5))
	rendered := formatSemanticPointsWithOptions(points, format)
	if rendered != "" {
		sb.WriteString(sectionTitle(titleBase, mode))
		sb.WriteString(rendered)
		sb.WriteString("\n\n")
		if err != nil && mode == "lexical-fallback" {
			sb.WriteString(degradedPrefix)
			sb.WriteString(" Busca semântica indisponível; usando fallback lexical.\n\n")
		}
		return
	}
	if err != nil {
		sb.WriteString(fmt.Sprintf("%s Busca semântica indisponível: %v\n\n", degradedPrefix, err))
	}
}

func (a *ContextAssembler) searchQdrant(ctx context.Context, botID, collection, text string, limit int) ([]SemanticPoint, string, error) {
	if a.qdrantURL == "" || strings.TrimSpace(collection) == "" || text == "" {
		return nil, "", nil
	}

	vector, semanticErr := EmbedText(ctx, a.client, a.ollamaURL, a.embeddingModel, text)
	if semanticErr == nil {
		points, err := SearchSemantic(ctx, a.client, a.qdrantURL, collection, a.qdrantAPIKey, vector, limit)
		if err == nil {
			points = FilterPointsByCanonicalBotID(points, botID)
			if len(points) > 0 {
				return points, "semantic", nil
			}
		} else {
			semanticErr = err
		}
	}

	points, scrollErr := ScrollPoints(ctx, a.client, a.qdrantURL, collection, a.qdrantAPIKey, max(limit*10, 50))
	if scrollErr != nil {
		if semanticErr != nil {
			return nil, "degraded", fmt.Errorf("%v; fallback lexical failed: %w", semanticErr, scrollErr)
		}
		return nil, "degraded", scrollErr
	}

	points = FilterPointsByCanonicalBotID(points, botID)
	filtered := FilterPointsLexical(points, text)
	if len(filtered) == 0 {
		return nil, "degraded", semanticErr
	}
	return filtered, "lexical-fallback", semanticErr
}

func (a *ContextAssembler) shouldIncludeMarkdownBrain(botID string) bool {
	if a == nil || strings.TrimSpace(a.markdownCollection) == "" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(botID)) {
	case "", "aurelia", "aurelia_code":
		return true
	default:
		return false
	}
}

func sectionTitle(base, mode string) string {
	switch mode {
	case "lexical-fallback":
		return base + " (Qdrant Fallback Lexical):\n"
	default:
		return base + " (Qdrant Semantic Search):\n"
	}
}

func formatSemanticPoints(points []SemanticPoint, minScore float64) string {
	return formatSemanticPointsWithOptions(points, semanticFormatOptions{minScore: minScore})
}

func formatSemanticPointsWithOptions(points []SemanticPoint, opts semanticFormatOptions) string {
	var out strings.Builder
	seenSources := make(map[string]struct{})
	written := 0

	for _, hit := range points {
		if hit.Score > 0 && hit.Score < opts.minScore {
			continue
		}
		if opts.maxItems > 0 && written >= opts.maxItems {
			break
		}

		payload := NormalizeSemanticPayload(hit.Payload)
		if opts.uniqueBySource {
			sourceID := firstString(payload, "source_id")
			if sourceID != "" {
				if _, ok := seenSources[sourceID]; ok {
					continue
				}
				seenSources[sourceID] = struct{}{}
			}
		}

		text := firstString(payload, "text", "content", "summary", "transcript", "message", "body")
		if text == "" {
			text = ExtractSearchableText(payload)
		}
		if text == "" {
			continue
		}

		var labels []string
		if path := firstString(payload, "repo_path", "vault_path"); path != "" {
			labels = append(labels, path)
		}
		if section := firstString(payload, "section"); section != "" {
			labels = append(labels, section)
		}
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
			written++
			continue
		}
		out.WriteString(fmt.Sprintf("> %s%s\n", prefix, text))
		written++
	}

	return out.String()
}

func (a *ContextAssembler) embed(ctx context.Context, text string) ([]float32, error) {
	return EmbedText(ctx, a.client, a.ollamaURL, a.embeddingModel, text)
}
