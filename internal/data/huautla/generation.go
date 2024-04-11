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

func (ha *HuautlaAdaptor) GetGenerationIndex(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetGenerationIndex")
	defer ms.end()

	if g, err := ha.db.SelectGenerationIndex(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generations")
	} else {
		ms.l.WithField("generations", g).Error("response?")
		ms.send(w, g, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetGeneration(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetGeneration")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostGeneration(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostGeneration")
	defer ms.end()
	defer r.Body.Close()

	var g types.Generation

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &g); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err = ha.db.InsertGeneration(r.Context(), g, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert generation")
	} else {
		ms.send(w, g, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchGeneration(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchGeneration")
	defer ms.end()
	defer r.Body.Close()

	var g types.Generation

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &g); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if g, err = ha.db.UpdateGeneration(r.Context(), g, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update generation")
	} else {
		ms.send(w, g, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteGeneration(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteGeneration")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete generation")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
