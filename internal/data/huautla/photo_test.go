package huautla

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

type photoerMock struct {
	getResult,
	addResult,
	changeResult,
	rmResult []types.Photo

	getErr,
	addErr,
	changeErr,
	rmErr error
}

func Test_PostPhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id       types.UUID
		data     []byte
		getErr   error
		updErr   error
		writeErr error
		sc       int
	}{
		"happy_path": {
			id:   "happy path",
			data: []byte{0x89, 0x50, 0x4e, 0x47},
			sc:   http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"get_error": {
			id:     "get error",
			getErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_body": {
			id: "missing body",
			sc: http.StatusBadRequest,
		},
		"write_error": {
			id:       "write_error",
			data:     []byte{0xff, 0xd8, 0xff, 0xe0},
			writeErr: fmt.Errorf("some error"),
			sc:       http.StatusBadRequest,
		},
		"post_error": {
			id:     "post error",
			data:   []byte{0x00, 0x00, 0x00, 0x00},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Photoer: &photoerMock{
					addErr: v.updErr,
					getErr: v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostPhoto", "case": k}),
			mtrcs: nil,
			filer: func(string, []byte, fs.FileMode) error {
				return v.writeErr
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"o_id"}, Values: []string{string(v.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(v.data))

			ha.PostPhoto(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchPhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id, oID  types.UUID
		data     []byte
		getErr   error
		updErr   error
		writeErr error
		sc       int
	}{
		"happy_path": {
			oID:  "happy path",
			id:   "happy path",
			data: []byte{0x89, 0x50, 0x4e, 0x47},
			sc:   http.StatusOK,
		},
		"missing_photo_id": {
			oID: "missing_photo_id",
			sc:  http.StatusBadRequest,
		},
		"get_error": {
			oID:    "get error",
			id:     "get error",
			getErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_body": {
			oID: "missing body",
			id:  "missing body",
			sc:  http.StatusBadRequest,
		},
		"write_error": {
			oID:      "write_error",
			id:       "write_error",
			data:     []byte{0x47, 0x49, 0x46, 0x38},
			writeErr: fmt.Errorf("some error"),
			sc:       http.StatusBadRequest,
		},
		"patch_error": {
			oID:    "post error",
			id:     "post error",
			data:   []byte{0x4d, 0x4d, 0x00, 0x2a},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Photoer: &photoerMock{
					changeErr: v.updErr,
					getErr:    v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchPhoto", "case": k}),
			mtrcs: nil,
			filer: func(string, []byte, fs.FileMode) error {
				return v.writeErr
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"o_id", "id"}, Values: []string{
				string(v.oID),
				string(v.id),
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(v.data))

			ha.PatchPhoto(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeletePhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id, oID types.UUID
		getErr  error
		updErr  error
		sc      int
	}{
		"happy_path": {
			oID: "happy path",
			id:  "happy path",
			sc:  http.StatusOK,
		},
		"missing_id": {
			oID: "happy path",
			sc:  http.StatusBadRequest,
		},
		"urldecode_error": {
			oID: "urldecode_error",
			id:  "%zzz",
			sc:  http.StatusBadRequest,
		},
		"get_error": {
			oID:    "get error",
			getErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
		"patch_error": {
			oID:    "post error",
			id:     "post error",
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Photoer: &photoerMock{
					rmErr:  v.updErr,
					getErr: v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeletePhoto", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"o_id", "id"}, Values: []string{
				string(v.oID),
				string(v.id),
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(nil))

			ha.DeletePhoto(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func (pm *photoerMock) GetPhotos(context.Context, types.UUID, types.CID) ([]types.Photo, error) {
	return pm.getResult, pm.getErr
}

func (pm *photoerMock) AddPhoto(context.Context, types.UUID, []types.Photo, types.Photo, types.CID) ([]types.Photo, error) {
	return pm.addResult, pm.addErr
}

func (pm *photoerMock) ChangePhoto(context.Context, []types.Photo, types.Photo, types.CID) ([]types.Photo, error) {
	return pm.changeResult, pm.changeErr
}

func (pm *photoerMock) RemovePhoto(context.Context, []types.Photo, types.UUID, types.CID) ([]types.Photo, error) {
	return pm.rmResult, pm.rmErr
}
