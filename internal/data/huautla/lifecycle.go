package huautla

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetLifecycleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetAllLifecycles")

	if lifecycles, err := ha.db.SelectLifecycleIndex(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycles")
	} else {
		ms.send(w, http.StatusOK, lifecycles)
	}
}

func (ha *HuautlaAdaptor) GetLifecycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetLifecycle")

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch lifecycle")
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else {
		ms.send(w, http.StatusOK, l)
	}
}

func (ha *HuautlaAdaptor) PostLifecycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostLifecycle")
	defer r.Body.Close()

	var l types.Lifecycle

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &l); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err = ha.db.InsertLifecycle(r.Context(), l, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert lifecycle")
	} else {
		ms.send(w, http.StatusCreated, l)
	}
}

func (ha *HuautlaAdaptor) PatchLifecycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchLifecycle")
	defer r.Body.Close()

	var l types.Lifecycle

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &l); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if l, err = ha.db.UpdateLifecycle(r.Context(), l, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, err.Error())
	} else {
		ms.send(w, http.StatusOK, l)
	}
}

func (ha *HuautlaAdaptor) DeleteLifecycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteLifecycle")

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if err = ha.db.DeleteLifecycle(r.Context(), id, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete lifecycle")
	} else {
		ms.send(w, http.StatusNoContent, nil)
	}
}

func (ha *HuautlaAdaptor) GetLifecycleReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "GetLifecycleReport")

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if l, err := ha.db.LifecycleReport(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch lifecycle")
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else {
		ms.send(w, http.StatusOK, l)
	}
}
