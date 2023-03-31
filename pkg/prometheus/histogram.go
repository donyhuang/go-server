package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.nongchangshijie.com/go-base/server/pkg/proc"
)

func NewHistogramVec(cfg *prometheus.HistogramOpts, labels []string, constLabel prometheus.Labels) *prometheus.HistogramVec {
	if cfg == nil {
		return nil
	}

	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.Namespace,
		Subsystem:   cfg.Subsystem,
		Name:        cfg.Name,
		Help:        cfg.Help,
		Buckets:     cfg.Buckets,
		ConstLabels: constLabel,
	}, labels)
	prometheus.MustRegister(vec)
	proc.AddDoneFn(func() {
		prometheus.Unregister(vec)
	})
	return vec
}
