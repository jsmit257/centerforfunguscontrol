package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
)

// var mtrcs = metrics.ServiceMetrics.MustCurryWith(prometheus.Labels{})

// func authnz(handler http.Handler) http.Handler {
// 	// check auth tokens and whatever other sanity
// 	return nil
// }

func main() {
	r := chi.NewRouter()
	// r.Use(authnz) // someday, maybe more too

	log.SetFormatter(&log.JSONFormatter{})

	log := log.WithFields(log.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	newHuautla(r, log)

	r.Get("/hc", hc)

	// r.Get("/metrics", mtrcs)

	_ = &http.Server{
		// Addr:    fmt.Sprintf("%s:%d", hostAddr, hostPort),
		Handler: r,
	}

}

// not much of a healthcheck, for now
func hc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
