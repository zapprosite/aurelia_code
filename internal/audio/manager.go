package audio

import (
	"context"
	"fmt"
	"sync"
)

// AudioManager orquestra o fluxo de dados do microfone/sistema para o Voice Loop.
type AudioManager struct {
	mu      sync.RWMutex
	buffers map[string]*StreamBuffer
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		buffers: make(map[string]*StreamBuffer),
	}
}

// RegisterBuffer cria um novo canal de stream identificado.
func (m *AudioManager) RegisterBuffer(id string, capacity int) *StreamBuffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	buf := NewStreamBuffer(capacity)
	m.buffers[id] = buf
	return buf
}

// GetBuffer retorna um buffer existente.
func (m *AudioManager) GetBuffer(id string) (*StreamBuffer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	buf, ok := m.buffers[id]
	if !ok {
		return nil, fmt.Errorf("buffer %s not found", id)
	}
	return buf, nil
}

// StartPipeline (Skeleton) inicia o processamento assíncrono.
func (m *AudioManager) StartPipeline(ctx context.Context) error {
	// Implementação futura: Integração com Whisper (STT) e Kokoro (TTS)
	return nil
}
