package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Общее количество HTTP запросов",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Время выполнения HTTP запроса",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	prometheus.MustRegister(m.httpRequestsTotal, m.httpRequestDuration)
	return m
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		m.httpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		m.httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}
