package huautla

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetAllStrains(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllStrains")
	defer ms.end()

	if Strains, err := ha.db.SelectAllStrains(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strains")
	} else {
		ms.send(w, Strains, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetStrain")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostStrain")
	defer ms.end()
	defer r.Body.Close()

	var s types.Strain

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if s, err = ha.db.InsertStrain(r.Context(), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert strain")
	} else {
		ms.send(w, s, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchStrain")
	defer ms.end()
	defer r.Body.Close()

	var s types.Strain

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateStrain(r.Context(), types.UUID(id), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update strain")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) DeleteStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteStrain")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete strain")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) GetGeneratedStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetGeneratedStrains")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.GeneratedStrain(r.Context(), types.UUID(id), ms.cid); err == sql.ErrNoRows {
		ms.send(w, nil, http.StatusNoContent)
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch generations")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PatchGeneratedStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchGeneratedStrains")
	defer ms.end()

	if gid := chi.URLParam(r, "gid"); gid == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if gid, err := url.QueryUnescape(gid); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else {
		ha.updateGeneratedStrain(w, r, (*types.UUID)(&gid), ms)
	}
}

func (ha *HuautlaAdaptor) DeleteGeneratedStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteGeneratedStrains")
	defer ms.end()

	ha.updateGeneratedStrain(w, r, nil, ms)
}

func (ha *HuautlaAdaptor) updateGeneratedStrain(w http.ResponseWriter, r *http.Request, gid *types.UUID, ms *methodStats) {
	if sid := chi.URLParam(r, "sid"); sid == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if sid, err := url.QueryUnescape(sid); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.UpdateGeneratedStrain(r.Context(), gid, types.UUID(sid), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update generation")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
