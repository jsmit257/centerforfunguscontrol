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

func (ha *HuautlaAdaptor) PatchTS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchTS")

	var patch types.Timestamp
	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed parameter")
	} else if table := chi.URLParam(r, "table"); table == "" {
		ms.error(w, fmt.Errorf("missing required table parameter"), http.StatusBadRequest, "missing required parameter")
	} else if table, err := url.QueryUnescape(table); err != nil {
		ms.error(w, fmt.Errorf("malformed table parameter"), http.StatusBadRequest, "malformed parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to read request body")
	} else if err := json.Unmarshal(body, &patch); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to parse request body")
	} else if err = patch.Validate(); err != nil {
		ms.error(w, err, http.StatusBadRequest, "timestamp validation failed")
	} else if err := ha.db.UpdateTimestamps(r.Context(), table, types.UUID(id), patch); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete vendor")
	} else {
		ms.send(w, http.StatusNoContent, nil)
	}
}

func (ha *HuautlaAdaptor) Undel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "Undel")

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed parameter")
	} else if table := chi.URLParam(r, "table"); table == "" {
		ms.error(w, fmt.Errorf("missing required table parameter"), http.StatusBadRequest, "missing required parameter")
	} else if table, err := url.QueryUnescape(table); err != nil {
		ms.error(w, fmt.Errorf("malformed table parameter"), http.StatusBadRequest, "malformed parameter")
	} else if err := ha.db.Undelete(r.Context(), table, types.UUID(id)); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to undelete row")
	} else {
		ms.send(w, http.StatusNoContent, nil)
	}
}
