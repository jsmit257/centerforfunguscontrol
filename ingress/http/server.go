package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"

	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"
)

func startServer(cfg *config.Config, r *chi.Mux, wg *sync.WaitGroup, log *log.Entry) *sync.WaitGroup {
	wg.Add(1)
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

func (g *global) staticContent(w http.ResponseWriter, r *http.Request) {
	l := g.l.WithField("method", "staticContent")
	l.Debug("starting work")

	mt := "text/html"
	f := chi.URLParam(r, "f")
	if f == "" {
		f = "./www/test-harness/index.html"
	} else if strings.HasPrefix(r.RequestURI, "/album/") {
		f = "." + r.RequestURI
	} else {
		f = "./www/test-harness" + r.RequestURI
	}

	l = l.WithField("resource", f)
	l.Info("fetching resource")

	// XXX: poor-man's mime typing; could at least use the parent directory
	if strings.HasSuffix(f, ".js") {
		mt = "text/javascript; charset=UTF-8"
	} else if strings.HasSuffix(f, ".css") {
		mt = "text/css; charset=UTF-8"
	} else if strings.HasSuffix(f, ".png") {
		mt = "image/png"
	} else if strings.HasSuffix(f, ".jpg") {
		mt = "image/jpg"
	} else if strings.HasSuffix(f, ".gif") {
		mt = "image/gif"
	} else if strings.HasSuffix(f, ".tiff") {
		mt = "image/tiff"
	}

	l = l.WithField("content-type", mt)
	l.Info("fetching resource")

	if result, err := os.ReadFile(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(f))
		l.WithError(err).Error("fetching resource")
	} else {
		w.Header().Add("Content-type", mt)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result)
	}
	l.Info("done work")
}
