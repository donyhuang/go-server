package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	prom "gitlab.nongchangshijie.com/go-base/server/pkg/prometheus"
	"strconv"
	"time"
)

const serverNamespace = "http_server"

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler(label map[string]string) gin.HandlerFunc {
	metricServerReqDur := prom.NewHistogramVec(&prometheus.HistogramOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"path"}, label)

	metricServerReqCodeTotal := prom.NewCounterVec(&prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
	}, []string{"path", "code"}, label)
	return func(c *gin.Context) {
		startTime := time.Now()
		defer func() {
			metricServerReqDur.WithLabelValues(c.Request.URL.Path).Observe(float64(time.Since(startTime).Milliseconds()))
			metricServerReqCodeTotal.WithLabelValues(c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
		}()
		c.Next()
	}
}
