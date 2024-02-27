package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
)

var traps = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGHUP,
	syscall.SIGQUIT}

// var mtrcs = metrics.ServiceMetrics.MustCurryWith(prometheus.Labels{})

// func authnz(handler http.Handler) http.Handler {
// 	// check auth tokens and whatever other sanity
// 	return nil
// }

func main() {
	cfg := config.NewConfig()

	log.SetLevel(log.InfoLevel) // TODO: grab this from the config
	log.SetFormatter(&log.JSONFormatter{})

	log := log.WithFields(log.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	wg := &sync.WaitGroup{}

	r := chi.NewRouter()
	// r.Use(authnz) // someday, maybe more too

	r.Get("/", staticContent)
	r.Get("/css/{f}", staticContent)
	r.Get("/css/images/{f}", staticContent)
	r.Get("/js/{f}", staticContent)
	r.Get("/images/{f}", staticContent)

	newHuautla(cfg, r, log)
	newHC(r)
	newServer(cfg, r, wg, log)

	wg.Wait()

	log.Info("done")
}

func trap(log *log.Entry) {
	trapped := make(chan os.Signal, len(traps))

	signal.Notify(trapped, traps...)

	log.WithField("sig", <-trapped).Info("stopping app with signal")
}
