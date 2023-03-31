package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.nongchangshijie.com/go-base/server/pkg/proc"
)

func NewCounterVec(cfg *prometheus.CounterOpts, labels []string, constLabel prometheus.Labels) *prometheus.CounterVec {
	if cfg == nil {
		return nil
	}

	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.Namespace,
		Subsystem:   cfg.Subsystem,
		Name:        cfg.Name,
		Help:        cfg.Help,
		ConstLabels: constLabel,
	}, labels)
	prometheus.MustRegister(vec)
	proc.AddDoneFn(func() {
		prometheus.Unregister(vec)
	})
	return vec
}
