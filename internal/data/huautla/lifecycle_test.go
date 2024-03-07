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

	updateErr error

	deleteErr error
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
			sc:     http.StatusOK,
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
			sc:    http.StatusNoContent,
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

func (vm *lifecyclerMock) SelectLifecycle(context.Context, types.UUID, types.CID) (types.Lifecycle, error) {
	return vm.selectResult, vm.selectErr
}

func (vm *lifecyclerMock) InsertLifecycle(context.Context, types.Lifecycle, types.CID) (types.Lifecycle, error) {
	return vm.insertResult, vm.insertErr
}

func (vm *lifecyclerMock) UpdateLifecycle(context.Context /*types.UUID,*/, types.Lifecycle, types.CID) error {
	return vm.updateErr
}

func (vm *lifecyclerMock) DeleteLifecycle(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
