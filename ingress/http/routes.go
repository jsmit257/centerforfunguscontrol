package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	us "github.com/jsmit257/userservice/shared/v1"
)

func loginRedirect(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusForbidden)
}

func authn(host string, port uint16) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := metrics.GetContextLog(r.Context())
			if c, err := r.Cookie("us-authn"); err == http.ErrNoCookie {
				loginRedirect(w, r)
			} else if newc, sc := us.CheckValid(host, port, c); sc != http.StatusFound {
				l.WithFields(logrus.Fields{
					"sc":     sc,
					"cookie": newc,
				}).Info("status")
				loginRedirect(w, r)
			} else {
				// resp := r.Response
				// if resp != nil && resp.Request != nil {
				// 	if found, _ := regexp.Match("/otp/.*", []byte(resp.Request.RequestURI)); !found {
				// 		l.
				// 			WithField("uri", resp.Request.RequestURI).
				// 			Warn("uri didn't match")
				// 		w.Header().Set("Authn-Pad", resp.Header.Get("Authn-Pad"))
				// 	} else {
				// 		l.Info("setting header")
				// 		w.Header().Set("Authn-Pad", resp.Header.Get("Authn-Pad"))
				// 	}
				// }
				http.SetCookie(w, newc)
				next.ServeHTTP(w, r)
			}
		})
	}
}

func newHuautla(cfg *config.Config, ha *huautla.HuautlaAdaptor, l *logrus.Entry) *chi.Mux {
	l = l.WithField("database", "huautla")
	r := chi.NewRouter()

	r.Use(metrics.WrapContext(l))

	if cfg.AuthnHost != "" && cfg.AuthnPort != 0 {
		r.Use(authn(cfg.AuthnHost, cfg.AuthnPort))
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

	r.Get("/reports/lifecycle/{id}", ha.GetLifecycleReport)
	r.Get("/reports/generation/{id}", ha.GetGenerationReport)
	r.Get("/reports/strain/{id}", ha.GetStrainReport)
	r.Get("/reports/substrate/{id}", ha.GetSubstrateReport)
	r.Get("/reports/eventtype/{id}", ha.GetEventTypeReport)
	r.Get("/reports/vendor/{id}", ha.GetVendorReport)

	r.Patch("/ts/{table}/{id}", ha.PatchTS)
	r.Patch("/undel/{table}/{id}", ha.Undel)

	r.Get("/metrics", metrics.NewHandler())

	return r
}
