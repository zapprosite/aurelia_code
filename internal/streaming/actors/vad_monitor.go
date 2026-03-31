package actors

import (
	"context"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/streaming"
)

// VADMonitor implementa a detecção de voz local para interrupção.
type VADMonitor struct {
	streaming.BaseActor
	interruptFunc func()
	mu            sync.Mutex
	Trigger       chan bool // Canal para notificar ativação de voz (Her-Mode)
}

func NewVADMonitor() *VADMonitor {
	return &VADMonitor{
		BaseActor: streaming.NewBaseActor("VADMonitor"),
		Trigger:   make(chan bool, 1),
	}
}

func (v *VADMonitor) OnInterrupt(callback func()) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.interruptFunc = callback
}

func (v *VADMonitor) Run(ctx context.Context) error {
	v.Logger.Info("VAD monitor searching for gateway socket", "path", "/tmp/aurelia-voice.sock")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := v.listenToSocket(ctx)
			if err != nil {
				v.Logger.Warn("Socket connection failed, retrying in 2s", "err", err)
				time.Sleep(2 * time.Second)
			}
		}
	}
}

func (v *VADMonitor) listenToSocket(ctx context.Context) error {
	conn, err := net.Dial("unix", "/tmp/aurelia-voice.sock")
	if err != nil {
		return err
	}
	defer conn.Close()

	v.Logger.Info("Connected to Voice Gateway socket")
	buf := make([]byte, 1024)
	
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := conn.Read(buf)
			if err != nil {
				return err
			}
			if n == 0 {
				return fmt.Errorf("socket closed by gateway")
			}

			// Processa energia básica para interrupção (Barge-in)
			if v.detectVoice(buf[:n]) {
				v.mu.Lock()
				// Ativa o gatilho se estiver em modo Her-Mode ouvindo
				select {
				case v.Trigger <- true:
				default:
				}

				if v.interruptFunc != nil {
					v.Logger.Warn("Voice activity detected! Triggering Barge-in.")
					v.interruptFunc()
				}
				v.mu.Unlock()
			}
		}
	}
}

func (v *VADMonitor) detectVoice(data []byte) bool {
	// Cálculo simples de energia para RMS (Root Mean Square)
	// Como os dados são int16 (S16_LE), processamos pares de bytes.
	var sum float64
	count := 0
	for i := 0; i < len(data)-1; i += 2 {
		val := int16(data[i]) | (int16(data[i+1]) << 8)
		sum += float64(val) * float64(val)
		count++
	}
	if count == 0 {
		return false
	}
	rms := math.Sqrt(sum / float64(count))
	
	// log de energia para debug em SOTA 2026.1
	if rms > 100 {
		v.Logger.Debug("Voice signal level", "rms", rms)
	}

	// Threshold industrial (ajustável via config futuramente)
	return rms > 500 
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
