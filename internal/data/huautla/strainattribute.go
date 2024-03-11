package huautla

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetStrainAttributeNames(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetStrainAttributeNames")
	defer ms.end()

	if result, err := ha.db.KnownAttributeNames(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch attribute names")
	} else {
		ms.send(w, result, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostStrainAttribute(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostStrainAttribute")
	defer ms.end()
	defer r.Body.Close()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required strain id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if name := chi.URLParam(r, "at_name"); name == "" {
		ms.error(w, fmt.Errorf("missing required attribute name parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if value := chi.URLParam(r, "at_value"); value == "" {
		ms.error(w, fmt.Errorf("missing required attribute value parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.AddAttribute(r.Context(), &s, name, value, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add strainattribute")
	} else {
		ms.send(w, s, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchStrainAttribute(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchStrainAttribute")
	defer ms.end()
	defer r.Body.Close()

	if stID := chi.URLParam(r, "st_id"); stID == "" {
		ms.error(w, fmt.Errorf("missing required strain id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if name := chi.URLParam(r, "at_name"); name == "" {
		ms.error(w, fmt.Errorf("missing required attribute name parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if value := chi.URLParam(r, "at_value"); value == "" {
		ms.error(w, fmt.Errorf("missing required attribute value parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(stID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.ChangeAttribute(r.Context(), &s, name, value, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change strainattribute")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteStrainAttribute(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteStrainAttribute")
	defer ms.end()

	if stID := chi.URLParam(r, "st_id"); stID == "" {
		ms.error(w, fmt.Errorf("missing required strain id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if atID := chi.URLParam(r, "at_id"); atID == "" {
		ms.error(w, fmt.Errorf("missing required attribute id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(stID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.RemoveAttribute(r.Context(), &s, types.UUID(atID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove strainattribute")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}
