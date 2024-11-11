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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type substraterMock struct {
	selectAllResult []types.Substrate
	selectAllErr    error

	selectResult types.Substrate
	selectErr    error

	insertResult types.Substrate
	insertErr    error

	updateErr error

	deleteErr error

	rpt    types.Entity
	rptErr error
}

func Test_GetAllSubstrates(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.Substrate
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.Substrate{},
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
				Substrater: &substraterMock{
					selectAllResult: v.result,
					selectAllErr:    v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_GetAllSubstrates", "case": k}),
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
			ha.GetAllSubstrates(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Substrate{}, &v.result)
			}
		})
	}
}

func Test_GetSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.Substrate
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.Substrate{},
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
				Substrater: &substraterMock{
					selectResult: v.result,
					selectErr:    v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_GetSubstrate", "case": k}),
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

			ha.GetSubstrate(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Substrate{}, &v.result)
			}
		})
	}
}

func Test_PostSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		substrate *types.Substrate
		result    types.Substrate
		err       error
		sc        int
	}{
		"happy_path": {
			substrate: &types.Substrate{},
			result:    types.Substrate{},
			sc:        http.StatusCreated,
		},
		"missing_stage": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			substrate: &types.Substrate{},
			err:       fmt.Errorf("db error"),
			sc:        http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
					insertResult: v.result,
					insertErr:    v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_PostSubstrate", "case": k}),
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
				bytes.NewReader(serializeSubstrate(v.substrate)))

			ha.PostSubstrate(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Substrate{}, &v.result)
			}
		})
	}
}

func Test_PatchSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		stage *types.Substrate
		err   error
		sc    int
	}{
		"happy_path": {
			id:    "1",
			stage: &types.Substrate{},
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
			stage: &types.Substrate{},
			err:   fmt.Errorf("db error"),
			sc:    http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
					updateErr: v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_PatchSubstrate", "case": k}),
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
				bytes.NewReader(serializeSubstrate(v.stage)))

			ha.PatchSubstrate(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteSubstrate(t *testing.T) {
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
				Substrater: &substraterMock{
					deleteErr: v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_DeleteSubstrate", "case": k}),
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

			ha.DeleteSubstrate(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetSubstrateReport(t *testing.T) {
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
				Substrater: &substraterMock{
					rpt:    v.result,
					rptErr: v.err,
				},
			},
			log:   logrus.WithFields(logrus.Fields{"test": "Test_GetSubstrateReport", "case": k}),
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
				nil)

			ha.GetSubstrateReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeSubstrate(s *types.Substrate) []byte {
	if s == nil {
		return []byte{}
	}
	result, _ := json.Marshal(s)
	return result
}

func (vm *substraterMock) SelectAllSubstrates(context.Context, types.CID) ([]types.Substrate, error) {
	return vm.selectAllResult, vm.selectAllErr
}
func (vm *substraterMock) SelectSubstrate(context.Context, types.UUID, types.CID) (types.Substrate, error) {
	return vm.selectResult, vm.selectErr
}
func (vm *substraterMock) InsertSubstrate(context.Context, types.Substrate, types.CID) (types.Substrate, error) {
	return vm.insertResult, vm.insertErr
}
func (vm *substraterMock) UpdateSubstrate(context.Context, types.UUID, types.Substrate, types.CID) error {
	return vm.updateErr
}
func (vm *substraterMock) DeleteSubstrate(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
func (vm *substraterMock) SubstrateReport(context.Context, types.UUID, types.CID) (types.Entity, error) {
	return vm.rpt, vm.rptErr
}
