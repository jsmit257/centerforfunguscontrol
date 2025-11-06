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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
					selectAllResult: tc.result,
					selectAllErr:    tc.err,
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

			ha.GetAllSubstrates(w, r)
			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.Substrate{}, &tc.result)
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
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
				nil)

			ha.GetSubstrate(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Substrate{}, &tc.result)
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
		"read_fails": {
			sc: http.StatusBadRequest,
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

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
					insertResult: tc.result,
					insertErr:    tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeSubstrate(tc.substrate)))
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

			ha.PostSubstrate(w, r)

			require.Equal(t, tc.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Substrate{}, &tc.result)
			}
		})
	}
}

func Test_PatchSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id        types.UUID
		substrate *types.Substrate
		err       error
		sc        int
	}{
		"happy_path": {
			id:        "1",
			substrate: &types.Substrate{},
			sc:        http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
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
			id:        "1",
			substrate: &types.Substrate{},
			err:       fmt.Errorf("db error"),
			sc:        http.StatusInternalServerError,
		},
	}

	for name, tc := range set {
		name, tc := name, tc
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Substrater: &substraterMock{
					updateErr: tc.err,
				},
			},
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bodyreader := io.Reader(bytes.NewReader(serializeSubstrate(tc.substrate)))
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

			ha.PatchSubstrate(w, r)

			require.Equal(t, tc.sc, w.Code)
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
