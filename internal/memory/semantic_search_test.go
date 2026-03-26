package memory

import (
	"testing"
	"time"
)

func TestExtractSearchableTextPrefersCanonicalFields(t *testing.T) {
	payload := map[string]any{
		"title":      "Status",
		"content":    "Conteudo canonico",
		"transcript": "Nao deveria vencer",
	}

	got := ExtractSearchableText(payload)
	want := "Status — Conteudo canonico"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeSemanticPayloadAddsCanonicalFields(t *testing.T) {
	payload := map[string]any{
		"bot_id":     "aurelia",
		"transcript": "texto herdado",
		"created_at": "2026-03-25T20:19:00Z",
	}

	normalized := NormalizeSemanticPayload(payload)
	if normalized["canonical_bot_id"] != "aurelia" {
		t.Fatalf("expected canonical_bot_id from bot_id, got %#v", normalized["canonical_bot_id"])
	}
	if normalized["text"] != "texto herdado" {
		t.Fatalf("expected text to be normalized from transcript, got %#v", normalized["text"])
	}
	if _, ok := normalized["ts"]; !ok {
		t.Fatalf("expected ts to be populated from created_at")
	}
}

func TestExtractPayloadTimestampSupportsUnixAndRFC3339(t *testing.T) {
	fromUnix := ExtractPayloadTimestamp(map[string]any{"ts": float64(1774396800)})
	if fromUnix.IsZero() {
		t.Fatalf("expected unix ts to be parsed")
	}

	fromRFC3339 := ExtractPayloadTimestamp(map[string]any{"created_at": "2026-03-25T20:19:00Z"})
	if fromRFC3339.IsZero() {
		t.Fatalf("expected created_at to be parsed")
	}

	if !fromRFC3339.Equal(time.Date(2026, 3, 25, 20, 19, 0, 0, time.UTC)) {
		t.Fatalf("unexpected parsed time: %s", fromRFC3339)
	}
}

func TestFilterPointsLexicalMatchesNormalizedText(t *testing.T) {
	points := []SemanticPoint{
		{ID: "1", Payload: NormalizeSemanticPayload(map[string]any{"transcript": "fluxo da caixa validado"})},
		{ID: "2", Payload: NormalizeSemanticPayload(map[string]any{"text": "pipeline de obras"})},
	}

	filtered := FilterPointsLexical(points, "caixa")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 lexical hit, got %d", len(filtered))
	}
	if filtered[0].ID != "1" {
		t.Fatalf("expected point 1, got %#v", filtered[0].ID)
	}
}

func TestFilterPointsByCanonicalBotIDSupportsLeaderAlias(t *testing.T) {
	points := []SemanticPoint{
		{ID: "1", Payload: NormalizeSemanticPayload(map[string]any{"canonical_bot_id": "aurelia_code", "text": "lider"})},
		{ID: "2", Payload: NormalizeSemanticPayload(map[string]any{"canonical_bot_id": "controle-db", "text": "governanca"})},
	}

	filtered := FilterPointsByCanonicalBotID(points, "aurelia")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 alias-matched point, got %d", len(filtered))
	}
	if filtered[0].ID != "1" {
		t.Fatalf("expected point 1, got %#v", filtered[0].ID)
	}
}
