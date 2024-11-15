package main

import (
	"os"
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
	syscall.SIGQUIT,
}

func main() {
	cfg := config.NewConfig()
	log.SetLevel(log.DebugLevel) // TODO: grab this from the config
	log.SetFormatter(&log.JSONFormatter{})

	log := log.WithFields(log.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	wg := &sync.WaitGroup{}

	r := chi.NewRouter()

	newHuautla(cfg, r, log)
	newHC(r)

	startServer(cfg, r, wg, log).Wait()

	log.Info("done")

	os.Exit(0)
}
