package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla/types"
)

func WrapContext(log *logrus.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			cid := types.CID(uuid.NewString())

			log = log.WithFields(logrus.Fields{
				"method": r.Method,
				"remote": r.RemoteAddr, // is this necessary?
				"url":    r.RequestURI, // this, or getRoutePattern and fill in the blanks with more fields?
				"cid":    cid,
			})

			w.Header().Set("Cid", string(cid))
			r = r.WithContext(context.WithValue(
				context.WithValue(
					context.WithValue(
						r.Context(),
						Metrics,
						ServiceMetrics.MustCurryWith(prometheus.Labels{
							"proto":  r.Proto,
							"method": r.Method,
							"url":    getRoutePattern(r),
						}),
					),
					Log,
					log,
				),
				Cid,
				cid,
			))

			log.Info("starting request")

			next.ServeHTTP(w, r)

			log.WithField("duration", time.Since(start).String()).Info("finished request")
		})
	}
}

// shamelessly copied from https://github.com/go-chi/chi/issues/270#issuecomment-479184559
func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		return "!" + pattern // leaving the bang to see if this ever happens (and how?)
	}

	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if rctx.Routes.Match(tctx, r.Method, routePath) {
		return tctx.RoutePattern()
	}

	// better than logging or panicing, as long as it never happens
	return "!!" + routePath
}

var MockServiceContext = context.WithValue(
	context.WithValue(
		context.WithValue(
			context.Background(),
			Metrics,
			ServiceMetrics.MustCurryWith(prometheus.Labels{
				"proto":  "r.Proto",
				"method": "r.Method",
				"url":    "url",
			}),
		),
		Log,
		logrus.WithFields(logrus.Fields{}),
	),
	Cid,
	types.CID("cid"),
)
