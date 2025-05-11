package metrics

import (
	"fmt"
	"strings"

	"github.com/enesanbar/go-service/core/info"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector interface {
	Hit(key string, cacheType string)
	Miss(key string, cacheType string)
}

type Instrumentor struct {
	cacheHits   *prometheus.CounterVec
	cacheMisses *prometheus.CounterVec
}

func NewInstrumentor() *Instrumentor {
	cacheHitsCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_cache_hits_total", strings.ReplaceAll(info.ServiceName, "-", "_")),
		Help: "Thu number of cache hits",
	}, []string{"key", "type"})

	cacheMissesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_cache_misses_total", strings.ReplaceAll(info.ServiceName, "-", "_")),
		Help: "Thu number of cache misses",
	}, []string{"key", "type"})

	err := prometheus.Register(cacheHitsCounter)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		panic(err)
	}

	err = prometheus.Register(cacheMissesCounter)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		panic(err)
	}

	return &Instrumentor{cacheHits: cacheHitsCounter, cacheMisses: cacheMissesCounter}
}

func (i Instrumentor) Hit(key string, cacheType string) {
	i.cacheHits.WithLabelValues(key, cacheType).Inc()
}

func (i Instrumentor) Miss(key string, cacheType string) {
	i.cacheMisses.WithLabelValues(key, cacheType).Inc()
}
