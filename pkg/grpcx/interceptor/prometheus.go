package interceptor

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	prom "gitlab.nongchangshijie.com/go-base/server/pkg/prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const serverNamespace = "grpc_server"

// UnaryPrometheusInterceptor reports the statistics to the prometheus server.
func UnaryPrometheusInterceptor(label map[string]string) func(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	metricServerReqDur := prom.NewHistogramVec(&prometheus.HistogramOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"method"}, label)

	metricServerReqCodeTotal := prom.NewCounterVec(&prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc server requests code count.",
	}, []string{"method", "code"}, label)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		metricServerReqDur.WithLabelValues(info.FullMethod).Observe(float64(time.Since(startTime).Milliseconds()))
		metricServerReqCodeTotal.WithLabelValues(info.FullMethod, strconv.Itoa(int(status.Code(err)))).Inc()
		return resp, err
	}
}
