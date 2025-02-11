package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jsmit257/huautla/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	Cid     ctxkey = "cid"
	Metrics ctxkey = "metrics"
	Log     ctxkey = "log"
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

func GetContextCID(ctx context.Context) types.CID {
	if result, ok := ctx.Value(Cid).(types.CID); ok {
		return result
	}
	return types.CID(fmt.Sprintf("context has no cid attribute: %#v", ctx))
}

func GetContextLog(ctx context.Context) *logrus.Entry {
	if result, ok := ctx.Value(Log).(*logrus.Entry); ok {
		return result
	}

	l := logrus.WithFields(logrus.Fields{
		"ctx":   ctx,
		"bogus": true,
	})

	l.
		WithError(fmt.Errorf("context has no log attribute: %#v", ctx)).
		Error("getting context")

	return l
}

func GetContextMetrics(ctx context.Context) *prometheus.CounterVec {
	if result, ok := ctx.Value(Metrics).(*prometheus.CounterVec); ok {
		return result
	}

	return ServiceMetrics.MustCurryWith(prometheus.Labels{
		"url":    "/missing/metrics/context/attribute",
		"proto":  "ERROR",
		"method": "metrics.GetContextMetrics",
	})
}
