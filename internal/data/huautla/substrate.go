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

func (ha *HuautlaAdaptor) GetAllSubstrates(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllSubstrates")
	defer ms.end()

	if substrates, err := ha.db.SelectAllSubstrates(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrates")
	} else {
		ms.send(w, substrates, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetSubstrate(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetSubstrate")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SelectSubstrate(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrate")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostSubstrate(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostSubstrate")
	defer ms.end()
	defer r.Body.Close()

	var s types.Substrate

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if s, err = ha.db.InsertSubstrate(r.Context(), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert substrate")
	} else {
		ms.send(w, s, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchSubstrate(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchSubstrate")
	defer ms.end()
	defer r.Body.Close()

	var s types.Substrate

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateSubstrate(r.Context(), types.UUID(id), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update substrate")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) DeleteSubstrate(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteSubstrate")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteSubstrate(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete substrate")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) GetSubstrateReport(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetSubstrateReport")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SubstrateReport(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrate")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}
