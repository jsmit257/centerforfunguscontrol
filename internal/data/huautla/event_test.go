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

type eventerMock struct {
	byTypeResult []types.Event
	byTypeErr    error

	selectResult types.Event
	selectErr    error

	addErr,
	changeErr,
	rmErr error
}

func Test_PostEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		l     types.Lifecycle
		e     *types.Event
		lcErr error
		evErr error
		sc    int
	}{
		"happy_path": {
			l:  types.Lifecycle{UUID: "happy_path"},
			e:  &types.Event{UUID: "happy_path"},
			sc: http.StatusCreated,
		},
		"event_error": {
			l:     types.Lifecycle{UUID: "event_error"},
			e:     &types.Event{UUID: "event_error"},
			evErr: fmt.Errorf("event_error"),
			sc:    http.StatusInternalServerError,
		},
		"lifecycle_error": {
			l:     types.Lifecycle{UUID: "lifecycle_error"},
			e:     &types.Event{UUID: "lifecycle_error"},
			lcErr: fmt.Errorf("lifecycle_error"),
			sc:    http.StatusInternalServerError,
		},
		"missing_body": {
			l:  types.Lifecycle{UUID: "missing_body"},
			sc: http.StatusBadRequest,
		},
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectResult: v.l,
					selectErr:    v.lcErr,
				},
				Eventer: &eventerMock{addErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.l.UUID)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeEvent(v.e)))

			ha.PostEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		l     types.Lifecycle
		id    string
		e     *types.Event
		lcErr error
		evErr error
		sc    int
	}{
		"happy_path": {
			l:  types.Lifecycle{UUID: "happy_path"},
			id: "happy_path",
			e:  &types.Event{UUID: "happy_path"},
			sc: http.StatusOK,
		},
		"event_error": {
			l:     types.Lifecycle{UUID: "event_error"},
			id:    "event_error",
			e:     &types.Event{UUID: "event_error"},
			evErr: fmt.Errorf("event_error"),
			sc:    http.StatusInternalServerError,
		},
		"lifecycle_error": {
			l:     types.Lifecycle{UUID: "lifecycle_error"},
			id:    "lifecycle_error",
			e:     &types.Event{UUID: "lifecycle_error"},
			lcErr: fmt.Errorf("lifecycle_error"),
			sc:    http.StatusInternalServerError,
		},
		"missing_body": {
			l:  types.Lifecycle{UUID: "missing_body"},
			id: "missing_body",
			sc: http.StatusBadRequest,
		},
		"missing_event_id": {
			l:  types.Lifecycle{UUID: "missing_event_id"},
			sc: http.StatusBadRequest,
		},
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectResult: v.l,
					selectErr:    v.lcErr,
				},
				Eventer: &eventerMock{changeErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"lc_id", "ev_id"}, Values: []string{
				string(v.l.UUID),
				v.id,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeEvent(v.e)))

			ha.PatchEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		l     types.Lifecycle
		id    string
		lcErr error
		evErr error
		sc    int
	}{
		"happy_path": {
			l:  types.Lifecycle{UUID: "happy_path"},
			id: "happy_path",
			sc: http.StatusOK,
		},
		"event_error": {
			l:     types.Lifecycle{UUID: "event_error"},
			id:    "event_error",
			evErr: fmt.Errorf("event_error"),
			sc:    http.StatusInternalServerError,
		},
		"lifecycle_error": {
			l:     types.Lifecycle{UUID: "lifecycle_error"},
			id:    "lifecycle_error",
			lcErr: fmt.Errorf("lifecycle_error"),
			sc:    http.StatusInternalServerError,
		},
		"missing_event_id": {
			l:  types.Lifecycle{UUID: "missing_event_id"},
			sc: http.StatusBadRequest,
		},
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Lifecycler: &lifecyclerMock{
					selectResult: v.l,
					selectErr:    v.lcErr,
				},
				Eventer: &eventerMock{rmErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"lc_id", "ev_id"}, Values: []string{
				string(v.l.UUID),
				v.id,
			}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeEvent(nil)))

			ha.DeleteEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func serializeEvent(e *types.Event) []byte {
	if e == nil {
		return []byte{}
	}
	result, _ := json.Marshal(e)
	return result
}

func (em *eventerMock) SelectByEventType(ctx context.Context, et types.EventType, cid types.CID) ([]types.Event, error) {
	return em.byTypeResult, em.byTypeErr
}
func (em *eventerMock) SelectEvent(ctx context.Context, id types.UUID, cid types.CID) (types.Event, error) {
	return em.selectResult, em.selectErr
}
func (em *eventerMock) GetLifecycleEvents(ctx context.Context, lc *types.Lifecycle, cid types.CID) error {
	return nil
}
func (em *eventerMock) AddEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	return em.addErr
}
func (em *eventerMock) ChangeEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	return em.changeErr
}
func (em *eventerMock) RemoveEvent(ctx context.Context, lc *types.Lifecycle, id types.UUID, cid types.CID) error {
	return em.rmErr
}
