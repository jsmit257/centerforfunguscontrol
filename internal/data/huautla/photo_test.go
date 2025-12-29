package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
)

type photoerMock struct {
	getIndexResult,
	getResult,
	addResult,
	changeResult,
	rmResult []types.Photo

	getIndexErr,
	getErr,
	addErr,
	changeErr,
	ndxErr error
}

func Test_GetPhotoIndex(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		getIndexResult []types.Photo
		getIndexErr    error
		sc             int
	}{
		"happy_path": {
			getIndexResult: []types.Photo{
				{
					UUID:     "getPhotoIndex_0",
					Filename: "getPhotoIndex_file",
					MTime:    time.Time{},
					CTime:    time.Time{},
					Owner: &types.PhotoOwner{
						ParentType: "generation",
						OwnerUUID:  "owner_0",
						ParentUUID: func(uuid types.UUID) *types.UUID { return &uuid }("parent_0"),
						Label:      "owner_label",
					},
				},
			},
			sc: http.StatusOK,
		},
		"get_error": {
			getIndexErr: fmt.Errorf("some error"),
			sc:          http.StatusInternalServerError,
		},
	}

	for name, rc := range set {
		name, rc := name, rc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodPost,
				"url",
				nil)

			(&HuautlaAdaptor{
				db: &huautlaMock{
					Photoer: &photoerMock{
						getIndexResult: rc.getIndexResult,
						getIndexErr:    rc.getIndexErr,
					},
				},
			}).GetPhotosIndex(w, r)

			require.Equal(t, rc.sc, w.Code)
			if w.Code != http.StatusOK {
				return
			}

			body, err := io.ReadAll(w.Body)
			require.Nil(t, err, "reading response body")
			result := []types.Photo{}
			err = json.Unmarshal(body, &result)
			require.Nil(t, err, "unmarshalling response body: %s", body)
			require.Equal(t, rc.getIndexResult, result)
		})
	}
}

func photoHelper(d []byte) (io.Reader, string) {
	var b = &bytes.Buffer{}
	var w = multipart.NewWriter(b)

	fw, _ := w.CreateFormFile("file", "sample.png")
	_, _ = io.Copy(fw, bytes.NewReader(d))
	w.Close()

	return b, w.FormDataContentType()
}

func sendPhoto(f func(w http.ResponseWriter, r *http.Request), data []byte, meth string, url chi.RouteParams) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	defer w.Result().Body.Close()
	rctx := chi.NewRouteContext()
	rctx.URLParams = url
	reader, contentType := photoHelper(data)
	r, _ := http.NewRequestWithContext(
		context.WithValue(
			metrics.MockServiceContext,
			chi.RouteCtxKey,
			rctx),
		meth,
		"url",
		reader)
	r.Header.Set("Content-Type", contentType)

	f(w, r)

	return w
}

func Test_getFormat(t *testing.T) {
	for frmt, head := range formats {
		result := getFormat(head.magic)
		require.Equal(t, frmt, result, fmt.Sprintf("test_%s", frmt))
	}

	result := getFormat([]byte("RIFF****WEBPVP8X"))
	require.Equal(t, format("image/webp"), result, "test_sparse")

	result = getFormat([]byte("head.magic"))
	require.Equal(t, format("unknown"), result, "test_random")
}

func Test_GetPhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		getErr error
		sc     int
	}{
		"happy_path": {
			id: "happy path",
			sc: http.StatusOK,
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
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Photoer: &photoerMock{
					getErr: v.getErr,
				},
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
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				nil)

			ha.GetPhotos(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
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
			sc:   http.StatusCreated,
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
			filer: func(string, []byte, fs.FileMode) error {
				return v.writeErr
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			w := sendPhoto(
				ha.PostPhoto,
				v.data,
				http.MethodPost,
				chi.RouteParams{Keys: []string{"o_id"}, Values: []string{string(v.id)}})

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
			filer: func(string, []byte, fs.FileMode) error {
				return v.writeErr
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			w := sendPhoto(
				ha.PatchPhoto,
				v.data,
				http.MethodPatch,
				chi.RouteParams{Keys: []string{"o_id", "id"}, Values: []string{
					string(v.oID),
					string(v.id),
				}})

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
					ndxErr: v.updErr,
					getErr: v.getErr,
				},
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
					metrics.MockServiceContext,
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

func (pm *photoerMock) AllPhotos(context.Context, types.CID) ([]types.Photo, error) {
	return pm.getIndexResult, pm.getIndexErr
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
	return pm.rmResult, pm.ndxErr
}
