package main

import (
	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla/types"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
)

func newHuautla(cfg *config.Config, r *chi.Mux, l *log.Entry) {
	ha, _ := huautla.New(
		&types.Config{
			PGHost: cfg.HuautlaHost,
			PGPort: uint(cfg.HuautlaPort),
			PGUser: cfg.HuautlaUser,
			PGPass: cfg.HuautlaPass,
			PGSSL:  cfg.HuautlaSSL,
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

	r.Get("/substrates", ha.GetAllSubstrates)
	r.Get("/substrate/{id}", ha.GetSubstrate)
	r.Post("/substrate", ha.PostSubstrate)
	r.Patch("/substrate/{id}", ha.PatchSubstrate)
	r.Delete("/substrate/{id}", ha.DeleteSubstrate)

	r.Get("/ingredients", ha.GetAllIngredients)
	r.Get("/ingredient/{id}", ha.GetIngredient)
	r.Post("/ingredient", ha.PostIngredient)
	r.Patch("/ingredient/{id}", ha.PatchIngredient)
	r.Delete("/ingredient/{id}", ha.DeleteIngredient)

	r.Post("/substrate/{id}/ingredients", ha.PostSubstrateIngredient)
	r.Patch("/substrate/{su_id}/ingredients/{ig_id}", ha.PatchSubstrateIngredient)
	r.Delete("/substrate/{su_id}/ingredients/{ig_id}", ha.DeleteSubstrateIngredient)

	r.Get("/strains", ha.GetAllStrains)
	r.Get("/strain/{id}", ha.GetStrain)
	r.Post("/strain", ha.PostStrain)
	r.Patch("/strain/{id}", ha.PatchStrain)
	r.Delete("/strain/{id}", ha.DeleteStrain)

	r.Get("/strainattributenames", ha.GetStrainAttributeNames)
	r.Post("/strain/{id}/attribute/{at_name}/{at_value}", ha.PostStrainAttribute)
	r.Patch("/strain/{st_id}/attribute/{at_name}/{at_value}", ha.PatchStrainAttribute)
	r.Delete("/strain/{st_id}/attribute/{at_id}", ha.DeleteStrainAttribute)

	r.Get("/lifecycles", nil) // TODO: needs to be implemented in huautla
	r.Get("/lifecycle/{id}", ha.GetLifecycle)
	r.Post("/lifecycle", ha.PostLifecycle)
	r.Patch("/lifecycle/{id}", ha.PatchLifecycle)
	r.Delete("/lifecycle/{id}", ha.DeleteLifecycle)

	r.Post("/lifecycle/{id}/events", ha.PostEvent)
	r.Patch("/lifecycle/{lc_id}/events/{ev_id}", ha.PatchEvent)
	r.Delete("/lifecycle/{lc_id}/events/{ev_id}", ha.DeleteEvent)
}
