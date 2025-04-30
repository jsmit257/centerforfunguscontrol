package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func NewHandler(reg *prometheus.Registry) http.HandlerFunc {
	reg.MustRegister(DataMetrics)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP
}
