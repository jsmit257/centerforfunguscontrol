package huautla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) PostEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostEvent")
	defer ms.end()

	var e types.Event

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.AddEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PatchEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchEvent")
	defer ms.end()

	var e types.Event

	if lcID := chi.URLParam(r, "lc_id"); lcID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if evID := chi.URLParam(r, "ev_id"); evID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.ChangeEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteEvent")
	defer ms.end()

	if lcID := chi.URLParam(r, "lc_id"); lcID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if evID := chi.URLParam(r, "ev_id"); evID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.RemoveEvent(r.Context(), &l, types.UUID(evID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add strainattribute")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}
