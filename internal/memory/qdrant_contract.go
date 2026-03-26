package memory

import (
	"fmt"
	"strings"
)

func ValidateCanonicalMemoryPayload(payload map[string]any) error {
	return validateQdrantPayload("memory", payload, []string{
		"app_id",
		"repo_id",
		"environment",
		"text",
		"canonical_bot_id",
		"source_system",
		"source_id",
		"domain",
		"ts",
		"version",
	})
}

func ValidateSkillIndexPayload(payload map[string]any) error {
	return validateQdrantPayload("skills", payload, []string{
		"app_id",
		"repo_id",
		"environment",
		"text",
		"name",
		"description",
		"source_system",
		"source_id",
		"domain",
		"ts",
		"version",
	})
}

func validateQdrantPayload(contract string, payload map[string]any, requiredKeys []string) error {
	if len(payload) == 0 {
		return fmt.Errorf("%s payload is empty", contract)
	}
	for _, key := range requiredKeys {
		value, ok := payload[key]
		if !ok {
			return fmt.Errorf("%s payload missing %q", contract, key)
		}
		if isZeroPayloadValue(value) {
			return fmt.Errorf("%s payload has empty %q", contract, key)
		}
	}
	return nil
}

func isZeroPayloadValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	case int:
		return typed == 0
	case int64:
		return typed == 0
	case float64:
		return typed == 0
	default:
		return false
	}
}
