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
	addErr       error
	changeResult types.Event
	changeErr    error
	rmErr        error

	addGenerationErr       error
	changeGenerationResult types.Event
	changeGenerationErr    error
	rmGenerationErr        error
}

func Test_PostLifecycleEvent(t *testing.T) {
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
		"urldecode_error": {
			l:  types.Lifecycle{UUID: "%zzz"},
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
				LifecycleEventer: &eventerMock{addErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostLifecycleEvent", "case": k}),
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

			ha.PostLifecycleEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchLifecycleEvent(t *testing.T) {
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
		"lc_decode_error": {
			l:  types.Lifecycle{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"ev_decode_error": {
			l:  types.Lifecycle{UUID: "happy_path"},
			id: "%zzz",
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
				LifecycleEventer: &eventerMock{changeErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchLifecycleEvent", "case": k}),
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

			ha.PatchLifecycleEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteLifecycleEvent(t *testing.T) {
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
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
		"lc_decode_error": {
			l:  types.Lifecycle{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"ev_decode_error": {
			l:  types.Lifecycle{UUID: "happy_path"},
			id: "%zzz",
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
				LifecycleEventer: &eventerMock{rmErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeleteLifecycleEvent", "case": k}),
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

			ha.DeleteLifecycleEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PostGenerationEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		e      *types.Event
		genErr error
		evtErr error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			e:  &types.Event{UUID: "happy_path"},
			sc: http.StatusCreated,
		},
		"event_error": {
			g:      types.Generation{UUID: "event_error"},
			e:      &types.Event{UUID: "event_error"},
			evtErr: fmt.Errorf("event_error"),
			sc:     http.StatusInternalServerError,
		},
		"generation_error": {
			g:      types.Generation{UUID: "generation_error"},
			e:      &types.Event{UUID: "generation_error"},
			genErr: fmt.Errorf("generation_error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_body": {
			g:  types.Generation{UUID: "missing_body"},
			sc: http.StatusBadRequest,
		},
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
		"urldecode_error": {
			g:  types.Generation{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					sel:    v.g,
					selErr: v.genErr,
				},
				GenerationEventer: &eventerMock{addGenerationErr: v.evtErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PostGenerationEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.g.UUID)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeEvent(v.e)))

			ha.PostGenerationEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchGenerationEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		e      *types.Event
		genErr error
		evtErr error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			e:  &types.Event{UUID: "happy_path"},
			sc: http.StatusOK,
		},
		"event_error": {
			g:      types.Generation{UUID: "event_error"},
			e:      &types.Event{UUID: "event_error"},
			evtErr: fmt.Errorf("event_error"),
			sc:     http.StatusInternalServerError,
		},
		"generation_error": {
			g:      types.Generation{UUID: "generation_error"},
			e:      &types.Event{UUID: "generation_error"},
			genErr: fmt.Errorf("generation_error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_body": {
			g:  types.Generation{UUID: "missing_body"},
			sc: http.StatusBadRequest,
		},
		"missing_event_id": {
			g:  types.Generation{UUID: "missing_event_id"},
			sc: http.StatusBadRequest,
		},
		"missing_generation": {
			sc: http.StatusBadRequest,
		},
		"lc_decode_error": {
			g:  types.Generation{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					sel:    v.g,
					selErr: v.genErr,
				},
				GenerationEventer: &eventerMock{changeGenerationErr: v.evtErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchGenerationEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"id"}, Values: []string{string(v.g.UUID)}}
			r, _ := http.NewRequestWithContext(
				context.WithValue(
					context.Background(),
					chi.RouteCtxKey,
					rctx),
				http.MethodPost,
				"url",
				bytes.NewReader(serializeEvent(v.e)))

			ha.PatchGenerationEvent(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteGenerationEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		id     string
		genErr error
		evErr  error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			id: "happy_path",
			sc: http.StatusOK,
		},
		"event_error": {
			g:     types.Generation{UUID: "event_error"},
			id:    "event_error",
			evErr: fmt.Errorf("event_error"),
			sc:    http.StatusInternalServerError,
		},
		"lifecycle_error": {
			g:      types.Generation{UUID: "lifecycle_error"},
			id:     "lifecycle_error",
			genErr: fmt.Errorf("lifecycle_error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_lifecycle": {
			sc: http.StatusBadRequest,
		},
		"lc_decode_error": {
			g:  types.Generation{UUID: "%zzz"},
			sc: http.StatusBadRequest,
		},
		"ev_decode_error": {
			g:  types.Generation{UUID: "happy_path"},
			id: "%zzz",
			sc: http.StatusBadRequest,
		},
	}

	for k, v := range set {
		k, v := k, v
		ha := &HuautlaAdaptor{
			db: &huautlaMock{
				Generationer: &generationerMock{
					sel:    v.g,
					selErr: v.genErr,
				},
				GenerationEventer: &eventerMock{rmGenerationErr: v.evErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeleteGenerationEvent", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"g_id", "ev_id"}, Values: []string{
				string(v.g.UUID),
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

			ha.DeleteGenerationEvent(w, r)

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

func (em *eventerMock) GetLifecycleEvents(ctx context.Context, lc *types.Lifecycle, cid types.CID) error {
	return nil
}
func (em *eventerMock) AddLifecycleEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	return em.addErr
}
func (em *eventerMock) ChangeLifecycleEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) (types.Event, error) {
	return em.changeResult, em.changeErr
}
func (em *eventerMock) RemoveLifecycleEvent(ctx context.Context, lc *types.Lifecycle, id types.UUID, cid types.CID) error {
	return em.rmErr
}

func (em *eventerMock) GetGenerationEvents(ctx context.Context, g *types.Generation, cid types.CID) error {
	return nil
}
func (em *eventerMock) AddGenerationEvent(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) error {
	return em.addGenerationErr
}
func (em *eventerMock) ChangeGenerationEvent(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) (types.Event, error) {
	return em.changeGenerationResult, em.changeGenerationErr
}
func (em *eventerMock) RemoveGenerationEvent(ctx context.Context, g *types.Generation, id types.UUID, cid types.CID) error {
	return em.rmGenerationErr
}
