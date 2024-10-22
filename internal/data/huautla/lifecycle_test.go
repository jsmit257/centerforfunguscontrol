package huautla

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type lifecyclerMock struct {
	selectIndexResult []types.Lifecycle
	selectIndexErr    error

	selectResult types.Lifecycle
	selectErr    error

	insertResult types.Lifecycle
	insertErr    error

	updateResult types.Lifecycle
	updateErr    error

	deleteErr error

	rpt    types.Entity
	rptErr error
}

func Test_GetLifecycleIndex(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Lifecycle
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Lifecycle{},
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
				Lifecycler: &lifecyclerMock{
					selectIndexResult: v.result,
					selectIndexErr:    v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_GetLifecycleIndex", "case": k}),
			mtrcs: nil,
		}

		t.Run(k, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				bytes.NewReader([]byte("")))
			ha.GetLifecycleIndex(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Lifecycle{}, &v.result)
			}
		})
	}
}

func Test_GetLifecyclesByAttrs(t *testing.T) {
	t.Parallel()

	type tc struct {
		query  string
		result []types.Lifecycle
		err    error
		sc     int
	}

	set := map[string]tc{
		"happy_strain": {
			query:  "strain-id=1234",
			result: []types.Lifecycle{},
			sc:     http.StatusOK,
		},
		"unparseable": {
			query:  "bulkID=%zzz",
			result: []types.Lifecycle{},
			sc:     http.StatusBadRequest,
		},
		"empty_value": {
			query:  "strainID",
			result: []types.Lifecycle{},
			sc:     http.StatusBadRequest,
		},
		"no_values": {
			result: []types.Lifecycle{},
			sc:     http.StatusBadRequest,
		},
		"db_error": {
			query: "strain-id=1234",
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectIndexResult: v.result,
					selectIndexErr:    v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_GetLifecyclesByAttrs2", "case": k}),
			mtrcs: nil,
		}

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				fmt.Sprintf("/reports/lifecycles?%s", v.query),
				nil)
			ha.GetLifecyclesByAttrs(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Lifecycle{}, &v.result)
			}
		})
	}
}

func Test_GetLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Lifecycle
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Lifecycle{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"url_decode_error": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"missing_row": {
			id:  "abcdefg",
			err: sql.ErrNoRows,
			sc:  http.StatusBadRequest,
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
				Lifecycler: &lifecyclerMock{
					selectResult: v.result,
					selectErr:    v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_GetLifecycle", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				bytes.NewReader([]byte("")))

			ha.GetLifecycle(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Lifecycle{}, &v.result)
			}
		})
	}
}

func Test_PostLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		stage  *types.Lifecycle
		result types.Lifecycle
		err    error
		sc     int
	}{
		"happy_path": {
			stage:  &types.Lifecycle{},
			result: types.Lifecycle{},
			sc:     http.StatusCreated,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			stage: &types.Lifecycle{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					insertResult: v.result,
					insertErr:    v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostLifecycle", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					chi.NewRouteContext()),
				http.MethodGet,
				"url",
				bytes.NewReader(serializeLifecycle(v.stage)))

			ha.PostLifecycle(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Lifecycle{}, &v.result)
			}
		})
	}
}

func Test_PatchLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		stage *types.Lifecycle
		err   error
		sc    int
	}{
		"happy_path": {
			id:    "1",
			stage: &types.Lifecycle{},
			sc:    http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"missing_stage": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:    "1",
			stage: &types.Lifecycle{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					updateErr: v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchLifecycle", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.id)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader(serializeLifecycle(v.stage)))

			ha.PatchLifecycle(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteLifecycle(t *testing.T) {
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

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					deleteErr: v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeleteLifecycle", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteLifecycle(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetLifecycleReport(t *testing.T) {
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
		"missing_row": {
			id:  "abcdefg",
			err: sql.ErrNoRows,
			sc:  http.StatusBadRequest,
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
				Lifecycler: &lifecyclerMock{
					rpt:    v.result,
					rptErr: v.err,
				},
			},
			log:   log.WithFields(log.Fields{"test": "Test_GetLifecycleReport", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{v.id}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodGet,
				"url",
				bytes.NewReader([]byte("")))

			ha.GetLifecycleReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeLifecycle(l *types.Lifecycle) []byte {
	if l == nil {
		return []byte{}
	}
	result, _ := json.Marshal(l)
	return result
}

func (vm *lifecyclerMock) SelectLifecycleIndex(context.Context, types.CID) ([]types.Lifecycle, error) {
	return vm.selectIndexResult, vm.selectIndexErr
}
func (vm *lifecyclerMock) SelectLifecyclesByAttrs(context.Context, types.ReportAttrs, types.CID) ([]types.Lifecycle, error) {
	return vm.selectIndexResult, vm.selectIndexErr
}
func (vm *lifecyclerMock) SelectLifecycle(context.Context, types.UUID, types.CID) (types.Lifecycle, error) {
	return vm.selectResult, vm.selectErr
}
func (vm *lifecyclerMock) InsertLifecycle(context.Context, types.Lifecycle, types.CID) (types.Lifecycle, error) {
	return vm.insertResult, vm.insertErr
}
func (vm *lifecyclerMock) UpdateLifecycle(context.Context, types.Lifecycle, types.CID) (types.Lifecycle, error) {
	return vm.updateResult, vm.updateErr
}
func (vm *lifecyclerMock) DeleteLifecycle(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
func (vm *lifecyclerMock) LifecycleReport(context.Context, types.UUID, types.CID) (types.Entity, error) {
	return vm.rpt, vm.rptErr
}
