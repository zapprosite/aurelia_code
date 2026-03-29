package agent

import "context"

type MockLLMProvider struct {
	response *ModelResponse
	err      error
}

func (m *MockLLMProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	return m.response, m.err
}

func (m *MockLLMProvider) GenerateStream(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (<-chan StreamResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	ch := make(chan StreamResponse, 10)
	go func() {
		defer close(ch)
		if m.response != nil {
			ch <- StreamResponse{Content: m.response.Content}
		}
		ch <- StreamResponse{Done: true}
	}()
	return ch, nil
}
