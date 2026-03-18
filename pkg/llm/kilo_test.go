package llm

import "testing"

func TestNewKiloProvider(t *testing.T) {
	t.Parallel()

	provider := NewKiloProvider("secret", "gpt-5.4")
	if provider == nil {
		t.Fatal("expected provider to be created")
	}
}
