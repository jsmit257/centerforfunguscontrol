package huautla

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

func (ha *HuautlaAdaptor) GetLifecyclesByAttrs(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetLifecyclesByAttrs")
	defer ms.end()

	if q, err := url.ParseQuery(r.URL.RawQuery); err != nil {
		ms.error(w, err, http.StatusBadRequest, "query string is malformed")
	} else if p, err := types.NewReportAttrs(q); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't parse report params")
	} else if !p.Contains("lifecycle-id", "strain-id", "grain-id", "bulk-id") {
		ms.error(w, fmt.Errorf("no report parameters supplied"), http.StatusBadRequest, "no report parameters supplied")
	} else if lifecycles, err := ha.db.SelectLifecyclesByAttrs(r.Context(), p, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycles")
	} else {
		ms.send(w, lifecycles, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetLifecycle(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetLifecycle")
	defer ms.end()

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch lifecycle")
	} else if err != nil {
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

	// id := getID(w, r, ms)
	if body, err := io.ReadAll(r.Body); err != nil {
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

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if err = ha.db.DeleteLifecycle(r.Context(), id, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete lifecycle")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) GetLifecycleReport(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetLifecycleReport")
	defer ms.end()

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if l, err := ha.db.LifecycleReport(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch lifecycle")
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else {
		ms.send(w, l, http.StatusOK)
	}
}
