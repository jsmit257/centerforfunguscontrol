package main

import (
	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla/types"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
)

func newHuautla(cfg *config.Config, r *chi.Mux, l *log.Entry) {
	ha, err := huautla.New(
		&types.Config{
			PGHost: cfg.HuautlaHost,
			PGPort: uint(cfg.HuautlaPort),
			PGUser: cfg.HuautlaUser,
			PGPass: cfg.HuautlaPass,
			PGSSL:  cfg.HuautlaSSL,
		},
		l.WithField("database", "huautla"),
		nil)
	if err != nil {
		panic(err)
	}

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
	r.Post("/strain/{id}/attribute", ha.PostStrainAttribute)
	r.Patch("/strain/{id}/attribute", ha.PatchStrainAttribute)
	r.Delete("/strain/{st_id}/attribute/{at_id}", ha.DeleteStrainAttribute)

	r.Get("/strain/{id}/generation", ha.GetGeneratedStrain)
	r.Patch("/strain/{sid}/generation/{gid}", ha.PatchGeneratedStrain)
	r.Delete("/strain/{sid}/generation", ha.DeleteGeneratedStrain)

	r.Get("/lifecycles", ha.GetLifecycleIndex)
	r.Get("/lifecycle/{id}", ha.GetLifecycle)
	r.Post("/lifecycle", ha.PostLifecycle)
	r.Patch("/lifecycle/{id}", ha.PatchLifecycle)
	r.Delete("/lifecycle/{id}", ha.DeleteLifecycle)

	r.Post("/lifecycle/{id}/events", ha.PostLifecycleEvent)
	r.Patch("/lifecycle/{lc_id}/events", ha.PatchLifecycleEvent)
	r.Delete("/lifecycle/{lc_id}/events/{ev_id}", ha.DeleteLifecycleEvent)

	r.Get("/generations", ha.GetGenerationIndex)
	r.Get("/generation/{id}", ha.GetGeneration)
	r.Post("/generation", ha.PostGeneration)
	r.Patch("/generation/{id}", ha.PatchGeneration)
	r.Delete("/generation/{id}", ha.DeleteGeneration)

	r.Post("/generation/{id}/events", ha.PostGenerationEvent)
	r.Patch("/generation/{id}/events", ha.PatchGenerationEvent)
	r.Delete("/generation/{g_id}/events/{ev_id}", ha.DeleteGenerationEvent)

	r.Post("/generation/{id}/sources/strain", ha.PostStrainSource)
	r.Post("/generation/{id}/sources/event", ha.PostEventSource)
	r.Patch("/generation/{id}/sources", ha.PatchSource)
	r.Delete("/generation/{g_id}/sources/{s_id}", ha.DeleteSource)

	r.Get("/notes/{o_id}", ha.GetNotes)
	r.Post("/notes/{o_id}", ha.PostNote)
	r.Patch("/notes/{o_id}", ha.PatchNote)
	r.Delete("/notes/{o_id}/{id}", ha.DeleteNote)

	r.Get("/photos/{o_id}", ha.GetPhotos)
	r.Post("/photos/{o_id}", ha.PostPhoto)
	r.Patch("/photos/{o_id}/{id}", ha.PatchPhoto)
	r.Delete("/photos/{o_id}/{id}", ha.DeletePhoto)

	r.Get("/reports/lifecycles", ha.GetLifecyclesByAttrs)
	r.Get("/reports/generations", ha.GetGenerationsByAttrs)

	r.Get("/reports/lifecycle/{id}", ha.GetLifecycleReport)
	r.Get("/reports/generation/{id}", ha.GetGenerationReport)
	r.Get("/reports/strain/{id}", ha.GetStrainReport)
	r.Get("/reports/substrate/{id}", ha.GetSubstrateReport)
	// r.Get("/reports/eventtype/{id}", ha.GetEventtypeReport)
	// r.Get("/reports/vendor/{id}", ha.GetVendorReport)
}
