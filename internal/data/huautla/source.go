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

func (ha *HuautlaAdaptor) PostSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostSource")
	defer r.Body.Close()

	var s types.Source

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.AddSource(r.Context(), &g, s, ms.cid); err != nil {
		// ms.error(w, err, http.StatusInternalServerError, "failed to add source")
		ms.error(w, fmt.Errorf("%w: %s", err, fmtSource(s)), http.StatusInternalServerError, err)
	} else {
		ms.send(w, http.StatusCreated, g)
	}
}

func (ha *HuautlaAdaptor) PostStrainSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostStrainSource")
	defer r.Body.Close()

	var s types.Source

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.AddStrainSource(r.Context(), &g, s, ms.cid); err != nil {
		// ms.error(w, err, http.StatusInternalServerError, "failed to add source")
		ms.error(w, err, http.StatusInternalServerError, func(s types.Source) string {
			result, _ := json.Marshal(&s)
			return string(result)
		}(s))
	} else {
		ms.send(w, http.StatusCreated, g)
	}
}

func (ha *HuautlaAdaptor) PostEventSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostEventSource")
	defer r.Body.Close()

	var e types.Event

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.AddEventSource(r.Context(), &g, e, ms.cid); err != nil {
		// ms.error(w, err, http.StatusInternalServerError, "failed to add source")
		ms.error(w, err, http.StatusInternalServerError, func(e types.Event) string {
			result, _ := json.Marshal(&e)
			return string(result)
		}(e))
	} else {
		ms.send(w, http.StatusCreated, g)
	}
}

func (ha *HuautlaAdaptor) PatchSourceNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchSourceNew")
	defer r.Body.Close()

	var s types.Source

	if gID, err := getUUIDByName("g_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if s.UUID, err = getUUIDByName("s_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: source id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.ChangeSource(r.Context(), &g, s, ms.cid); err != nil {
		ms.error(w, fmt.Errorf("%w: %s", err, fmtSource(s)), http.StatusInternalServerError, err)
	} else {
		ms.send(w, http.StatusOK, g)
	}
}

func (ha *HuautlaAdaptor) PatchSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchSource")
	defer r.Body.Close()

	var s types.Source

	if gID := chi.URLParam(r, "id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required parameter"), http.StatusBadRequest, "missing required parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
		ms.error(w, fmt.Errorf("malformed parameter"), http.StatusBadRequest, "malformed parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generation")
	} else if err := ha.db.ChangeSource(r.Context(), &g, s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, func(s types.Source) string {
			result, _ := json.Marshal(&s)
			return string(result)
		}(s))
	} else {
		ms.send(w, http.StatusOK, g)
	}
}

func (ha *HuautlaAdaptor) DeleteSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteSource")

	if gID := chi.URLParam(r, "g_id"); gID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gID, err := url.QueryUnescape(gID); err != nil {
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
		ms.send(w, http.StatusOK, g)
	}
}
