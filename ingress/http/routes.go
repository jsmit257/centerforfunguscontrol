package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	us "github.com/jsmit257/userservice/shared/v1"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
)

func authn(host string, port uint16, logon string) func(next http.Handler) http.Handler {
	noCookie := http.Header{}
	noCookie.Set("Location", logon)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := metrics.GetContextLog(r.Context())

			c, _ := r.Cookie("us-authn") // only err is ErrNoCookie and we don't care
			if newc, header, sc := us.CheckValid(host, port, c); sc == http.StatusNoContent {
				http.SetCookie(w, newc)
				next.ServeHTTP(w, r)
			} else if sc == http.StatusTemporaryRedirect {
				l.WithFields(logrus.Fields{
					"sc":     sc,
					"cookie": newc,
					"header": header,
				}).Warn("cookie isn't valid, redirecting")

				if newc != nil {
					http.SetCookie(w, newc)
				}

				for key, values := range header {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.WriteHeader(http.StatusForbidden)
			} else if sc == http.StatusInternalServerError {
				l.WithFields(logrus.Fields{
					"sc":     sc,
					"cookie": newc,
					"header": header,
				}).Error("userservice internal error")

				// happens when:
				// - a new request can't be created (bad host/port?)
				// - or sent (lots of reasons)
				// - or if the returned cookie is absent or couldn't be parsed
				// - also, redis errors
				// most of these come from userservice itself, and none of them
				// are likely to get better if we try again, so ???
				http.SetCookie(w, newc)
				w.WriteHeader(sc)
				// here's a strong argument for including err in the return from CheckValid
				_, _ = w.Write([]byte("userservice internal error, probably"))
			} else {
				// is there any reason not to let 500s fall through to this case?
				l.WithFields(logrus.Fields{
					"sc":     sc,
					"cookie": newc,
					"header": header,
				}).Error("unexpected error")

				http.SetCookie(w, newc)
				w.WriteHeader(sc)
				_, _ = w.Write([]byte("unexpected error"))
			}
		})
	}
}

func newHuautla(cfg *config.Config, ha *huautla.HuautlaAdaptor, l *logrus.Entry) *chi.Mux {
	l = l.WithField("database", "huautla")
	r := chi.NewRouter()

	r.Use(metrics.WrapContext(l))

	if cfg.AuthnHost != "" && cfg.AuthnPort != 0 {
		r.Use(authn(cfg.AuthnHost, cfg.AuthnPort, cfg.AuthnPath))
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
	r.Patch("/strain/{st_id}/attribute/{at_id}", ha.PatchStrainAttribute)
	r.Delete("/strain/{st_id}/attribute/{at_id}", ha.DeleteStrainAttribute)

	r.Get("/strain/{id}/generation", ha.GetGeneratedStrain)
	r.Patch("/strain/{sid}/generation/{gid}", ha.PatchGeneratedStrain)
	r.Delete("/strain/{sid}/generation", ha.DeleteGeneratedStrain)

	r.Get("/lifecycles", ha.GetLifecycleIndex)
	r.Get("/lifecycle/{id}", ha.GetLifecycle)
	r.Post("/lifecycle", ha.PostLifecycle)
	r.Patch("/lifecycle/{id}", ha.PatchLifecycle)
	r.Delete("/lifecycle/{id}", ha.DeleteLifecycle)

	r.Get("/events/{o_id}", ha.GetObservableEvents)
	r.Get("/event/{ev_id}", ha.GetEvent)
	r.Post("/event/{o_id}", ha.PostEvent)
	// FIXME: one of the following two is wrong; the first makes the most sense
	// since the url param is never read, the second makes the front end URL
	// construction simpler and follows the pattern of all the other PATCHs on
	// this router; if we decide to go with the first versoin, then all of the
	// patch routes and their corresponding handlers should be updated if needed;
	// also the front end needs to appending an id on PATCH, so we can't really
	// remove anything here until cffc-web::rc1 is complete
	r.Patch("/events/{o_id}", ha.PatchEvent)
	r.Patch("/event/{o_id}/{ignored}", ha.PatchEvent)
	r.Delete("/events/{o_id}/{ev_id}", ha.DeleteEvent)

	// DEPRECTED API; will be removed in the future
	r.Post("/lifecycle/{lc_id}/events", ha.PostLifecycleEvent)
	// same note as the generations version of patch-events below
	r.Patch("/lifecycle/{lc_id}/events/{ev_id}", ha.PatchEvent)
	r.Patch("/lifecycle/{lc_id}/events", ha.PatchLifecycleEvent)
	r.Delete("/lifecycle/{lc_id}/events/{ev_id}", ha.DeleteLifecycleEvent)
	// END DEPRECATED API

	r.Get("/generations", ha.GetGenerationIndex)
	r.Get("/generation/{id}", ha.GetGeneration)
	r.Post("/generation", ha.PostGeneration)
	r.Patch("/generation/{id}", ha.PatchGeneration)
	r.Delete("/generation/{id}", ha.DeleteGeneration)

	// DEPRECTED API; will be removed in the future
	r.Post("/generation/{g_id}/events", ha.PostGenerationEvent)
	// keeping the /generation/{g_id} prefix b/c it matches the pattern used
	// by the front end post/patch pattern; not sure if that's the best idea;
	// either way the generation part is ignored
	r.Patch("/generation/{g_id}/events/{ev_id}", ha.PatchEvent)
	r.Patch("/generation/{g_id}/events", ha.PatchGenerationEvent)
	r.Delete("/generation/{g_id}/events/{ev_id}", ha.DeleteGenerationEvent)
	// END DEPRECATED API

	r.Post("/generation/{id}/sources/{origin}", ha.PostSource)
	r.Patch("/generation/{g_id}/sources/{origin}/{s_id}", ha.PatchSource)
	r.Delete("/generation/{g_id}/sources/{s_id}", ha.DeleteSource)

	r.Get("/notes/{o_id}", ha.GetNotes)
	r.Post("/notes/{o_id}", ha.PostNote)
	r.Patch("/notes/{o_id}/{n_id}", ha.PatchNote)
	r.Delete("/notes/{o_id}/{id}", ha.DeleteNote)

	r.Get("/photoalbum", ha.GetPhotosIndex)
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

	// units save copies of config values when they're initialized, so there's no
	// POST /settings; it would be a bad idea anyway
	r.Get("/settings/{name}", settings(cfg))

	r.Get("/metrics", metrics.NewHandler())

	return r
}

func settings(cfg *config.Config) http.HandlerFunc {
	temp, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	var result map[string]any
	err = json.Unmarshal(temp, &result)
	if err != nil {
		panic(err)
	}

	all, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s string
		if key := chi.URLParam(r, "name"); key == "*" {
			w.WriteHeader(http.StatusOK)
			s = string(all)
		} else if body, ok := result[key]; !ok {
			w.WriteHeader(http.StatusNotFound) // dicey: NoContent or NotFound?
			s = fmt.Sprintf(`{"error": "no value for key: %s"}`, key)
		} else if value, err := json.Marshal(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s = fmt.Sprintf(`{"error": "%q"}`, err)
		} else {
			w.WriteHeader(http.StatusOK)
			s = fmt.Sprintf(`{"value": %s}`, value)
		}
		_, _ = w.Write([]byte(s))
	}
}
