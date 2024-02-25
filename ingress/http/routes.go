package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla/types"

	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
)

// var mtrcs = metrics.ServiceMetrics.MustCurryWith(prometheus.Labels{})

func authnz(handler http.Handler) http.Handler {
	// check auth tokens and whatever other sanity
	return nil
}

func NewInstance(hostAddr string, hostPort uint16, mtrcs http.HandlerFunc, log *log.Entry) *http.Server {
	r := chi.NewRouter()
	// r.Use(authnz) // someday, maybe more too

	log = log.WithField("ingress", "http")

	newHuautla(r, log)

	r.Get("/hc", hc)

	r.Get("/metrics", mtrcs)

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", hostAddr, hostPort),
		Handler: r,
	}
}

func newHuautla(r *chi.Mux, l *log.Entry) {
	ha, _ := huautla.New(
		&types.Config{
			PGHost: "localhost",
			PGPort: 5432,
			PGUser: "postgres",
			PGPass: "root",
			PGSSL:  "disable",
		},
		l.WithField("database", "huautla"),
		nil)

	r.Get("/vendors", ha.GetAllVendors)
	r.Get("/vendor/{id}", ha.GetVendor)
	r.Post("/vendor", ha.PostVendor)
	r.Patch("/vendor/{id}", ha.PatchVendor)
	r.Delete("/vendor/{id}", ha.DeleteVendor)

	r.Get("/stages", ha.GetAllStages)
	r.Get("/stage/{id}", ha.GetStage)
	r.Post("/stage", ha.PostStage)
	r.Patch("/stage/{id}", ha.PatchStage)
	r.Delete("/stage/{id}", ha.DeleteStage)

	r.Get("/eventtypes", ha.GetAllEventTypes)
	r.Get("/eventtype/{id}", ha.GetEventType)
	r.Post("/eventtype", ha.PostEventType)
	r.Patch("/eventtype/{id}", ha.PatchEventType)
	r.Delete("/eventtype/{id}", ha.DeleteEventType)

	r.Get("/ingredients", ha.GetAllIngredients)
	r.Get("/ingredient/{id}", ha.GetIngredient)
	r.Post("/ingredient", ha.PostIngredient)
	r.Patch("/ingredient/{id}", ha.PatchIngredient)
	r.Delete("/ingredient/{id}", ha.DeleteIngredient)

	r.Get("/strains", nil)
	r.Get("/strain/{id}", nil)
	r.Post("/strain", nil)
	r.Patch("/strain/{id}", nil)
	r.Delete("/strain/{id}", nil)

	r.Get("/lifecycle/{id}", nil)
	r.Post("/lifecycle", nil)
	r.Patch("/lifecycle/{id}", nil)
	r.Delete("/lifecycle/{id}", nil)

	r.Post("/lifecycle/{id}/events", nil)
	r.Patch("/lifecycle/{id}/events", nil)
	r.Delete("/lifecycle/{id}/events/{id}", nil)
}

// not much of a healthcheck, for now
func hc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
