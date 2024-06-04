package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

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

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					getErr: v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_GetNotes", "case": k}),
			mtrcs: nil,
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
				bytes.NewReader(nil))

			ha.GetNotes(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PostNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		p      *types.Note
		getErr error
		updErr error
		sc     int
	}{
		"happy_path": {
			id: "happy path",
			p:  &types.Note{},
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
		"missing_body": {
			id: "missing body",
			sc: http.StatusBadRequest,
		},
		"post_error": {
			id:     "post error",
			p:      &types.Note{},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					addErr: v.updErr,
					getErr: v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostNote", "case": k}),
			mtrcs: nil,
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
				bytes.NewReader(serializeNote(v.p)))

			ha.PostNote(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_ChangeNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		p      *types.Note
		getErr error
		updErr error
		sc     int
	}{
		"happy_path": {
			id: "happy path",
			p:  &types.Note{},
			sc: http.StatusOK,
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
		"patch_error": {
			id:     "post error",
			p:      &types.Note{},
			updErr: fmt.Errorf("some error"),
			sc:     http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Noter: &noterMock{
					changeErr: v.updErr,
					getErr:    v.getErr,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_ChangeNote", "case": k}),
			mtrcs: nil,
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
				bytes.NewReader(serializeNote(v.p)))

			ha.PatchNote(w, r)

			require.Equal(t, v.sc, w.Code)
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
			log:   log.WithFields(log.Fields{"test": "Test_DeleteNote", "case": k}),
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
