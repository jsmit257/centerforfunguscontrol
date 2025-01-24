package huautla

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) writePhoto(r *http.Request) (string, error) {
	var err error
	var data []byte
	var ct string

	if err = r.ParseMultipartForm(1 << 23); err != nil {
		return "", err
	} else if f, fh, err := r.FormFile("file"); err != nil {
		return "", err
	} else if data, err = io.ReadAll(f); err != nil {
		return "", err
	} else if len(data) < 4 {
		return "", fmt.Errorf("invalid request body")
	} else {
		ct = fh.Header.Get("Content-Type")
	}

	filetype := map[[4]byte]string{
		{0x89, 'P', 'N', 'G'}:    "image/png",
		{0xff, 0xd8, 0xff, 0xe0}: "image/jpeg",
		{0xff, 0xd8, 0xff, 0xe1}: "image/jpeg", // old format??
		{'G', 'I', 'F', '8'}:     "image/gif",
		{'M', 'M', 0, '*'}:       "image/tiff",
		{'I', 'I', '*', 0}:       "image/tiff",
	}[[4]byte(data[:4])]
	if filetype == "" {
		filetype = append(r.Header[http.CanonicalHeaderKey("Content-Type")], "image/x-unknown")[0]
	}

	// r.Context().Value(metrics.Log).(*logrus.Entry).WithFields(log.Fields{
	logrus.WithFields(log.Fields{
		"from-request": ct,
		"from-app":     filetype,
	}).
		Warn("comparing types")

	ext := map[string]string{
		"image/jpeg": "jpg",
		"image/jpg":  "jpg",
		"image/png":  "png",
		"image/gif":  "gif",
		"image/tiff": "tiff",
	}[filetype]
	if ext == "" {
		ext = "unk"
	}

	name := fmt.Sprintf("%s.%s", uuid.New().String(), ext)

	return name, ha.filer("album/"+name, data, 0644)
}

func (ha *HuautlaAdaptor) GetPhotos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()
	ms := ha.start(ctx, "GetPhotos")

	if _, photos, err := ha.getPhotos(w, r, ms); err != nil {
		return
	} else {
		ms.send(w, http.StatusOK, photos)
	}
}

func (ha *HuautlaAdaptor) getPhotos(w http.ResponseWriter, r *http.Request, ms *methodStats) (olID string, photos []types.Photo, err error) {

	if olID = chi.URLParam(r, "o_id"); olID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if olID, err = url.QueryUnescape(olID); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if photos, err = ha.db.GetPhotos(r.Context(), types.UUID(olID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch photos")
	}
	return olID, photos, err
}

func (ha *HuautlaAdaptor) PostPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostPhoto")
	defer r.Body.Close()

	var p types.Photo

	if oID, photos, err := ha.getPhotos(w, r, ms); err != nil {
		return
	} else if p.Filename, err = ha.writePhoto(r); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read/write request body")
	} else if photos, err = ha.db.AddPhoto(r.Context(), types.UUID(oID), photos, p, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add photo")
	} else {
		ms.send(w, http.StatusOK, photos)
	}
}

func (ha *HuautlaAdaptor) PatchPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchPhoto")
	defer r.Body.Close()

	var p types.Photo

	if _, photos, err := ha.getPhotos(w, r, ms); err != nil {
		return
	} else if p.UUID = types.UUID(chi.URLParam(r, "id")); p.UUID == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if p.Filename, err = ha.writePhoto(r); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if photos, err = ha.db.ChangePhoto(r.Context(), photos, p, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change photo")
	} else {
		ms.send(w, http.StatusOK, photos)
	}
}

func (ha *HuautlaAdaptor) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeletePhoto")

	if _, photos, err := ha.getPhotos(w, r, ms); err != nil {
		return
	} else if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if photos, err = ha.db.RemovePhoto(r.Context(), photos, types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove photo")
	} else {
		ms.send(w, http.StatusOK, photos)
	}
}
