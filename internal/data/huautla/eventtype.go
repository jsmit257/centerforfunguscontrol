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

func (ha *HuautlaAdaptor) GetAllEventTypes(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllEventTypes")
	defer ms.end()

	if stages, err := ha.db.SelectAllEventTypes(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch eventtypes")
	} else {
		ms.send(w, stages, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetEventType(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetEventType")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if eventtype, err := ha.db.SelectEventType(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch eventtype")
	} else {
		ms.send(w, eventtype, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostEventType(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostEventType")
	defer ms.end()
	defer r.Body.Close()

	var et types.EventType

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &et); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if et, err = ha.db.InsertEventType(r.Context(), et, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert eventtype")
	} else {
		ms.send(w, et, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchEventType(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchEventType")
	defer ms.end()
	defer r.Body.Close()

	var et types.EventType

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal([]byte(body), &et); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if err = ha.db.UpdateEventType(r.Context(), types.UUID(id), et, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update eventtype")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) DeleteEventType(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteEventType")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteEventType(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete eventtype")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
