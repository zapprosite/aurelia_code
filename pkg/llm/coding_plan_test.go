package llm

import (
	"reflect"
	"testing"
)

func TestNewZAIProvider_UsesCodingPlanEndpoint(t *testing.T) {
	t.Parallel()

	provider := NewZAIProvider("secret", "glm-5")
	baseURL := reflect.ValueOf(provider).Elem().FieldByName("baseURL").String()
	if baseURL != "https://api.z.ai/api/coding/paas/v4/chat/completions" {
		t.Fatalf("baseURL = %q", baseURL)
	}
}

func TestNewAlibabaProvider_UsesCodingPlanEndpoint(t *testing.T) {
	t.Parallel()

	provider := NewAlibabaProvider("secret", "qwen3-coder-plus")
	baseURL := reflect.ValueOf(provider).Elem().FieldByName("baseURL").String()
	if baseURL != "https://coding-intl.dashscope.aliyuncs.com/v1/chat/completions" {
		t.Fatalf("baseURL = %q", baseURL)
	}
}
