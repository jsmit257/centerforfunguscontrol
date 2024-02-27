package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newHC(r *chi.Mux) {
	r.Get("/hc", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}
