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

func (ha *HuautlaAdaptor) GetAllIngredients(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllIngredients")
	defer ms.end()

	if Ingredients, err := ha.db.SelectAllIngredients(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch ingredients")
	} else {
		ms.send(w, Ingredients, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetIngredient")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SelectIngredient(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch ingredient")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostIngredient")
	defer ms.end()
	defer r.Body.Close()

	var i types.Ingredient

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &i); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if i, err = ha.db.InsertIngredient(r.Context(), i, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert ingredient")
	} else {
		ms.send(w, i, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchIngredient")
	defer ms.end()
	defer r.Body.Close()

	var i types.Ingredient

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &i); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateIngredient(r.Context(), types.UUID(id), i, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update ingredient")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteIngredient")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteIngredient(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete ingredient")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}
