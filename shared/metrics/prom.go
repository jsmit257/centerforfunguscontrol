package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/jsmit257/userservice/shared/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	DataMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "cffc",
		Subsystem:   "api",
		Name:        "database",
		Help:        "The packages, methods and possible errors when accessing data",
		ConstLabels: prometheus.Labels{},
	}, []string{"db", "pkg", "function", "status"})

	ServiceMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "cffc",
		Subsystem:   "api",
		Name:        "router",
		Help:        "Service requests tracked by ???",
		ConstLabels: prometheus.Labels{},
	}, []string{"url", "proto", "method", "sc"})
)

func NewHandler() http.HandlerFunc {
	reg := prometheus.NewRegistry()

	reg.MustRegister(DataMetrics)
	reg.MustRegister(ServiceMetrics)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP
}

func (v *core) NewDataTracker(ctx context.Context, fn string) Tracker {
	cid := ctx.Value(Cid).(shared.CID)

	l := v.log.WithFields(logrus.Fields{
		"function": fn,
		"cid":      cid,
	})
	l.Info("starting work")

	return &track{
		l: l,
		m: v.metrics.MustCurryWith(prometheus.Labels{"function": fn}),
		r: &trackresults{},
		s: time.Now().UTC(),
	}
}
