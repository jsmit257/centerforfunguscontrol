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

type eventtyperMock struct {
	selectAllResult []types.EventType
	selectAllErr    error

	selectResult types.EventType
	selectErr    error

	insertResult types.EventType
	insertErr    error

	updateErr error

	deleteErr error

	etr    types.Entity
	etrErr error
}

func Test_GetAllEventTypes(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		result []types.EventType
		err    error
		sc     int
	}{
		"happy_path": {
			result: []types.EventType{},
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
				EventTyper: &eventtyperMock{
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
			ha.GetAllEventTypes(w, r)
			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &[]types.EventType{}, &v.result)
			}
		})
	}
}

func Test_GetEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     string
		result types.EventType
		err    error
		sc     int
	}{
		"happy_path": {
			id:     "1",
			result: types.EventType{},
			sc:     http.StatusOK,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"urldecode_fails": {
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
				EventTyper: &eventtyperMock{
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

			ha.GetEventType(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.EventType{}, &v.result)
			}
		})
	}
}

func Test_PostEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		et     *types.EventType
		result types.EventType
		err    error
		sc     int
	}{
		"happy_path": {
			et:     &types.EventType{},
			result: types.EventType{},
			sc:     http.StatusCreated,
		},
		"missing_evnttype": {
			sc: http.StatusBadRequest,
		},
		"db_error": {
			et:  &types.EventType{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				EventTyper: &eventtyperMock{
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
				bytes.NewReader(serializeEventType(v.et)))

			ha.PostEventType(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.EventType{}, &v.result)
			}
		})
	}
}

func Test_PatchEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		et  *types.EventType
		err error
		sc  int
	}{
		"happy_path": {
			id: "1",
			et: &types.EventType{},
			sc: http.StatusNoContent,
		},
		"missing_id": {
			sc: http.StatusBadRequest,
		},
		"missing_eventtype": {
			id: "1",
			sc: http.StatusBadRequest,
		},
		"urldecode_fails": {
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
		"db_error": {
			id:  "1",
			et:  &types.EventType{},
			err: fmt.Errorf("db error"),
			sc:  http.StatusInternalServerError,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				EventTyper: &eventtyperMock{
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
				bytes.NewReader(serializeEventType(v.et)))

			ha.PatchEventType(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteEventType(t *testing.T) {
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
		"urldecode_fails": {
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
				EventTyper: &eventtyperMock{
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

			ha.DeleteEventType(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_GetEventTypeReport(t *testing.T) {
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
		"urldecode_fails": {
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
				EventTyper: &eventtyperMock{
					etr:    v.result,
					etrErr: v.err,
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

			ha.GetEventTypeReport(w, r)

			require.Equal(t, v.sc, w.Code)
			if w.Code == http.StatusOK {
				checkResult(t, w.Body, &types.Entity{}, &v.result)
			}
		})
	}
}

func serializeEventType(s *types.EventType) []byte {
	if s == nil {
		return []byte{}
	}
	result, _ := json.Marshal(s)
	return result
}

func (vm *eventtyperMock) SelectAllEventTypes(context.Context, types.CID) ([]types.EventType, error) {
	return vm.selectAllResult, vm.selectAllErr
}
func (vm *eventtyperMock) SelectEventType(context.Context, types.UUID, types.CID) (types.EventType, error) {
	return vm.selectResult, vm.selectErr
}
func (vm *eventtyperMock) InsertEventType(context.Context, types.EventType, types.CID) (types.EventType, error) {
	return vm.insertResult, vm.insertErr
}
func (vm *eventtyperMock) UpdateEventType(context.Context, types.UUID, types.EventType, types.CID) error {
	return vm.updateErr
}
func (vm *eventtyperMock) DeleteEventType(context.Context, types.UUID, types.CID) error {
	return vm.deleteErr
}
func (vm *eventtyperMock) EventTypeReport(context.Context, types.UUID, types.CID) (types.Entity, error) {
	return vm.etr, vm.etrErr
}
