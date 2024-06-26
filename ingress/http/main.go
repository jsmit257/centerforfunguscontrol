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
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGQUIT}

// var mtrcs = metrics.ServiceMetrics.MustCurryWith(prometheus.Labels{})

// func authnz(handler http.Handler) http.Handler {
// 	// check auth tokens and whatever other sanity
// 	return nil
// }

type global struct {
	l *log.Entry
}

func main() {
	cfg := config.NewConfig()
	log.SetLevel(log.DebugLevel) // TODO: grab this from the config
	log.SetFormatter(&log.JSONFormatter{})

	log := log.WithFields(log.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	g := &global{log}

	wg := &sync.WaitGroup{}

	r := chi.NewRouter()
	// r.Use(authnz) // someday, maybe more too

	r.Get("/", g.staticContent)
	r.Get("/css/{f}", g.staticContent)
	r.Get("/css/images/{f}", g.staticContent)
	r.Get("/css/images/background/{f}", g.staticContent)
	r.Get("/js/{f}", g.staticContent)
	r.Get("/images/{f}", g.staticContent)
	r.Get("/photos/{f}", g.staticContent)

	newHC(r)
	newHuautla(cfg, r, log)
	newServer(cfg, r, wg, log)

	wg.Wait()

	log.Info("done")
}

func trap(log *log.Entry) {
	trapped := make(chan os.Signal, len(traps))

	signal.Notify(trapped, traps...)

	log.WithField("sig", <-trapped).Info("stopping app with signal")
}
