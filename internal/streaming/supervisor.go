package streaming

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Supervisor gerencia a saúde e o reinício automático de atores.
type Supervisor struct {
	actors []Actor
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
	logger *slog.Logger
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		logger: slog.Default().With("component", "supervisor", "system", "aurelia"),
	}
}

func (s *Supervisor) Add(a Actor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.actors = append(s.actors, a)
}

func (s *Supervisor) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	
	s.logger.Info("Supervisor starting SAP actor monitoring")
	
	for _, a := range s.actors {
		go s.monitorActor(a)
	}
	
	return nil
}

func (s *Supervisor) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Supervisor) monitorActor(a Actor) {
	name := a.Name()
	for {
		s.logger.Info("Starting actor", "name", name)
		
		// Executa o ator
		err := a.Run(s.ctx)
		
		select {
		case <-s.ctx.Done():
			s.logger.Info("Supervisor stopping, actor exited", "name", name)
			return
		default:
			if err != nil {
				s.logger.Error("Actor crashed, initiating self-healing", "name", name, "err", err)
			} else {
				s.logger.Warn("Actor exited unexpectedly without error, restarting", "name", name)
			}
			
			// Backoff exponencial simples para evitar loops de crash infinitos
			time.Sleep(2 * time.Second)
		}
	}
}

func (s *Supervisor) HealthCheck() error {
	// Futura implementação: verificar se cada ator responde a um ping interno
	return nil
}
