package memory

import (
"context"
)

type MemoryOS interface {
GetContext(ctx context.Context, sessionID string, role string) (string, error)
Push(ctx context.Context, sessionID string, contribution string) error
}

type SimpleMemoryOS struct {
// Add backend (Redis/VectorDB) later
}

func NewSimpleMemoryOS() *SimpleMemoryOS {
return &SimpleMemoryOS{}
}

func (m *SimpleMemoryOS) GetContext(ctx context.Context, sessionID string, role string) (string, error) {
// Mock implementation for now
return "Recent context for " + role, nil
}

func (m *SimpleMemoryOS) Push(ctx context.Context, sessionID string, contribution string) error {
// Mock implementation for now
return nil
}
