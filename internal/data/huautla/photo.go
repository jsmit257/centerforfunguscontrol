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

func (ha *HuautlaAdaptor) getPhotos(w http.ResponseWriter, r *http.Request) (oID string, photos []types.Photo, err error) {
	ms := ha.start("GetPhotos")
	defer ms.end()

	if oID = chi.URLParam(r, "o_id"); oID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if oID, err = url.QueryUnescape(oID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if photos, err = ha.db.GetPhotos(r.Context(), types.UUID(oID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch photos")
	}
	return oID, photos, err
}

func (ha *HuautlaAdaptor) PostPhoto(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostPhoto")
	defer ms.end()
	defer r.Body.Close()

	var p types.Photo

	if oID, photos, err := ha.getPhotos(w, r); err != nil {
		return
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &p); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if photos, err = ha.db.AddPhoto(r.Context(), types.UUID(oID), photos, p, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add photo")
	} else {
		ms.send(w, photos, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PatchPhoto(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchPhoto")
	defer ms.end()
	defer r.Body.Close()

	var p types.Photo

	if _, photos, err := ha.getPhotos(w, r); err != nil {
		return
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &p); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if photos, err = ha.db.ChangePhoto(r.Context(), photos, p, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change photo")
	} else {
		ms.send(w, photos, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeletePhoto")
	defer ms.end()

	if _, photos, err := ha.getPhotos(w, r); err != nil {
		return
	} else if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if photos, err = ha.db.RemovePhoto(r.Context(), photos, types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove photo")
	} else {
		ms.send(w, photos, http.StatusOK)
	}
}
