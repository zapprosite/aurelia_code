package actors

import (
	"context"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/streaming"
)

// VADMonitor implementa a detecção de voz local para interrupção.
type VADMonitor struct {
	streaming.BaseActor
	interruptFunc func()
	mu            sync.Mutex
}

func NewVADMonitor() *VADMonitor {
	return &VADMonitor{
		BaseActor: streaming.NewBaseActor("VADMonitor"),
	}
}

func (v *VADMonitor) OnInterrupt(callback func()) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.interruptFunc = callback
}

func (v *VADMonitor) Run(ctx context.Context) error {
	v.Logger.Info("VAD monitor started")
	
	// TODO: Integrar com driver de áudio nativo (portaudio/alsa) ou 
	// polling de nível de decibéis do microfone.
	// Por enquanto, simulamos o loop de escuta.
	
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(100 * time.Millisecond):
			// Simulação de detecção: 
			// Se o Jarvis está falando (Speaker ativo) e detectamos voz, interrompemos.
			if v.isVoiceDetected() {
				v.mu.Lock()
				if v.interruptFunc != nil {
					v.Logger.Warn("Barge-in detected! Triggering interrupt")
					v.interruptFunc()
				}
				v.mu.Unlock()
			}
		}
	}
}

func (v *VADMonitor) isVoiceDetected() bool {
	// TODO: Em SOTA 2026.2, integraremos portaudio para detecção real.
	// Para o "go !" inicial, permitimos interrupção manual ou via signal.
	return false 
}

// TriggerVoice detectado manualmente para testes de integração.
func (v *VADMonitor) TriggerVoice() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.interruptFunc != nil {
		v.Logger.Warn("Manual Barge-in triggered")
		v.interruptFunc()
	}
}
