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

func (ha *HuautlaAdaptor) GetLifecycleIndex(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllLifecycles")
	defer ms.end()

	if lifecycles, err := ha.db.SelectLifecycleIndex(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycles")
	} else {
		ms.send(w, lifecycles, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetLifecycle(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetLifecycle")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostLifecycle(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostLifecycle")
	defer ms.end()
	defer r.Body.Close()

	var l types.Lifecycle

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &l); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err = ha.db.InsertLifecycle(r.Context(), l, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert lifecycle")
	} else {
		ms.send(w, l, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchLifecycle(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchLifecycle")
	defer ms.end()
	defer r.Body.Close()

	var l types.Lifecycle

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &l); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if l, err = ha.db.UpdateLifecycle(r.Context(), l, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update lifecycle")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteLifecycle(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteLifecycle")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteLifecycle(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete lifecycle")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
