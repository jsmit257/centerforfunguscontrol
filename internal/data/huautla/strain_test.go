package huautla

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

type strainerMock struct {
	selectAllResult []types.Strain
	selectAllErr    error

	selectResult types.Strain
	selectErr    error

	insertResult types.Strain
	insertErr    error

	updateErr error

	deleteErr error

	str    types.Strain
	strErr error

	patchErr error

	rpt    types.Entity
	rptErr error
}

func Test_GetAllStrains(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Strain
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Strain{},
			sc:     http.StatusOK,
		},
		"db_error": {
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					selectAllResult: v.result,
					selectAllErr:    v.err,
				},
			},
		}

		t.Run(k, func(t *testing.T) {
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
				bytes.NewReader([]byte("")))
			ha.GetAllStrains(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Strain{}, &v.result)
			}
		})
	}
}

func Test_GetStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Strain
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Strain{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
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
				Strainer: &strainerMock{
					selectResult: v.result,
					selectErr:    v.err,
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
				bytes.NewReader([]byte("")))

			ha.GetStrain(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Strain{}, &v.result)
			}
		})
	}
}

func Test_PostStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		stage  *types.Strain
		result types.Strain
		err    error
		sc     int
	}{
		"happy_path": {
			stage:  &types.Strain{},
			result: types.Strain{},
			sc:     http.StatusCreated,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			stage: &types.Strain{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					insertResult: v.result,
					insertErr:    v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
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
				bytes.NewReader(serializeStrain(v.stage)))

			ha.PostStrain(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Strain{}, &v.result)
			}
		})
	}
}

func Test_PatchStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		stage *types.Strain
		err   error
		sc    int
	}{
		"happy_path": {
			id:    "1",
			stage: &types.Strain{},
			sc:    http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_stage": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:    "1",
			stage: &types.Strain{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Strainer: &strainerMock{
					updateErr: v.err,
				},
			},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader(serializeStrain(v.stage)))

			ha.PatchStrain(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteStrain(t *testing.T) {
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
		"urldecode_error": {
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
				Strainer: &strainerMock{
					deleteErr: v.err,
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
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteStrain(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetGeneratedStrain(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		id     string
		result types.Strain
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Strain{},
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
				Strainer: &strainerMock{
					str:    v.result,
					strErr: v.err,
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
				bytes.NewReader([]byte{}))
			ha.GetGeneratedStrain(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Strain{}, &v.result)
			}
		})
	}
}

func Test_PatchGeneratedStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gid,
		sid types.UUID
		err error
		sc  int
	}{
		"happy_path": {
			gid: "1",
			sid: "1",
			sc:  http.StatusNoContent,
		},
		"missing_generation": {
			sc:  http.StatusBadRequest,
			sid: "1",
		},
		"bad_generation": {
			sc:  http.StatusBadRequest,
			gid: "%zzz",
		},
		"missing_strain": {
			gid: "1",
			sc:  http.StatusBadRequest,
		},
		"bad_strain": {
			sc:  http.StatusBadRequest,
			gid: "1",
			sid: "%zzz",
		},
		"db_error": {
			gid: "1",
			sid: "1",
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{Strainer: &strainerMock{patchErr: v.err}},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"gid", "sid"}, Values: []string{string(v.gid), string(v.sid)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPatch,
				"url",
				nil)

			ha.PatchGeneratedStrain(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteGeneratedStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		sid types.UUID
		err error
		sc  int
	}{
		"happy_path": {
			sid: "1",
			sc:  http.StatusNoContent,
		},
		// // the rest are checked in Test_PatchGeneratedStrain
		// "missing_strain": {
		// 	sc: http.StatusBadRequest,
		// },
		// "bad_strain": {
		// 	sc:  http.StatusBadRequest,
		// 	sid: "%zzz",
		// },
		// "db_error": {
		// 	sid: "1",
		// 	err: fmt.Errorf("db error"),
		// 	sc:  http.StatusInternalServerError,
		// },
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{Strainer: &strainerMock{patchErr: v.err}},
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"gid", "sid"}, Values: []string{"1", string(v.sid)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					metrics.MockServiceContext,
					chi.RouteCtxKey,
					rctx),
				http.MethodPatch,
				"url",
				nil)

			ha.DeleteGeneratedStrain(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetStrainReport(t *testing.T) {
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
		"urldecode_error": {
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
				Strainer: &strainerMock{
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
				bytes.NewReader([]byte("")))

			ha.GetStrainReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeStrain(s *types.Strain) []byte {
	if s == nil {
		return []byte{}
	}
	result, _ := json.Marshal(s)
	return result
}

func (sm *strainerMock) SelectAllStrains(context.Context, types.CID) ([]types.Strain, error) {
	return sm.selectAllResult, sm.selectAllErr
}
func (sm *strainerMock) SelectStrain(context.Context, types.UUID, types.CID) (types.Strain, error) {
	return sm.selectResult, sm.selectErr
}
func (sm *strainerMock) InsertStrain(context.Context, types.Strain, types.CID) (types.Strain, error) {
	return sm.insertResult, sm.insertErr
}
func (sm *strainerMock) UpdateStrain(context.Context, types.UUID, types.Strain, types.CID) error {
	return sm.updateErr
}
func (sm *strainerMock) DeleteStrain(context.Context, types.UUID, types.CID) error {
	return sm.deleteErr
}
func (sm *strainerMock) GeneratedStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	return sm.str, sm.strErr
}
func (sm *strainerMock) UpdateGeneratedStrain(ctx context.Context, gid *types.UUID, sid types.UUID, cid types.CID) error {
	return sm.patchErr
}
func (sm *strainerMock) StrainReport(context.Context, types.UUID, types.CID) (types.Entity, error) {
	return sm.rpt, sm.rptErr
}
