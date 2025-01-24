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

type stagerMock struct {
	selectAllResult []types.Stage
	selectAllErr    error

	selectResult types.Stage
	selectErr    error

	insertResult types.Stage
	insertErr    error

	updateErr error

	deleteErr error
}

func Test_GetAllStages(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Stage
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Stage{},
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
				Stager: &stagerMock{
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
			ha.GetAllStages(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Stage{}, &v.result)
			}
		})
	}
}

func Test_GetStage(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Stage
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Stage{},
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
				Stager: &stagerMock{
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

			ha.GetStage(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Stage{}, &v.result)
			}
		})
	}
}

func Test_PostStage(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		stage  *types.Stage
		result types.Stage
		err    error
		sc     int
	}{
		"happy_path": {
			stage:  &types.Stage{},
			result: types.Stage{},
			sc:     http.StatusCreated,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			stage: &types.Stage{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Stager: &stagerMock{
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
				bytes.NewReader(serializeStage(v.stage)))

			ha.PostStage(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Stage{}, &v.result)
			}
		})
	}
}

func Test_PatchStage(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		stage *types.Stage
		err   error
		sc    int
	}{
		"happy_path": {
			id:    "1",
			stage: &types.Stage{},
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
			stage: &types.Stage{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Stager: &stagerMock{
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
				bytes.NewReader(serializeStage(v.stage)))

			ha.PatchStage(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteStage(t *testing.T) {
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
				Stager: &stagerMock{
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

			ha.DeleteStage(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func serializeStage(s *types.Stage) []byte {
	if s == nil {
		return []byte{}
	}
	result, _ := json.Marshal(s)
	return result
}

func (vm *stagerMock) SelectAllStages(context.Context, types.CID) ([]types.Stage, error) {
	return vm.selectAllResult, vm.selectAllErr
}

func (vm *stagerMock) SelectStage(context.Context, types.UUID, types.CID) (types.Stage, error) {
	return vm.selectResult, vm.selectErr
}

func (vm *stagerMock) InsertStage(context.Context, types.Stage, types.CID) (types.Stage, error) {
	return vm.insertResult, vm.insertErr
}

func (vm *stagerMock) UpdateStage(context.Context, types.UUID, types.Stage, types.CID) error {
	return vm.updateErr
}

func (vm *stagerMock) DeleteStage(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
