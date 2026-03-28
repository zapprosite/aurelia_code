package aflow

import (
	"context"
	"github.com/kocar/aurelia/internal/purity/alog"
)

// Task é a unidade atômica de trabalho no aflow
type Task func(ctx context.Context) error

// Flow gerencia o encadeamento de tarefas
type Flow struct {
	name  string
	steps []Task
}

// New inicia um novo fluxo soberano
func New(name string) *Flow {
	return &Flow{
		name: name,
	}
}

// Step adiciona uma tarefa ao fluxo (inspirado no PicoFlow >>)
func (f *Flow) Step(task Task) *Flow {
	f.steps = append(f.steps, task)
	return f
}

// Run executa o fluxo sequencialmente com stop-on-error
func (f *Flow) Run(ctx context.Context) error {
	alog.Info("starting sovereign flow", alog.With("flow", f.name), alog.With("steps", len(f.steps)))
	
	for i, step := range f.steps {
		alog.Debug("executing step", alog.With("step", i+1))
		if err := step(ctx); err != nil {
			alog.Error("flow interrupted by step error", alog.With("flow", f.name), alog.With("step", i+1), alog.With("err", err))
			return err
		}
	}

	alog.Info("flow completed successfully", alog.With("flow", f.name))
	return nil
}
