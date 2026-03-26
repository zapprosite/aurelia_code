package memory

import "testing"

func TestValidateCanonicalMemoryPayload_RequiresCanonicalFields(t *testing.T) {
	payload := map[string]any{
		"app_id":           "aurelia",
		"repo_id":          "aurelia",
		"environment":      "local",
		"text":             "ok",
		"canonical_bot_id": "aurelia_code",
		"source_system":    "voice",
		"source_id":        "voice:job-1",
		"domain":           "system",
		"ts":               int64(1),
		"version":          1,
	}

	if err := ValidateCanonicalMemoryPayload(payload); err != nil {
		t.Fatalf("ValidateCanonicalMemoryPayload() error = %v", err)
	}
}

func TestValidateCanonicalMemoryPayload_FailsWithoutSourceID(t *testing.T) {
	payload := map[string]any{
		"app_id":           "aurelia",
		"repo_id":          "aurelia",
		"environment":      "local",
		"text":             "ok",
		"canonical_bot_id": "aurelia_code",
		"source_system":    "voice",
		"domain":           "system",
		"ts":               int64(1),
		"version":          1,
	}

	if err := ValidateCanonicalMemoryPayload(payload); err == nil {
		t.Fatal("expected missing source_id to fail")
	}
}
