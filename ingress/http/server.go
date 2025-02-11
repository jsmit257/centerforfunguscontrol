package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"

	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"
)

func startServer(cfg *config.Config, r *chi.Mux, wg *sync.WaitGroup, log *log.Entry) *sync.WaitGroup {
	defer wg.Done()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort),
		Handler: r,
	}

	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Fatal("http server didn't start properly")
		}
	}(srv)

	log.Info("server started, waiting for traps")

	trapped := make(chan os.Signal, len(traps))

	signal.Notify(trapped, traps...)

	log.WithField("sig", <-trapped).Info("stopping app with signal")

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeout); err != nil {
		log.WithError(err).Error("error stopping server")
	}

	log.Debug("http server shutdown complete")

	return wg
}
