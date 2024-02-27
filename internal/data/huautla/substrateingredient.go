package huautla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) PostSubstrateIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostSubstrateIngredient")
	defer ms.end()
	defer r.Body.Close()

	var i types.Ingredient

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &i); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if s, err := ha.db.SelectSubstrate(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrate")
	} else if err = ha.db.AddIngredient(r.Context(), &s, i, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add substrateingredient")
	} else {
		ms.send(w, s, http.StatusOK)
	}

}

func (ha *HuautlaAdaptor) PatchSubstrateIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostSubstrateIngredient")
	defer ms.end()
	defer r.Body.Close()

	var newI types.Ingredient

	if suID := chi.URLParam(r, "su_id"); suID == "" {
		ms.error(w, fmt.Errorf("missing required substrate id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if igID := chi.URLParam(r, "ig_id"); igID == "" {
		ms.error(w, fmt.Errorf("missing required ingredient id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &newI); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if s, err := ha.db.SelectSubstrate(r.Context(), types.UUID(suID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrate")
	} else if err = ha.db.ChangeIngredient(r.Context(), &s, types.Ingredient{UUID: types.UUID(igID)}, newI, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change substrateingredient")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteSubstrateIngredient(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostSubstrateIngredient")
	defer ms.end()

	// 	RemoveIngredient(ctx context.Context, s *Substrate, i Ingredient, cid CID) error
	if suID := chi.URLParam(r, "su_id"); suID == "" {
		ms.error(w, fmt.Errorf("missing required substrate id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if igID := chi.URLParam(r, "ig_id"); igID == "" {
		ms.error(w, fmt.Errorf("missing required ingredient id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if s, err := ha.db.SelectSubstrate(r.Context(), types.UUID(suID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch substrate")
	} else if err = ha.db.RemoveIngredient(r.Context(), &s, types.Ingredient{UUID: types.UUID(igID)}, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove substrateingredient")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}
