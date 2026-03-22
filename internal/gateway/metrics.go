package gateway

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type gatewayMetrics struct {
	requests  *prometheus.CounterVec
	failures  *prometheus.CounterVec
	fallbacks *prometheus.CounterVec
	latency   *prometheus.HistogramVec
	breakers  *prometheus.GaugeVec
	budgets   *prometheus.GaugeVec
	tokens    *prometheus.CounterVec
	costUSD   *prometheus.CounterVec
}

var (
	metricsOnce sync.Once
	metricsSet  *gatewayMetrics
)

func defaultMetrics() *gatewayMetrics {
	metricsOnce.Do(func() {
		metricsSet = &gatewayMetrics{
			requests: mustRegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "requests_total",
				Help:      "Total de requests roteadas pelo gateway.",
			}, []string{"lane", "provider", "model", "result"})),
			failures: mustRegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "failures_total",
				Help:      "Falhas por lane/provider/model.",
			}, []string{"lane", "provider", "model"})),
			fallbacks: mustRegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "fallbacks_total",
				Help:      "Quantidade de fallbacks disparados pelo gateway.",
			}, []string{"from_lane", "to_lane"})),
			latency: mustRegisterCollector(prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "route_latency_seconds",
				Help:      "Latencia por lane/provider/model.",
				Buckets:   prometheus.DefBuckets,
			}, []string{"lane", "provider", "model", "result"})),
			breakers: mustRegisterCollector(prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "breaker_state",
				Help:      "Estado do circuit breaker por provider/model. closed=0, half-open=0.5, open=1.",
			}, []string{"provider", "model"})),
			budgets: mustRegisterCollector(prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "budget_usage_ratio",
				Help:      "Uso do budget por lane em relacao ao limite hard.",
			}, []string{"lane"})),
			tokens: mustRegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "tokens_total",
				Help:      "Total de tokens processados pelo gateway.",
			}, []string{"lane", "direction"})),
			costUSD: mustRegisterCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "aurelia",
				Subsystem: "gateway",
				Name:      "cost_usd_total",
				Help:      "Custo total em USD por lane.",
			}, []string{"lane"})),
		}
	})
	return metricsSet
}

func mustRegisterCollector[T prometheus.Collector](collector T) T {
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
