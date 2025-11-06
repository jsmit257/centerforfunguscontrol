package huautla

import (
	"bytes"
	"context"
	"database/sql"
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectIndexResult: tc.result,
					selectIndexErr:    tc.err,
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
				bytes.NewReader([]byte("")))
			ha.GetLifecycleIndex(w, r)
			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Lifecycle{}, &tc.result)
			}
		})
	}
}

// func Test_GetLifecyclesByAttrs(t *testing.T) {
// 	t.Parallel()

// 	type tc struct {
// 		query  string
// 		result []types.Lifecycle
// 		err    error
// 		sc     int
// 	}

// 	set := map[string]tc{
// 		"happy_strain": {
// 			query:  "strain-id=1234",
// 			result: []types.Lifecycle{},
// 			sc:     http.StatusOK,
// 		},
// 		"unparseable": {
// 			query:  "bulkID=%zzz",
// 			result: []types.Lifecycle{},
// 			sc:     http.StatusBadRequest,
// 		},
// 		"empty_value": {
// 			query:  "strainID",
// 			result: []types.Lifecycle{},
// 			sc:     http.StatusBadRequest,
// 		},
// 		"no_values": {
// 			result: []types.Lifecycle{},
// 			sc:     http.StatusBadRequest,
// 		},
// 		"db_error": {
// 			query: "strain-id=1234",
// 			err:   fmt.Errorf("db error"),
// 			sc:    http.StatusInternalServerError,
// 		},
// 	}

// 	for k, v := range set {
// 		k, v := k, v
// 		ha := &HuautlaAdaptor{
// 			db: &huautlaMock{
// 				Lifecycler: &lifecyclerMock{
// 					selectIndexResult: v.result,
// 					selectIndexErr:    v.err,
// 				},
// 			},
// 			log:   log.WithFields(log.Fields{"test": "Test_GetLifecyclesByAttrs2", "case": k}),
// // 		}

// 		t.Run(k, func(t *testing.T) {
// 			t.Parallel()

// 			w := httptest.NewRecorder()
// 			defer w.Result().Body.Close()
// 			r, _ := http.NewRequestWithContext(
// 				context.WithValue(
// 					metrics.MockServiceContext,
// 					chi.RouteCtxKey,
// 					chi.NewRouteContext()),
// 				http.MethodGet,
// 				fmt.Sprintf("/reports/lifecycles?%s", v.query),
// 				nil)
// 			ha.GetLifecyclesByAttrs(w, r)
// 			require.Equal(t, v.sc, w.Code)
// 			if w.Code == http.StatusOK {
// 				checkResult(t, w.Body, &[]types.Lifecycle{}, &v.result)
// 			}
// 		})
// 	}
// }

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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectResult: tc.result,
					selectErr:    tc.err,
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
				bytes.NewReader([]byte("")))

			ha.GetLifecycle(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Lifecycle{}, &tc.result)
			}
		})
	}
}

func Test_PostLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		lc     *types.Lifecycle
		result types.Lifecycle
		err    error
		sc     int
	}{
		"happy_path": {
			lc:     &types.Lifecycle{},
			result: types.Lifecycle{},
			sc:     http.StatusCreated,
		},
		"read_fails": {
			sc: http.StatusBadRequest,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			lc:  &types.Lifecycle{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					insertResult: tc.result,
					insertErr:    tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeLifecycle(tc.lc)))
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

			ha.PostLifecycle(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Lifecycle{}, &tc.result)
			}
		})
	}
}

func Test_PatchLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		lc  *types.Lifecycle
		err error
		sc  int
	}{
		"happy_path": {
			id: "1",
			lc: &types.Lifecycle{},
			sc: http.StatusOK,
		},
		"missing_id": {
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
			lc:  &types.Lifecycle{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					updateErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeLifecycle(tc.lc)))
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
				http.MethodDelete,
				"url",
				bodyreader)

			ha.PatchLifecycle(w, r)

			require.Equal(t, tc.sc, w.Code)
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					deleteErr: tc.err,
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
				http.MethodDelete,
				"url",
				bytes.NewReader([]byte("")))

			ha.DeleteLifecycle(w, r)

			require.Equal(t, tc.sc, w.Code)
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					rpt:    tc.result,
					rptErr: tc.err,
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
				bytes.NewReader([]byte("")))

			ha.GetLifecycleReport(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &tc.result)
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
