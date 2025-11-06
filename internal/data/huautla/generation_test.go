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
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type generationerMock struct {
	all    []types.Generation
	allErr error

	sel    types.Generation
	selErr error

	ins    types.Generation
	insErr error

	upd    types.Generation
	updErr error

	str    types.Strain
	strErr error

	rmErr error

	patchErr error

	rpt    types.Entity
	rptErr error
}

func Test_GetGenerationIndex(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Generation
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Generation{},
			sc:     http.StatusOK,
		},
		"db_error": {
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					all:    tc.result,
					allErr: tc.err,
				},
			},
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				nil)
			ha.GetGenerationIndex(w, r)
			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Generation{}, &tc.result)
			}
		})
	}
}

func Test_GetGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Generation
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Generation{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"url_decode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					sel:    tc.result,
					selErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{tc.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				nil)

			ha.GetGeneration(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Generation{}, &tc.result)
			}
		})
	}
}

func Test_PostGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gen    *types.Generation
		result types.Generation
		err    error
		sc     int
	}{
		"happy_path": {
			gen:    &types.Generation{},
			result: types.Generation{},
			sc:     http.StatusCreated,
		},
		"read_fails": {
			sc: http.StatusBadRequest,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			gen: &types.Generation{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					ins:    tc.result,
					insErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeGeneration(tc.gen)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				bodyreader)

			ha.PostGeneration(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Generation{}, &tc.result)
			}
		})
	}
}

func Test_PatchGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		gen *types.Generation
		err error
		sc  int
	}{
		"happy_path": {
			id:  "1",
			gen: &types.Generation{},
			sc:  http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"malformed_id": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"read_fails": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"missing_stage": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			gen: &types.Generation{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					upd:    types.Generation{},
					updErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeGeneration(tc.gen)))
			if name == "read_fails" {
				bodyreader = errReader(name)
			}

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(tc.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPatch,
				"url",
				bodyreader)

			ha.PatchGeneration(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_DeleteGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  string
		err error
		sc  int
	}{
		"happy_path": {
			id: "1",
			sc: http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"url_decode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{rmErr: tc.err},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{tc.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				nil)

			ha.DeleteGeneration(w, r)

			require.Equal(t, tc.sc, w.Code)
		})
	}
}

func Test_GetGenerationReport(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Entity
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Entity{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"url_decode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					rpt:    v.result,
					rptErr: v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				nil)

			ha.GetGenerationReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeGeneration(l *types.Generation) []byte {
	if l == nil {
		return []byte{}
	}
	result, _ := json.Marshal(l)
	return result
}

func (gm *generationerMock) SelectGenerationIndex(ctx context.Context, cid types.CID) ([]types.Generation, error) {
	return gm.all, gm.allErr
}
func (gm *generationerMock) SelectGenerationsByAttrs(context.Context, types.ReportAttrs, types.CID) ([]types.Generation, error) {
	return gm.all, gm.allErr
}
func (gm *generationerMock) SelectGeneration(ctx context.Context, id types.UUID, cid types.CID) (types.Generation, error) {
	return gm.sel, gm.selErr
}
func (gm *generationerMock) InsertGeneration(ctx context.Context, g types.Generation, cid types.CID) (types.Generation, error) {
	return gm.ins, gm.insErr
}
func (gm *generationerMock) UpdateGeneration(ctx context.Context, g types.Generation, cid types.CID) (types.Generation, error) {
	return gm.upd, gm.updErr
}
func (gm *generationerMock) DeleteGeneration(ctx context.Context, id types.UUID, cid types.CID) error {
	return gm.rmErr
}
func (gm *generationerMock) GeneratedStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	return gm.str, gm.strErr
}
func (gm *generationerMock) UpdateGeneratedStrain(ctx context.Context, gid *types.UUID, sid types.UUID, cid types.CID) error {
	return gm.patchErr
}
func (gm *generationerMock) GenerationReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	return gm.rpt, gm.rptErr
}
