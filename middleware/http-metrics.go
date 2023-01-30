package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	labels = []string{"status", "method", "path", "host"}

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count_total",
			Help: "Total number of HTTP requests made.",
		}, labels,
	)

	// DefBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request latencies in seconds.",
		}, labels,
	)

	requestInflight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_inflight",
			Help: "HTTP request in flight count",
		}, []string{"method", "path", "host"},
	)
)

func init() {
	prometheus.MustRegister(requestCount, requestDuration, requestInflight)
}

func HttpMetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path
		host := c.Request.Host

		requestInflight.WithLabelValues(method, path, host).Add(1)
		defer requestInflight.WithLabelValues(method, path, host).Add(-1)

		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		lvs := []string{status, method, path, host}
		requestCount.WithLabelValues(lvs...).Inc()
		requestDuration.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
	}
}

func HttpMetricHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
