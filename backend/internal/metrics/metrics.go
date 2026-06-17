// Package metrics — все Prometheus-коллекторы в одном месте (observability §2, конвенции лейблов).
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Кол-во HTTP-запросов",
	}, []string{"route", "status"})

	HTTPDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Латентность HTTP-запросов",
		Buckets: prometheus.DefBuckets,
	}, []string{"route"})

	IngestJobs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ingest_jobs_total",
		Help: "Завершённые задачи обработки фото по статусу",
	}, []string{"status"})

	IngestLastSuccess = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ingest_last_success_timestamp_seconds",
		Help: "Unix-время последней успешной обработки фото",
	})

	MatchesCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "matches_created_total",
		Help: "Созданные совпадения по типу",
	}, []string{"type"})

	IngestQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ingest_queue_depth",
		Help: "Глубина очереди обработки (pending+processing)",
	})
)
