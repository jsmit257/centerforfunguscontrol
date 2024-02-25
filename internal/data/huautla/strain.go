package huautla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetAllStrains(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllStrains")
	defer ms.end()

	if Strains, err := ha.db.SelectAllStrains(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch Strains")
	} else {
		ms.send(w, Strains, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetStrain")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch Strain")
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
		ms.error(w, err, http.StatusInternalServerError, "failed to insert Strain")
	}

	ms.send(w, s, http.StatusOK)
}

func (ha *HuautlaAdaptor) PatchStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchStrain")
	defer ms.end()
	defer r.Body.Close()

	var s types.Strain

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if err := json.Unmarshal(body, &s); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateStrain(r.Context(), types.UUID(id), s, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update Strain")
	}

	ms.send(w, nil, http.StatusNoContent)
}

func (ha *HuautlaAdaptor) DeleteStrain(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteStrain")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if err := ha.db.DeleteStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete Strain")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
