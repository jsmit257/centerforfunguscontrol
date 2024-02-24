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
		Name: "datailed_data_access_metrics",
		Help: "The packages, methods and possible errors when accessing data",
	}, []string{"pkg", "method", "err"})

}

func NewHandler(reg *prometheus.Registry) http.HandlerFunc {
	reg.MustRegister(DataMetrics)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP
}
