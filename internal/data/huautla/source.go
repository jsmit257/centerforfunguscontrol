package huautla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) PostStrainSource(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostStrainSource")
	defer ms.end()
	defer r.Body.Close()

	var s types.Source

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.AddStrainSource(r.Context(), &g, s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add source")
	} else {
		ms.send(w, g, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PostEventSource(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostEventSource")
	defer ms.end()
	defer r.Body.Close()

	var e types.Event

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.AddEventSource(r.Context(), &g, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add source")
	} else {
		ms.send(w, g, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchSource(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchLifecycleEvent")
	defer ms.end()
	defer r.Body.Close()

	var s types.Source

	if gID := chi.URLParam(r, "id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.ChangeSource(r.Context(), &g, s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change source")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteSource(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteSource")
	defer ms.end()

	if gID := chi.URLParam(r, "g_id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if sID := chi.URLParam(r, "s_id"); sID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if sID, err := url.QueryUnescape(sID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.RemoveSource(r.Context(), &g, types.UUID(sID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove source")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}
