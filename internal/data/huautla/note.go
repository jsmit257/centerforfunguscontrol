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

func (ha *HuautlaAdaptor) getNotes(w http.ResponseWriter, r *http.Request) (oID string, notes []types.Note, err error) {
	ms := ha.start("GetNotes")
	defer ms.end()

	if oID = chi.URLParam(r, "o_id"); oID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if oID, err = url.QueryUnescape(oID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if notes, err = ha.db.GetNotes(r.Context(), types.UUID(oID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch notes")
	}
	return oID, notes, err
}

func (ha *HuautlaAdaptor) PostNote(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostNote")
	defer ms.end()
	defer r.Body.Close()

	var n types.Note

	if oID, notes, err := ha.getNotes(w, r); err != nil {
		return
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &n); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if notes, err = ha.db.AddNote(r.Context(), types.UUID(oID), notes, n, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add note")
	} else {
		ms.send(w, notes, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PatchNote(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchNote")
	defer ms.end()
	defer r.Body.Close()

	var n types.Note

	if _, notes, err := ha.getNotes(w, r); err != nil {
		return
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &n); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if notes, err = ha.db.ChangeNote(r.Context(), notes, n, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change note")
	} else {
		ms.send(w, notes, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeleteNote(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteNote")
	defer ms.end()

	if _, notes, err := ha.getNotes(w, r); err != nil {
		return
	} else if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if notes, err = ha.db.RemoveNote(r.Context(), notes, types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove note")
	} else {
		ms.send(w, notes, http.StatusOK)
	}
}
