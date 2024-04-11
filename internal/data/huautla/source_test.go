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

type sourcerMock struct {
	addStrainErr,
	addEventErr,
	changeErr,
	removeErr error
}

func Test_PostStrainSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		s      *types.Source
		genErr error
		srcErr error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			s:  &types.Source{UUID: "happy_path"},
			sc: http.StatusCreated,
		},
		"strain_error": {
			g:      types.Generation{UUID: "strain_error"},
			s:      &types.Source{UUID: "strain_error"},
			srcErr: fmt.Errorf("strain_error"),
			sc:     http.StatusInternalServerError,
		},
		"generation_error": {
			g:      types.Generation{UUID: "generation_error"},
			s:      &types.Source{UUID: "generation_error"},
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
				Sourcer: &sourcerMock{addStrainErr: v.srcErr},
			},
			log:   log.WithFields(log.Fields{"test": "PostStrainSource", "case": k}),
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
				bytes.NewReader(serializeSource(v.s)))

			ha.PostStrainSource(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PostEventSource(t *testing.T) {
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
				Sourcer: &sourcerMock{addEventErr: v.evtErr},
			},
			log:   log.WithFields(log.Fields{"test": "PostStrainSource", "case": k}),
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

			ha.PostEventSource(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_PatchSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		id     string
		s      *types.Source
		genErr error
		srcErr error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			id: "happy_path",
			s:  &types.Source{UUID: "happy_path"},
			sc: http.StatusOK,
		},
		"source_error": {
			g:      types.Generation{UUID: "source_error"},
			id:     "source_error",
			s:      &types.Source{UUID: "source_error"},
			srcErr: fmt.Errorf("source_error"),
			sc:     http.StatusInternalServerError,
		},
		"generation_error": {
			g:      types.Generation{UUID: "generation_error"},
			id:     "generation_error",
			s:      &types.Source{UUID: "generation_error"},
			genErr: fmt.Errorf("generation_error"),
			sc:     http.StatusInternalServerError,
		},
		"missing_body": {
			g:  types.Generation{UUID: "missing_body"},
			id: "missing_body",
			sc: http.StatusBadRequest,
		},
		"missing_event_id": {
			g:  types.Generation{UUID: "missing_event_id"},
			sc: http.StatusBadRequest,
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
				Sourcer: &sourcerMock{changeErr: v.srcErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_PatchSource", "case": k}),
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
				bytes.NewReader(serializeSource(v.s)))

			ha.PatchSource(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func Test_DeleteSource(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		id     string
		genErr error
		srcErr error
		sc     int
	}{
		"happy_path": {
			g:  types.Generation{UUID: "happy_path"},
			id: "happy_path",
			sc: http.StatusOK,
		},
		"source_error": {
			g:      types.Generation{UUID: "source_error"},
			id:     "source_error",
			srcErr: fmt.Errorf("source_error"),
			sc:     http.StatusInternalServerError,
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
				Sourcer: &sourcerMock{removeErr: v.srcErr},
			},
			log:   log.WithFields(log.Fields{"test": "Test_DeleteSource", "case": k}),
			mtrcs: nil,
		}
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			defer w.Result().Body.Close()
			rctx := chi.NewRouteContext()
			rctx.URLParams = chi.RouteParams{Keys: []string{"g_id", "s_id"}, Values: []string{
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
				bytes.NewReader(serializeSource(nil)))

			ha.DeleteSource(w, r)

			require.Equal(t, v.sc, w.Code)
		})
	}
}

func serializeSource(e *types.Source) []byte {
	if e == nil {
		return []byte{}
	}
	result, _ := json.Marshal(e)
	return result
}

func (sm *sourcerMock) AddStrainSource(ctx context.Context, g *types.Generation, s types.Source, cid types.CID) error {
	return sm.addStrainErr
}
func (sm *sourcerMock) AddEventSource(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) error {
	return sm.addEventErr
}
func (sm *sourcerMock) ChangeSource(ctx context.Context, g *types.Generation, s types.Source, cid types.CID) error {
	return sm.changeErr
}
func (sm *sourcerMock) RemoveSource(ctx context.Context, g *types.Generation, id types.UUID, cid types.CID) error {
	return sm.removeErr
}
