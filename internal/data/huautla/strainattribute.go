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

	a := types.StrainAttribute{}

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required strain id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal([]byte(body), &a); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if a.Name == "" {
		ms.error(w, fmt.Errorf("incomplete strainattribute body"), http.StatusBadRequest, "incomplete strainattribute body")
	} else if a.Value == "" {
		ms.error(w, fmt.Errorf("incomplete strainattribute body"), http.StatusBadRequest, "incomplete strainattribute body")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if a, err := ha.db.AddAttribute(r.Context(), &s, a, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add strainattribute")
	} else {
		ms.send(w, a, http.StatusCreated)
	}
}

func (ha *HuautlaAdaptor) PatchStrainAttribute(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchStrainAttribute")
	defer ms.end()
	defer r.Body.Close()

	a := types.StrainAttribute{}

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required strain id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal([]byte(body), &a); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if a.Name == "" {
		ms.error(w, fmt.Errorf("incomplete strainattribute body"), http.StatusBadRequest, "incomplete strainattribute body")
	} else if a.Value == "" {
		ms.error(w, fmt.Errorf("incomplete strainattribute body"), http.StatusBadRequest, "incomplete strainattribute body")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.ChangeAttribute(r.Context(), &s, a, ms.cid); err != nil {
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
	} else if stID, err := url.QueryUnescape(stID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if atID := chi.URLParam(r, "at_id"); atID == "" {
		ms.error(w, fmt.Errorf("missing required attribute id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if atID, err := url.QueryUnescape(atID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if s, err := ha.db.SelectStrain(r.Context(), types.UUID(stID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch strain")
	} else if err := ha.db.RemoveAttribute(r.Context(), &s, types.UUID(atID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove strainattribute")
	} else {
		ms.send(w, s, http.StatusOK)
	}
}
