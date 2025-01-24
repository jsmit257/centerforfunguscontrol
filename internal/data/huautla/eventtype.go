package huautla

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetAllEventTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetAllEventTypes")

	if stages, err := ha.db.SelectAllEventTypes(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch eventtypes")
	} else {
		ms.send(w, http.StatusOK, stages)
	}
}

func (ha *HuautlaAdaptor) GetEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetEventType")

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if eventtype, err := ha.db.SelectEventType(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch eventtype")
	} else {
		ms.send(w, http.StatusOK, eventtype)
	}
}

func (ha *HuautlaAdaptor) PostEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostEventType")
	defer r.Body.Close()

	var et types.EventType

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &et); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if et, err = ha.db.InsertEventType(r.Context(), et, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert eventtype")
	} else {
		ms.send(w, http.StatusCreated, et)
	}
}

func (ha *HuautlaAdaptor) PatchEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchEventType")
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
		ms.send(w, http.StatusNoContent, nil)
	}
}

func (ha *HuautlaAdaptor) DeleteEventType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteEventType")

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteEventType(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete eventtype")
	} else {
		ms.send(w, http.StatusNoContent, nil)
	}
}

func (ha *HuautlaAdaptor) GetEventTypeReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetEventTypeReport")

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if v, err := ha.db.EventTypeReport(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch eventtype")
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch eventtype")
	} else {
		ms.send(w, http.StatusOK, v)
	}
}
