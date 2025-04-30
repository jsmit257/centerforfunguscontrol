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

func fmtSource(s types.Source) string {
	result, _ := json.Marshal(&s)
	return string(result)
}

var origins = map[string]struct{}{
	"event":  {},
	"strain": {},
}

func (ha *HuautlaAdaptor) PostSource(w http.ResponseWriter, r *http.Request) {
	ms := ha.start(r.Context(), "PostSource")
	defer r.Body.Close()

	var s types.Source

	if genID, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if origin := chi.URLParam(r, "origin"); origin == "" {
		ms.error(w, fmt.Errorf("missing required parameter: origin"), http.StatusBadRequest, "missing required parameter")
	} else if origin, err = url.QueryUnescape(origin); err != nil {
		ms.error(w, fmt.Errorf("malformed parameter: origin"), http.StatusBadRequest, "malformed parameter")
	} else if _, ok := origins[origin]; !ok {
		ms.error(w, fmt.Errorf("origin value not allowed: %s", origin), http.StatusBadRequest, "origin value not allowed")

		// } else if err := bodyHelper(r, s); err != nil {
		// 	ms.error(w, err, http.StatusBadRequest, "couldn't read request body")

	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")

	} else if s, err := ha.db.InsertSource(r.Context(), genID, origin, s, ms.cid); err != nil {
		ms.error(w, fmt.Errorf("%w: %s", err, fmtSource(s)), http.StatusInternalServerError, err)
	} else {
		ms.created(w, s)
	}
}

func (ha *HuautlaAdaptor) PatchSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchSource")
	defer r.Body.Close()

	var s types.Source

	if _, err := getUUIDByName("g_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if origin := chi.URLParam(r, "origin"); origin == "" {
		ms.error(w, fmt.Errorf("missing required parameter: origin"), http.StatusBadRequest, "missing required parameter")
	} else if origin, err = url.QueryUnescape(origin); err != nil {
		ms.error(w, fmt.Errorf("malformed parameter: origin"), http.StatusBadRequest, "malformed parameter")
	} else if _, ok := origins[origin]; !ok {
		ms.error(w, fmt.Errorf("origin value not allowed: %s", origin), http.StatusBadRequest, "origin value not allowed")
	} else if s.UUID, err = getUUIDByName("s_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: source id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if err := ha.db.UpdateSource(r.Context(), origin, s, ms.cid); err != nil {
		ms.error(w, fmt.Errorf("%w: %s", err, fmtSource(s)), http.StatusInternalServerError, err)
	} else {
		ms.empty(w)
	}
}

func (ha *HuautlaAdaptor) DeleteSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteSource")

	if genID := chi.URLParam(r, "g_id"); genID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(genID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if sID := chi.URLParam(r, "s_id"); sID == "" {
		ms.error(w, fmt.Errorf("missing required event id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if sID, err := url.QueryUnescape(sID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.RemoveSource(r.Context(), &g, types.UUID(sID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove source")
	} else {
		ms.ok(w, g)
	}
}
