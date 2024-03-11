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

func (ha *HuautlaAdaptor) GetAllStages(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllStages")
	defer ms.end()

	if stages, err := ha.db.SelectAllStages(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch stages")
	} else {
		ms.send(w, stages, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetStage(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetStage")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SelectStage(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch stage")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostStage(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostStage")
	defer ms.end()
	defer r.Body.Close()

	var s types.Stage

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if s, err = ha.db.InsertStage(r.Context(), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert stage")
	}

	ms.send(w, s, http.StatusCreated)
}

func (ha *HuautlaAdaptor) PatchStage(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchStage")
	defer ms.end()
	defer r.Body.Close()

	var s types.Stage

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateStage(r.Context(), types.UUID(id), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update stage")
	}

	ms.send(w, nil, http.StatusNoContent)
}

func (ha *HuautlaAdaptor) DeleteStage(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteStage")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteStage(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete stage")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
