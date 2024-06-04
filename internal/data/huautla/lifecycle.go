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

func (ha *HuautlaAdaptor) GetLifecycleIndex(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllLifecycles")
	defer ms.end()

	if lifecycles, err := ha.db.SelectLifecycleIndex(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycles")
	} else {
		ms.send(w, lifecycles, http.StatusOK)
	}
}

func getUUIDByName(name string, w http.ResponseWriter, r *http.Request, ms *methodStats) (uuid types.UUID, err error) {
	if id := chi.URLParam(r, name); id == "" {
		err = fmt.Errorf("missing required id parameter")
	} else if id, err = url.QueryUnescape(id); err != nil {
		err = fmt.Errorf("malformed id parameter")
	} else {
		uuid = types.UUID(id)
	}
	return uuid, err
}

func (ha *HuautlaAdaptor) GetLifecyclesByAttrs(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetLifecyclesByAttrs")
	defer ms.end()

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		ms.error(w, err, http.StatusBadRequest, "query string is malformed")
		return
	}

	p := types.ReportAttrs{}
	for k, v := range q {
		if v[0] == "" {
			ms.l.Errorf("query value is empty for: %s", k)
		} else {
			p[k] = types.UUID(v[0])
		}
	}

	if len(p) == 0 {
		ms.error(w, fmt.Errorf("on report parameters supplied"), http.StatusBadRequest, "on report parameters supplied")
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
