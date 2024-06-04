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

func (ha *HuautlaAdaptor) PostLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostEvent")
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
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.AddLifecycleEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.send(w, l, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchLifecycleEvent")
	defer ms.end()
	defer r.Body.Close()

	var e types.Event

	if lcID := chi.URLParam(r, "lc_id"); lcID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if lcID, err := url.QueryUnescape(lcID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if _, err := ha.db.ChangeLifecycleEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change event")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteLifecycleEvent")
	defer ms.end()

	if lcID := chi.URLParam(r, "lc_id"); lcID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if lcID, err := url.QueryUnescape(lcID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if evID := chi.URLParam(r, "ev_id"); evID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if evID, err := url.QueryUnescape(evID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.RemoveLifecycleEvent(r.Context(), &l, types.UUID(evID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove event")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostGenerationEvent")
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
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.AddGenerationEvent(r.Context(), &g, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.send(w, g, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchGenerationEvent")
	defer ms.end()
	defer r.Body.Close()

	var e types.Event

	if gID := chi.URLParam(r, "id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if _, err := ha.db.ChangeGenerationEvent(r.Context(), &g, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change event")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteGenerationEvent")
	defer ms.end()

	if gID := chi.URLParam(r, "g_id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required lifecycle id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if evID := chi.URLParam(r, "ev_id"); evID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if evID, err := url.QueryUnescape(evID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.RemoveGenerationEvent(r.Context(), &g, types.UUID(evID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove lifecycle")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}
