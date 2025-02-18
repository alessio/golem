package metrics

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

// StartMetricsServer starts a prometheus server.
// Data Url is at localhost:<port>/metrics/<path>
// Normally you would use /metrics as endpoint and 9090 as port

type TaskMetrics interface {
	Start() error
	RegisterMetric(name string, help string, labels []string, handler MetricHandler) error
	UpdateMetric(name string, value float64, labels ...string) error
	IncrementMetric(name string, labels ...string) error
	DecrementMetric(name string, labels ...string) error
	Name() string
	Stop() error
}

type MetricDetail struct {
	Collector prometheus.Collector
	Handler   MetricHandler
}

type taskMetrics struct {
	path    string
	port    string
	metrics map[string]MetricDetail
	mux     sync.RWMutex
}

func NewTaskMetrics(path string, port string) TaskMetrics {
	return &taskMetrics{
		path:    path,
		port:    port,
		metrics: make(map[string]MetricDetail),
	}
}

func (t *taskMetrics) Name() string {
	return "metrics"
}

func (t *taskMetrics) Start() error {
	router := chi.NewRouter()

	zap.S().Infof("Metrics (prometheus) starting: %v", t.port)

	// Prometheus path
	router.Get(t.path, promhttp.Handler().(http.HandlerFunc))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", t.port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		zap.S().Errorf("Prometheus server error: %v", err)
	} else {
		zap.S().Infof("Prometheus server serving at port %s", t.port)
	}

	return err
}

func (t *taskMetrics) Stop() error {
	return nil
}
