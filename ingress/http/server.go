package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"

	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"
)

func newServer(cfg *config.Config, r *chi.Mux, wg *sync.WaitGroup, log *log.Entry) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort),
		Handler: r,
	}

	wg.Add(1)

	go func(srv *http.Server, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.
				// WithField("cfg", cfg).
				WithError(err).
				Fatal("http server didn't start properly")
			panic(err)
		}
		log.Debug("http server shutdown complete")
	}(srv, wg)

	trap(log)

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeout); err != nil {
		log.
			// WithField("cfg", cfg).
			WithError(err).
			Error("error stopping server")
	}
}

func staticContent(w http.ResponseWriter, r *http.Request) {
	mt := "text/html"
	f := chi.URLParam(r, "f")
	if f == "" {
		f = "./www/test-harness/index.html"
	} else {
		f = "./www/test-harness" + r.RequestURI
	}

	// XXX: poor-man's mime typing; could at least use the parent directory
	if strings.HasSuffix(f, ".js") {
		mt = "text/javascript; charset=UTF-8"
	} else if strings.HasSuffix(f, ".css") {
		mt = "text/css; charset=UTF-8"
	} else if strings.HasSuffix(f, ".png") {
		mt = "image/png"
	}

	if result, err := os.ReadFile(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(f))
	} else {
		w.Header().Add("Content-type", mt)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result)
	}
}
