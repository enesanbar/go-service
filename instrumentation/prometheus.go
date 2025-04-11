package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Service implements UseCase interface
type Service struct {
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.SummaryVec
	requestTotal    *prometheus.CounterVec
	responseSize    *prometheus.SummaryVec
}

// NewPrometheusService create a new prometheus service
func NewPrometheusService() (*Service, error) {

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "host", "method", "code"})

	requestSize := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_request_size_bytes",
		Help: "The size of HTTP requests",
	}, []string{"handler", "host", "method", "code"})

	requestTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Thu number of HTTP requests received",
	}, []string{"handler", "host", "method", "code"})

	responseSize := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_response_size_bytes",
		Help: "The size of HTTP responses",
	}, []string{"handler", "host", "method", "code"})

	err := prometheus.Register(requestDuration)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(requestSize)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(responseSize)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(requestTotal)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	return &Service{
		requestDuration: requestDuration,
		requestSize:     requestSize,
		requestTotal:    requestTotal,
		responseSize:    responseSize,
	}, nil
}

func (s *Service) SaveMetrics(h *HTTP) {
	s.requestDuration.
		WithLabelValues(h.Handler, h.Host, h.Method, h.StatusCode).
		Observe(h.Duration)

	s.requestSize.
		WithLabelValues(h.Handler, h.Host, h.Method, h.StatusCode).
		Observe(float64(h.RequestSize))

	s.responseSize.
		WithLabelValues(h.Handler, h.Host, h.Method, h.StatusCode).
		Observe(float64(h.ResponseSize))

	s.requestTotal.
		WithLabelValues(h.Handler, h.Host, h.Method, h.StatusCode).
		Inc()
}
