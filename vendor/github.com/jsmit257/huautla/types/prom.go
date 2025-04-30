package types

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DataMetrics *prometheus.CounterVec
)

func init() {
	DataMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "cffc",
		Subsystem:   "huautla",
		Name:        "database",
		Help:        "The packages, methods and possible errors when accessing data",
		ConstLabels: prometheus.Labels{},
	}, []string{"db", "pkg", "function", "status"})
}
