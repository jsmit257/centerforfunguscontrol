package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
)

type noterMock struct {
	getResult,
	addResult,
	changeResult,
	rmResult []types.Note

	getErr,
	addErr,
	changeErr,
	rmErr error
}

func Test_GetNotes(t *testing.T) {
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					getErr: tc.getErr,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"o_id"}, Values: []string{string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(nil))

			ha.GetNotes(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_PostNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		note   *types.Note
		getErr error
		updErr error
		sc     int
	}{
		"happy_path": {
			id:   "happy path",
			note: &types.Note{},
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
		"read_fails": {
			id: "missing body",
			sc: http.StatusBadRequest,
		},
		"missing_body": {
			id: "missing body",
			sc: http.StatusBadRequest,
		},
		"post_error": {
			id:     "post error",
			note:   &types.Note{},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					addErr: tc.updErr,
					getErr: tc.getErr,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeNote(tc.note)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"o_id"}, Values: []string{string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bodyreader)

			ha.PostNote(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_ChangeNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		noteID types.UUID
		note   *types.Note
		getErr error
		updErr error
		sc     int
	}{
		"happy_path": {
			id:     "happy path",
			noteID: "happy path",
			note:   &types.Note{},
			sc:     http.StatusOK,
		},
		"get_error": {
			id:     "get error",
			getErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
		"noteid_error": {
			id:     "noteid error",
			noteID: "%zzz",
			sc:     http.StatusBadRequest,
		},
		"read_fails": {
			id:     "read fails",
			noteID: "read fails",
			sc:     http.StatusBadRequest,
		},
		"missing_body": {
			id:     "missing body",
			noteID: "missing body",
			sc:     http.StatusBadRequest,
		},
		"patch_error": {
			id:     "patch_error",
			noteID: "patch_error",
			note:   &types.Note{},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					changeErr: tc.updErr,
					getErr:    tc.getErr,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeNote(tc.note)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{
				"o_id",
				"n_id",
			}, Values: []string{
				string(tc.id),
				string(tc.noteID),
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bodyreader)

			ha.PatchNote(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_DeleteNote(t *testing.T) {
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
			oID: "happy path",
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
				Noter: &noterMock{
					rmErr:  v.updErr,
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
				bytes.NewReader(serializeNote(nil)))

			ha.DeleteNote(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func serializeNote(p *types.Note) []byte {
	if p == nil {
		return []byte{}
	}
	result, _ := json.Marshal(p)
	return result
}

func (nm *noterMock) GetNotes(context.Context, types.UUID, types.CID) ([]types.Note, error) {
	return nm.getResult, nm.getErr
}

func (nm *noterMock) AddNote(context.Context, types.UUID, []types.Note, types.Note, types.CID) ([]types.Note, error) {
	return nm.addResult, nm.addErr
}

func (nm *noterMock) ChangeNote(context.Context, []types.Note, types.Note, types.CID) ([]types.Note, error) {
	return nm.changeResult, nm.changeErr
}

func (nm *noterMock) RemoveNote(context.Context, []types.Note, types.UUID, types.CID) ([]types.Note, error) {
	return nm.rmResult, nm.rmErr
}
