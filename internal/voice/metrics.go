package voice

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type voiceMetrics struct {
	queueDepth    prometheus.Gauge
	processed     *prometheus.CounterVec
	dispatches    *prometheus.CounterVec
	fallbacks     *prometheus.CounterVec
	mirrorFailure prometheus.Counter
	budgetUsage   prometheus.Gauge
}

var (
	voiceMetricsOnce sync.Once
	voiceMetricsSet  *voiceMetrics
)

func defaultVoiceMetrics() *voiceMetrics {
	voiceMetricsOnce.Do(func() {
		voiceMetricsSet = &voiceMetrics{
			queueDepth: mustRegisterVoiceCollector(prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "queue_depth",
				Help:      "Quantidade de jobs pendentes/processando no spool de voz.",
			})),
			processed: mustRegisterVoiceCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "jobs_total",
				Help:      "Jobs de voz processados por resultado.",
			}, []string{"result"})),
			dispatches: mustRegisterVoiceCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "dispatch_total",
				Help:      "Despachos do pipeline de voz para o runtime principal.",
			}, []string{"result"})),
			fallbacks: mustRegisterVoiceCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "fallback_total",
				Help:      "Uso de fallback no pipeline de voz.",
			}, []string{"reason"})),
			mirrorFailure: mustRegisterVoiceCollector(prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "mirror_failures_total",
				Help:      "Falhas de mirror para Supabase/Qdrant no pipeline de voz.",
			})),
			budgetUsage: mustRegisterVoiceCollector(prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "aurelia",
				Subsystem: "voice",
				Name:      "budget_usage_ratio",
				Help:      "Uso do budget diario do STT em relacao ao limite hard.",
			})),
		}
	})
	return voiceMetricsSet
}

func mustRegisterVoiceCollector[T prometheus.Collector](collector T) T {
	if err := prometheus.Register(collector); err != nil {
		if already, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := already.ExistingCollector.(T); ok {
				return existing
			}
		}
		panic(err)
	}
	return collector
}
