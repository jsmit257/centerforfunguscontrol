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

type strainerMock struct {
	selectAllResult []types.Strain
	selectAllErr    error

	selectResult types.Strain
	selectErr    error

	insertResult types.Strain
	insertErr    error

	updateErr error

	deleteErr error
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
			log:   log.WithFields(log.Fields{"test": "Test_GetAllStrains", "case": k}),
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
			log:   log.WithFields(log.Fields{"test": "Test_GetStrain", "case": k}),
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
			sc:     http.StatusOK,
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
			log:   log.WithFields(log.Fields{"test": "Test_PostStrain", "case": k}),
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
			log:   log.WithFields(log.Fields{"test": "Test_PatchStrain", "case": k}),
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
			log:   log.WithFields(log.Fields{"test": "Test_DeleteStrain", "case": k}),
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

			ha.DeleteStrain(w, r)

			require.Equal(t, v.sc, w.Code)
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

func (vm *strainerMock) SelectAllStrains(context.Context, types.CID) ([]types.Strain, error) {
	return vm.selectAllResult, vm.selectAllErr
}

func (vm *strainerMock) SelectStrain(context.Context, types.UUID, types.CID) (types.Strain, error) {
	return vm.selectResult, vm.selectErr
}

func (vm *strainerMock) InsertStrain(context.Context, types.Strain, types.CID) (types.Strain, error) {
	return vm.insertResult, vm.insertErr
}

func (vm *strainerMock) UpdateStrain(context.Context, types.UUID, types.Strain, types.CID) error {
	return vm.updateErr
}

func (vm *strainerMock) DeleteStrain(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
